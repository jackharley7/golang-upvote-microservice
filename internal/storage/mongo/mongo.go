package mongoservice

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type key string

const (
	MongoConnectionUrlKey          = key("mongoConnectionUrlKey")
	MongoDBName                    = key("mongoDBName")
	UpvoteMongoCollectionName      = key("upvoteMongoCollectionName")
	UpvoteCountMongoCollectionName = key("upvoteCountMongoCollectionName")
)

type Collections struct {
	Upvotes      *mongo.Collection
	UpvoteCounts *mongo.Collection
}

type Index struct {
	Key   string
	Value int
}

func Connect(ctx context.Context) (*mongo.Client, *Collections, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(ctx.Value(MongoConnectionUrlKey).(string)))
	if err != nil {
		return nil, nil, fmt.Errorf("todo: couldn't connect to mongo: %v", err)
	}

	err = client.Connect(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("todo: mongo client couldn't connect with background context: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("todo: couldn't connect to mongo: %v", err)
	}

	upvoteDB := client.Database(ctx.Value(MongoDBName).(string)).Collection(ctx.Value(UpvoteMongoCollectionName).(string))
	upvoteCountDB := client.Database(ctx.Value(MongoDBName).(string)).Collection(ctx.Value(UpvoteCountMongoCollectionName).(string))

	return client, &Collections{
		Upvotes:      upvoteDB,
		UpvoteCounts: upvoteCountDB,
	}, nil
}

func newContext() (context.Context, context.CancelFunc) {
	ctx, _ := context.WithTimeout(context.Background(), 85*time.Second)
	return context.WithCancel(ctx)
}
