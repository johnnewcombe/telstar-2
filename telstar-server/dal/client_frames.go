package dal

// Database access layer methods that are not restricted by user Id

import (
	"context"
	"errors"
	"fmt"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar-library/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"slices"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBNAME               = "telstardb"
	REGEXP               = "p[0-9]"
	REGEXS               = "s[0-9]"
	AUTH_COLLECTION      = "system-auth"
	ERROR_SCOPE          = "user does not have sufficient scope to perform this task"
	ERROR_AUTHENTICATION = "user has not have authenticated"
	ERROR_FRAMEDECODE    = "decoding frame %s: %v"
)

func GetFrames(connectionUrl string, primaryDb bool) ([]types.Frame, error) {

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
			//TODO should this n=be a panic??
		}
	}()

	var (
		//frame            Frame
		result           []types.Frame
		collectionNames  []string
		pCollectionNames []string
		sCollectionNames []string
		frameDocs        *mongo.Cursor
	)

	if pCollectionNames, sCollectionNames, err = getCollectionNames(connectionUrl); err != nil {
		return result, err
	}

	if primaryDb {
		collectionNames = pCollectionNames
	} else {
		collectionNames = sCollectionNames
	}

	// create the filter
	filter := bson.D{{}}

	for _, collection := range collectionNames {
		if frameDocs, err = client.Database(DBNAME).Collection(collection).Find(ctx, filter); err != nil {
			return result, err
		}
		for frameDocs.Next(ctx) {
			var frame types.Frame
			err = frameDocs.Decode(&frame)
			result = append(result, frame)
		}
	}

	return result, nil
}

func GetFramesByUser(connectionUrl string, primaryDb bool, user types.User) ([]types.Frame, error) {

	var (
		result           []types.Frame
		collectionNames  []string
		pCollectionNames []string
		sCollectionNames []string
		frameDocs        *mongo.Cursor
		//user             User
	)
	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
			//TODO should this n=be a panic??
		}
	}()

	if pCollectionNames, sCollectionNames, err = getCollectionNames(connectionUrl); err != nil {
		return result, err
	}

	if primaryDb {
		collectionNames = pCollectionNames
	} else {
		collectionNames = sCollectionNames
	}

	// sort them, not essential but helps with debugging
	slices.Sort(collectionNames)

	// create the filter
	filter := bson.D{{}}

	for _, collection := range collectionNames {

		if frameDocs, err = client.Database(DBNAME).Collection(collection).Find(ctx, filter); err != nil {
			return result, err
		}

		for frameDocs.Next(ctx) {

			var frame types.Frame

			err = frameDocs.Decode(&frame)

			if err != nil {
				return nil, fmt.Errorf(ERROR_FRAMEDECODE, frame.GetPageId(), err)
			}

			if !user.Authenticated {
				return []types.Frame{}, errors.New(ERROR_AUTHENTICATION)
			}

			if !frame.IsValid() {
				// TODO specific DB clean type command
				// FIXME this does not delete anything is the data actually in the database?
				//       The dodgy data stems from a completely blank frame
				//       i.e. page-no: 0, frame-id: "" but early indications suggest that
				//       there is an ID associated with the frame.
				//coll := client.Database(DBNAME).Collection(collection)
				//delFilter := bson.M{"pid.page-no": frame.PID.PageNumber, "pid.frame-id": frame.PID.FrameId}
				//res, err := coll.DeleteMany(ctx, delFilter)
				//fmt.Printf("%d, %v",res.DeletedCount, err)

				// dodgy data (see above) so ignore.
				//continue
			}

			if user.IsInScope(frame.PID.PageNumber) {
				result = append(result, frame)
			}
		}
	}

	types.SortFrames(result)
	return result, nil

}

func GetFramesByCollection(connectionUrl string, collectionName string) ([]types.Frame, error) {

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
			//TODO should this n=be a panic??
		}
	}()

	// define the result type
	var result []types.Frame

	// get the collection
	collection := client.Database(DBNAME).Collection(collectionName)

	// create the filter
	filter := bson.D{{}}

	// perform the search and get a cursor
	cur, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("finding collection: %s: %v", collectionName, err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		var f types.Frame
		err := cur.Decode(&f)
		if err != nil {
			return nil, fmt.Errorf("decoding frame %s: %v", f.GetPageId(), err)

		}

		result = append(result, f)

	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return result, nil

}

