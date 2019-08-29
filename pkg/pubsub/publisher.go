package pubsub

import (
	"cloud.google.com/go/pubsub"
	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
)

type topics struct {
	voted *pubsub.Topic
}

func NewPublisher(v *pubsub.Topic) upvote.Publisher {
	return &topics{
		voted: v,
	}
}
