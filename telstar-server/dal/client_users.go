package dal

import (
	"context"
	"errors"
	"fmt"
	"github.com/johnnewcombe/telstar-library/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InsertOrReplaceUser(connectionUrl string, user types.User) error {

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	var (
		data []byte
	)

	if !ValidateUser(user) {
		return errors.New("userId or password does not meet the requirements")
	}
	if user.Password, err = HashPassword(user.Password); err != nil {
		return err
	}

	filter := bson.M{"user-id": user.UserId}
	collection := client.Database(DBNAME).Collection(AUTH_COLLECTION)
	// marshall the data
	if data, err = bson.Marshal(user); err != nil {
		return fmt.Errorf("converting user data for user id %v to BSON: %v", user.UserId, err)
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
			return fmt.Errorf("inserting user %s: %v", user.UserId, err)
		}
	}
	return err
}

func InsertOrReplaceUserByUser(connectionUrl string, newUser types.User, user types.User) error {

	// FIXME FIXME user param needs to be checked for permissions etc.
	//  i.e. admin and new users base page is in scope with current user.

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	var (
		data []byte
	)

	if !ValidateUser(newUser) {
		return errors.New("userId or password does not meet the requirements")
	}
	if newUser.Password, err = HashPassword(newUser.Password); err != nil {
		return err
	}

	filter := bson.M{"user-id": newUser.UserId}
	collection := client.Database(DBNAME).Collection(AUTH_COLLECTION)
	// marshall the data
	if data, err = bson.Marshal(newUser); err != nil {
		return fmt.Errorf("converting user data for user id %v to BSON: %v", newUser.UserId, err)
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
			return fmt.Errorf("inserting user %s: %v", newUser.UserId, err)
		}
	}
	return err
}

func DeleteUser(connectionUrl string, userId string) (int64, error) {

	// FIXME FIXME user param needs to be checked for permissions etc.
	//  i.e. admin and new users base page is in scope with current user.

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		return 0, err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"user-id": userId}
	collection := client.Database(DBNAME).Collection(AUTH_COLLECTION)

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, err
}

func DeleteUserByUser(connectionUrl string, userId string, user types.User) (int64, error) {

	// FIXME FIXME user param needs to be checked for permissions etc.
	//  i.e. admin and new users base page is in scope with current user.

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		return 0, err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"user-id": userId}
	collection := client.Database(DBNAME).Collection(AUTH_COLLECTION)

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, err
}

func GetUser(connectionUrl string, userId string) (types.User, error) {

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result types.User

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		return result, err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"user-id": userId}
	collection := client.Database(DBNAME).Collection(AUTH_COLLECTION)

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("finding user %s: %v", userId, err)
	}

	return result, nil
}

func GetUserByUser(connectionUrl string, userId string, user types.User) (types.User, error) {

	var result types.User

	// FIXME FIXME user param needs to be checked for permissions etc.
	//  i.e. admin and new users base page is in scope with current user.

	// get a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		return result, err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	filter := bson.M{"user-id": userId}
	collection := client.Database(DBNAME).Collection(AUTH_COLLECTION)

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("finding user %s: %v", userId, err)
	}

	return result, nil
}

func GetAllUsers(connectionUrl string, authUser string, pageNumber int) []types.User {
	return []types.User{}
}

func IsUserInScope(connection string, userId string, pageNumber int) bool {

	var (
		user types.User
		err  error
	)
	// get the user
	if user, err = GetUser(connection, userId); err != nil {
		return false
	}

	if user.IsInScope(pageNumber) {
		return true
	}
	return false
}

func IsUserAdmin(connection string, userId string) bool {

	var (
		user types.User
		err  error
	)
	// get the user
	if user, err = GetUser(connection, userId); err != nil {
		return false
	}

	return user.Admin
}

func IsUserGuest(connection string, userId string) bool {

	var (
		user types.User
		err  error
	)
	// get the user
	if user, err = GetUser(connection, userId); err != nil {
		return false
	}
	if user.IsGuest() {
		return true
	}
	return false
}
