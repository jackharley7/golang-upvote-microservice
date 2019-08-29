package mongoservice

import (
	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	client      *mongo.Client
	upvotes     *mongo.Collection
	upvoteCount *mongo.Collection
}

func NewUpvoteRepository(client *mongo.Client, upvotes *mongo.Collection, upvoteCount *mongo.Collection) upvote.Repository {
	return &repository{
		client,
		upvotes,
		upvoteCount,
	}
}
