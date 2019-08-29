package transport

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"

	mongoservice "github.com/jackharley7/golang-upvote-microservice/internal/storage/mongo"
	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
	mock_upvote "github.com/jackharley7/golang-upvote-microservice/internal/upvote/mock"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc/metadata"
)

type DbConfig struct {
	ConnectionUrl             string `default:"mongodb://localhost:30001/test"`
	DbName                    string `default:"testUpvote"`
	UpvoteCollectionName      string `default:"upvotes"`
	UpvoteCountCollectionName string `default:"upvoteCounts"`
}

type Data struct {
	Upvotes      []upvote.Upvote
	UpvoteCounts []upvote.UpvoteCount
}

type Mocks struct {
	Publisher *mock_upvote.Publisher
}

type Services struct {
	upvoteService *upvote.Service
}

// InitDB connect, clear and seed DB
func initDB(ctx context.Context, t *testing.T) (*mongo.Client, *mongoservice.Collections, *Data) {
	client, collections, data, err := createMongoDbConnection(ctx)
	if err != nil {
		t.Errorf("Unable to connect to db. %v", err)
	}
	return client, collections, data
}

func initUpvoteGrpcService(collections *mongoservice.Collections, client *mongo.Client, mocks Mocks) (*UpvoteGrpcService, Services) {

	if mocks.Publisher == nil {
		mocks.Publisher = new(mock_upvote.Publisher)
	}

	upvoteRepository := mongoservice.NewUpvoteRepository(client, collections.Upvotes, collections.UpvoteCounts)

	upvoteService := upvote.NewService(upvoteRepository, mocks.Publisher)

	svs := Services{
		upvoteService: upvoteService,
	}

	s := NewUpvoteGrpcService(upvoteService)

	return s, svs
}

func initUpvoteTestSetup(t *testing.T) (*mongo.Client, *mongoservice.Collections, Data) {

	var c DbConfig
	_ = envconfig.Process("db", &c)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, mongoservice.MongoConnectionUrlKey, c.ConnectionUrl)
	ctx = context.WithValue(ctx, mongoservice.MongoDBName, c.DbName)
	ctx = context.WithValue(ctx, mongoservice.UpvoteMongoCollectionName, c.UpvoteCollectionName)
	ctx = context.WithValue(ctx, mongoservice.UpvoteCountMongoCollectionName, c.UpvoteCountCollectionName)

	client, collections, data := initDB(ctx, t)

	return client, collections, *data
}

func createMongoDbConnection(ctx context.Context) (*mongo.Client, *mongoservice.Collections, *Data, error) {
	client, collections, err := mongoservice.Connect(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	data, err := seedDB(ctx, collections.Upvotes, collections.UpvoteCounts)
	if err != nil {
		return nil, nil, nil, err
	}
	return client, collections, data, nil
}

func seedDB(ctx context.Context, upvoteC *mongo.Collection, upvoteCountC *mongo.Collection) (*Data, error) {
	upvoteC.DeleteMany(ctx, bson.D{{}})
	upvoteCountC.DeleteMany(ctx, bson.D{{}})

	savedUpvotes, err := seedUpvotes(ctx, upvoteC)
	savedUpvoteCounts, err := seedUpvoteCounts(ctx, upvoteCountC)
	if err != nil {
		return nil, err
	}

	data := &Data{
		Upvotes:      savedUpvotes,
		UpvoteCounts: savedUpvoteCounts,
	}
	return data, nil
}

func createContext(userID int64) context.Context {
	ctx := context.Background()

	// add key-value pairs of metadata to context
	ctx = metadata.NewIncomingContext(
		ctx,
		metadata.Pairs("authuserid", strconv.Itoa(int(userID))),
	)
	return ctx
}

func seedUpvotes(ctx context.Context, db *mongo.Collection) ([]upvote.Upvote, error) {
	var ups []upvote.Upvote

	for i := 0; i < 10; i++ {
		u := upvote.Upvote{
			UserID:    12,
			CatItemID: "cat." + strconv.Itoa(i),
			Type:      upvote.UP,
		}
		ups = append(ups, u)
	}

	for i := 10; i < 20; i++ {
		u := upvote.Upvote{
			UserID:    12,
			CatItemID: "cat." + strconv.Itoa(i),
			Type:      upvote.DOWN,
		}
		ups = append(ups, u)
	}

	var ui []interface{}
	for _, u := range ups {
		ui = append(ui, u)
	}
	_, err := db.InsertMany(ctx, ui)
	if err != nil {
		return nil, fmt.Errorf("create: couldn't be created: %v", err)
	}

	var savedUpvotes []upvote.Upvote
	cursor, _ := db.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var t upvote.Upvote
		cursor.Decode(&t)
		savedUpvotes = append(savedUpvotes, t)
	}

	return savedUpvotes, nil
}

func seedUpvoteCounts(ctx context.Context, db *mongo.Collection) ([]upvote.UpvoteCount, error) {
	var ups []upvote.UpvoteCount

	for i := 0; i < 10; i++ {
		u := upvote.UpvoteCount{
			CatItemID: "cat." + strconv.Itoa(i),
			Count:     1,
		}
		ups = append(ups, u)
	}

	for i := 10; i < 20; i++ {
		u := upvote.UpvoteCount{
			CatItemID: "cat." + strconv.Itoa(i),
			Count:     -1,
		}
		ups = append(ups, u)
	}

	var ui []interface{}
	for _, u := range ups {
		ui = append(ui, u)
	}
	_, err := db.InsertMany(ctx, ui)
	if err != nil {
		return nil, fmt.Errorf("create: couldn't be created: %v", err)
	}

	var savedUpvoteCounts []upvote.UpvoteCount
	cursor, _ := db.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var t upvote.UpvoteCount
		cursor.Decode(&t)
		savedUpvoteCounts = append(savedUpvoteCounts, t)
	}

	return savedUpvoteCounts, nil
}
