package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

func CreateSubscription(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, subName string) (*pubsub.Subscription, error) {
	subscription := client.Subscription(subName)
	exists, err := subscription.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		if _, err = client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{Topic: topic}); err != nil {
			return nil, err
		}
	}
	return subscription, nil
}
