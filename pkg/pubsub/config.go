package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

func ConfigurePubsub(ctx context.Context, projectID string) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return client, nil
}
