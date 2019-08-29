package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

func CreateTopic(ctx context.Context, client *pubsub.Client, topicName string) (*pubsub.Topic, error) {
	topic := client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		if _, err := client.CreateTopic(ctx, topicName); err != nil {
			return nil, err
		}
	}
	return topic, nil
}
