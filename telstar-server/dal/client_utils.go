package dal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
	"time"
)

func getCollectionNames(connectionUrl string) (pNames, sNames []string, err error) {

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

	names, err := client.Database(DBNAME).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, nil, fmt.Errorf("listing collection names: %v", err)
	}

	p, err := regexp.Compile(REGEXP)
	if err != nil {
		return nil, nil, fmt.Errorf("compiling regex expression %s: %v", REGEXP, err)
	}
	s, err := regexp.Compile(REGEXS)

	if err != nil {
		return nil, nil, fmt.Errorf("compiling regex expression %s: %v", REGEXS, err)
	}

	for i := range names {

		if s.MatchString(names[i]) {
			sNames = append(sNames, names[i])
		}
		if p.MatchString(names[i]) {
			pNames = append(pNames, names[i])
		}
	}
	return pNames, sNames, err
}

func createIndex(connectionUrl string, collectionName string) error {

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

	collection := client.Database(DBNAME).Collection(collectionName)

	var indexModel = mongo.IndexModel{
		Keys: bson.D{
			{"pid", 1},
			//{"PID.frame-id", -1},
		},
	}
	name, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Println("Name of Index Created: " + name)

	return nil
}
