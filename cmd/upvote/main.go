package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	mongoservice "github.com/jackharley7/golang-upvote-microservice/internal/storage/mongo"
	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
	"github.com/jackharley7/golang-upvote-microservice/pkg/pubsub"
	"github.com/jackharley7/golang-upvote-microservice/pkg/transport"

	"github.com/jackharley7/discussproto"
	"google.golang.org/grpc"

	"github.com/urfave/cli"
	health_proto "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "upvoteservice"
	app.Usage = "upvote service"
	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "start server",
			Action: runStart(),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "pubsub-project-id",
					Usage:  "pubsub-project-id",
					EnvVar: "PUBSUB-PROJECT-ID",
				},
				cli.StringFlag{
					Name:   "pubsub-upvote-topic",
					Usage:  "pubsub-upvote-topic",
					EnvVar: "PUBSUB-UPVOTE-TOPIC",
				},
				cli.StringFlag{
					Name:   "pubsub-upvote-subscription",
					Usage:  "pubsub-upvote-subscription",
					EnvVar: "PUBSUB-UPVOTE-SUBSCRIPTION",
				},
				cli.StringFlag{
					Name:   "mongo-connection-url",
					Usage:  "mongo-connection-url",
					EnvVar: "MONGO_CONNECTION_URL",
				},
				cli.StringFlag{
					Name:   "mongo-db-name",
					Usage:  "mongo-db-name",
					EnvVar: "MONGO_DB_NAME",
				},
				cli.StringFlag{
					Name:   "mongo-collection-upvote",
					Usage:  "mongo-collection-upvote",
					EnvVar: "MONGO_COLLECTION_UPVOTE",
				},
				cli.StringFlag{
					Name:   "mongo-collection-count",
					Usage:  "mongo-collection-count",
					EnvVar: "MONGO_COLLECTION_COUNT",
				},
				cli.IntFlag{
					Name:   "server-port",
					Usage:  "server port",
					EnvVar: "SERVER_PORT",
					Value:  9000,
				},
				cli.StringFlag{
					Name:   "service-env",
					Usage:  "service-env",
					EnvVar: "SERVICE_ENV",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Failed to start app")
		fmt.Print(err)
		os.Exit(1)
	}
}

func runStart() func(*cli.Context) error {
	return func(c *cli.Context) error {
		fmt.Println("Upvote Service Started")
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		// Setup Server
		// Connect to MongoDB
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)

		defer cancel()
		ctx = context.WithValue(ctx, mongoservice.MongoConnectionUrlKey, os.Getenv("MONGO_CONNECTION_URL"))
		ctx = context.WithValue(ctx, mongoservice.MongoDBName, os.Getenv("MONGO_DB_NAME"))
		ctx = context.WithValue(ctx, mongoservice.UpvoteMongoCollectionName, os.Getenv("MONGO_COLLECTION_UPVOTE"))
		ctx = context.WithValue(ctx, mongoservice.UpvoteCountMongoCollectionName, os.Getenv("MONGO_COLLECTION_COUNT"))

		mongoclient, collections, err := mongoservice.Connect(ctx)
		if err != nil {
			log.Fatalf("Failed to connect to mongoDB: %v", err)
		}

		// Init Publisher
		publisherCtx := context.Background()
		publisherCtx, cancelPublisher := context.WithCancel(publisherCtx)
		defer cancelPublisher()

		client, err := pubsub.ConfigurePubsub(publisherCtx, os.Getenv("PUBSUB-PROJECT-ID"))
		if err != nil {
			log.Fatalf("fatal pubsub.ConfigurePubsub: %v", err)
		}

		upvoteTopic, err := pubsub.CreateTopic(publisherCtx, client, os.Getenv("PUBSUB-UPVOTE-TOPIC"))
		if err != nil {
			log.Fatalf("fatal pubsub.CreateTopic %v", err)
		}
		publisher := pubsub.NewPublisher(upvoteTopic)

		// Init repository
		upvoteRepository := mongoservice.NewUpvoteRepository(mongoclient, collections.Upvotes, collections.Upvotes)

		// Init service
		upvoteService := upvote.NewService(upvoteRepository, publisher)

		s := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
			)),
		)

		// init health grpc service for kubenetes health checks
		healthGrpcService := transport.NewHealthGrpcService(mongoclient)

		// init upvote grpc service
		grpcUpvoteService := transport.NewUpvoteGrpcService(upvoteService)

		health_proto.RegisterHealthServer(s, healthGrpcService)
		discussproto.RegisterUpvoteServiceServer(s, grpcUpvoteService)

		// Start Server
		return s.Serve(lis)
	}
}
