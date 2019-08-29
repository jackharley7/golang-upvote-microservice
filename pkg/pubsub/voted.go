package pubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

type VotedMessage struct {
	Cat    string
	ItemID string
	Votes  int64
}

func (t *topics) UserUpvoted(cat string, itemID string, votes int64) error {
	ctx := context.Background()

	voted := VotedMessage{
		Cat:    cat,
		ItemID: itemID,
		Votes:  votes,
	}

	n, err := json.Marshal(voted)
	if err != nil {
		return err
	}
	_, err = t.voted.Publish(ctx, &pubsub.Message{Data: n}).Get(ctx)
	return nil
}