func GetFrame(connectionUrl string, pageNo int, frameId string, primaryDb bool, visibleOnly bool) (types.Frame, error) {

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// define the result type
	var result types.Frame
	var filter bson.M
	if visibleOnly {
		filter = bson.M{"pid.page-no": pageNo, "pid.frame-id": frameId, "visible": true}
	} else {
		filter = bson.M{"pid.page-no": pageNo, "pid.frame-id": frameId}
	}

	collectionName, err := getCollectionName(pageNo, primaryDb)
	if err != nil {
		return result, fmt.Errorf("getting collection name for frame %d%v: %v", pageNo, frameId, err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("finding frame %d%v: %v", pageNo, frameId, err)
	}

	// TODO check for invalid PID e.g. zero page-no and "" frame-id
	return result, nil
}

func GetFrameByUser(connectionUrl string, pageNo int, frameId string, primaryDb bool, user types.User) (types.Frame, error) {

	var (
		//user   User
		result types.Frame
	)

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//if user, err = GetUser(connectionUrl, authUser); err != nil {
	//	return result, err
	//}

	filter := bson.M{"pid.page-no": pageNo, "pid.frame-id": frameId}

	collectionName, err := getCollectionName(pageNo, primaryDb)
	if err != nil {
		return result, fmt.Errorf("getting collection name for frame %d%v: %v", pageNo, frameId, err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("finding frame %d%v: %v", pageNo, frameId, err)
	}

	if !user.Authenticated {
		return result, errors.New(ERROR_AUTHENTICATION)
	}

	if !user.IsInScope(result.PID.PageNumber) {
		return result, errors.New(ERROR_SCOPE)
	}
	return result, nil
}

func InsertFrame(connectionUrl string, frame types.Frame, primaryDb bool) (bool, error) {

	// This could insert multiple tmp with the same ID if called multiple times
	// unless prevented with a unique index. Consider using InsertOrReplace or adding a
	// unique index on pageNo/frameId

	var err error

	if !frame.IsValid() {
		return false, fmt.Errorf("invalid frame data for frame %s", frame.GetPageId())
	}

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collectionName, err := getCollectionName(frame.PID.PageNumber, primaryDb)
	if err != nil {
		return false, fmt.Errorf("getting collection name for frame %v: %v", frame.GetPageId(), err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	// marshall the data
	data, err := bson.Marshal(frame)
	if err != nil {
		return false, fmt.Errorf("converting frame data for frame %v to BSON: %v", frame.GetPageId(), err)
	}
	res, err := collection.InsertOne(ctx, data)
	if err != nil || res.InsertedID == nil {
		return false, fmt.Errorf("inserting frame %s: %v", frame.GetPageId(), err)
	}
	return true, err
}

func InsertOrReplaceFrame(connectionUrl string, frame types.Frame, primaryDb bool) error {

	var (
		err error
	)

	if !frame.IsValid() {
		return fmt.Errorf("invalid frame data for frame %s", frame.GetPageId())
	}

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collectionName, err := getCollectionName(frame.PID.PageNumber, primaryDb)
	if err != nil {
		return fmt.Errorf("getting collection name for frame %v: %v", frame.GetPageId(), err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	filter := bson.M{"pid.page-no": frame.PID.PageNumber, "pid.frame-id": frame.PID.FrameId}

	// remove the ID or an update might fail i.e. when taking frames from one db to another
	frame.ID = primitive.NilObjectID

	// marshall the data
	data, err := bson.Marshal(frame)
	if err != nil {
		return fmt.Errorf("converting frame data for frame %v to BSON: %v", frame.GetPageId(), err)
	}
	// data good so replace
	res, err := collection.ReplaceOne(ctx, filter, data)
	if err != nil {
		// error detected
		return err
	}

	// if we matched 0 items then insert
	if res.MatchedCount == 0 {
		res, err := collection.InsertOne(ctx, data)
		if err != nil || res.InsertedID == nil {
			return fmt.Errorf("inserting frame %s: %v", frame.GetPageId(), err)
		}
	}

	return err

}

func InsertOrReplaceFrameByUser(connectionUrl string, frame types.Frame, primaryDb bool, user types.User) error {

	var (
		//user User
		err error
	)

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//if user, err = GetUser(connectionUrl, authUser); err != nil {
	//	return err
	//}

	if !user.Authenticated {
		return errors.New(ERROR_AUTHENTICATION)
	}

	if !user.IsInScope(frame.PID.PageNumber) {
		return errors.New(ERROR_SCOPE)
	}

	if !frame.IsValid() {
		return fmt.Errorf("invalid frame data for frame %s", frame.GetPageId())
	}

	collectionName, err := getCollectionName(frame.PID.PageNumber, primaryDb)
	if err != nil {
		return fmt.Errorf("getting collection name for frame %v: %v", frame.GetPageId(), err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	filter := bson.M{"pid.page-no": frame.PID.PageNumber, "pid.frame-id": frame.PID.FrameId}

	// remove the ID or an update might fail i.e. when taking frames from one db to another
	frame.ID = primitive.NilObjectID

	// marshall the data
	data, err := bson.Marshal(frame)
	if err != nil {
		return fmt.Errorf("converting frame data for frame %v to BSON: %v", frame.GetPageId(), err)
	}
	// data good so replace
	res, err := collection.ReplaceOne(ctx, filter, data)
	if err != nil {
		// error detected
		return err
	}
	if res.MatchedCount == 0 {
		res, err := collection.InsertOne(ctx, data)
		if err != nil || res.InsertedID == nil {
			return fmt.Errorf("inserting frame %s: %v", frame.GetPageId(), err)
		}
	}

	return err

}

func DeleteFrame(connectionUrl string, pageNo int, frameId string, primaryDb bool) (int64, error) {

	var (
		err error
	)

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collectionName, err := getCollectionName(pageNo, primaryDb)
	if err != nil {
		return 0, fmt.Errorf("getting collection name for frame %d%v: %v", pageNo, frameId, err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	filter := bson.M{"pid.page-no": pageNo, "pid.frame-id": frameId}

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, err
}

func DeleteFrameByUser(connectionUrl string, pageNo int, frameId string, primaryDb bool, user types.User) (int64, error) {

	// FIXME sort out user stuff
	var (
		//user User
		err error
	)

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//if user, err = GetUser(connectionUrl, authUser); err != nil {
	//	return 0, err
	//}

	if !user.Authenticated {
		return 0, errors.New(ERROR_AUTHENTICATION)
	}

	if !user.IsInScope(pageNo) {
		return 0, errors.New(ERROR_SCOPE)
	}

	collectionName, err := getCollectionName(pageNo, primaryDb)
	if err != nil {
		return 0, fmt.Errorf("getting collection name for frame %d%v: %v", pageNo, frameId, err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	filter := bson.M{"pid.page-no": pageNo, "pid.frame-id": frameId}

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, err
}

func PurgeFramesByUser(connectionUrl string, pageNo int, frameId string, primaryDb bool, user types.User) (int64, error) {

	var (
		deletedCount   int64
		rFrameId       rune
		collectionName string
		result         *mongo.DeleteResult
		err            error
	)

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if !user.Authenticated {
		return 0, errors.New(ERROR_AUTHENTICATION)
	}

	if !user.IsInScope(pageNo) {
		return 0, errors.New(ERROR_SCOPE)
	}

	if collectionName, err = getCollectionName(pageNo, primaryDb); err != nil {
		return deletedCount, err
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	for {

		filter := bson.M{"pid.page-no": pageNo, "pid.frame-id": frameId}
		logger.LogInfo.Printf("Purging: %d%s", pageNo, frameId)

		result, err = collection.DeleteMany(ctx, filter)
		if err != nil {
			return deletedCount, err
		}

		deletedCount += result.DeletedCount

		// all frames within the page number completed so do any zero routed frames
		// this is repeated until while pageNo is less that 16 chars long
		//pageNo *= 10
		//pageNoS = strconv.Itoa(pageNo)
		if pageNo, rFrameId, err = utils.GetFollowOnPID(pageNo, []rune(frameId)[0]); err != nil {
			return deletedCount, err
		}
		frameId = string(rFrameId)

		if len(strconv.Itoa(pageNo)) > 9 {
			break
		}
	}
	return deletedCount, nil
}

func PublishFrameByUser(connectionUrl string, pageNo int, frameId string, user types.User) error {

	var (
		//user User
		frame     types.Frame
		err       error
		frameData []byte
	)

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//if user, err = GetUser(connectionUrl, authUser); err != nil {
	//	return 0, err
	//}

	if !user.Authenticated {
		return errors.New(ERROR_AUTHENTICATION)
	}

	if !user.IsInScope(pageNo) {
		return errors.New(ERROR_SCOPE)
	}

	// get the secondary collection
	collectionName, err := getCollectionName(pageNo, false)
	if err != nil {
		return fmt.Errorf("getting collection name for frame %d%v: %v", pageNo, frameId, err)
	}
	collection := client.Database(DBNAME).Collection(collectionName)

	// set the filter
	filter := bson.M{"pid.page-no": pageNo, "pid.frame-id": frameId}

	// get the frame
	err = collection.FindOne(ctx, filter).Decode(&frame)
	if err != nil {
		return fmt.Errorf("finding frame %d%v: %v", pageNo, frameId, err)
	}

	// get the primary collection
	collectionName, err = getCollectionName(pageNo, true)
	if err != nil {
		return fmt.Errorf("getting collection name for frame %d%v: %v", pageNo, frameId, err)
	}
	collection = client.Database(DBNAME).Collection(collectionName)

	// create the filter
	filter = bson.M{"pid.page-no": frame.PID.PageNumber, "pid.frame-id": frame.PID.FrameId}

	// remove the ID or an update might fail i.e. a frame that already exists in the primary will have its own ID
	frame.ID = primitive.NilObjectID

	// marshall the frame
	frameData, err = bson.Marshal(frame)
	if err != nil {
		return fmt.Errorf("converting frame frameData for frame %v to BSON: %v", frame.GetPageId(), err)
	}

	res, err := collection.ReplaceOne(ctx, filter, frameData)
	if err != nil {
		// error detected
		return err
	}

	// if we matched 0 items then insert
	if res.MatchedCount == 0 {
		res, err := collection.InsertOne(ctx, frameData)
		if err != nil || res.InsertedID == nil {
			return fmt.Errorf("inserting frame %s: %v", frame.GetPageId(), err)
		}
	}

	return err
}
