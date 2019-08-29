package upvote

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpvoteType int

const (
	NONE UpvoteType = iota
	UP
	DOWN
)

type UpvoteCount struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CatItemID string             `bson:"cat_item_id,omitempty"`
	Count     int64
}

type Upvote struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CatItemID string             `bson:"cat_item_id,omitempty"`
	UserID    int64              `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
	Type      UpvoteType
}

func (s *Service) Upvote(userID int64, cat string, itemID string) error {
	// check if user has voted

	catItemID := buildId(cat, itemID)

	vote, err := s.repository.GetVote(userID, catItemID)
	if err != nil {
		return err
	}
	if vote == nil {
		// has not voted - create upvote, increment comments upvotes
		err := s.repository.Upvote(userID, catItemID, 1)
		return err
	}
	if vote.Type == DOWN {
		// has already downvoted - change downvote to upvote, increment comments by 2
		return s.repository.Upvote(userID, catItemID, 2)
	}

	// has already upvoted - remove and decrement
	if err = s.repository.RemoveVote(userID, catItemID, -1); err != nil {
		return err
	}

	// get votes
	return err
}

func (s *Service) Downvote(userID int64, cat string, itemID string) error {

	catItemID := buildId(cat, itemID)

	vote, err := s.repository.GetVote(userID, catItemID)
	if err != nil {
		return err
	}

	if vote == nil {
		// has not voted - create downvote, decrement comments upvotes
		return s.repository.Downvote(userID, catItemID, -1)
	}

	if vote.Type == UP {
		// has already upvoted - change upvote to downvote, decrement comments by 2
		return s.repository.Downvote(userID, catItemID, -2)
	}

	// has already downvoted - remove and increment
	return s.repository.RemoveVote(userID, catItemID, 1)
}

func (s *Service) RemoveVote(userID int64, cat string, itemID string) error {

	catItemID := buildId(cat, itemID)

	// check if user has voted
	vote, err := s.repository.GetVote(userID, catItemID)
	if err != nil {
		return err
	}

	if vote == nil {
		// has not voted - do nothing
		return nil
	}

	if vote.Type == UP {
		// has already upvoted - delete vote, decrement comments by 1
		return s.repository.RemoveVote(userID, catItemID, -1)
	}

	if vote.Type == DOWN {
		// has already downvoted - delete vote, increment comments by 1
		return s.repository.RemoveVote(userID, catItemID, 1)
	}

	// cleanup (should never be hit)
	return nil
}

func (s *Service) PublishVote(cat string, itemID string) error {
	count, err := s.repository.GetVoteCount(buildId(cat, itemID))
	if err != nil {
		return err
	}
	return s.publisher.UserUpvoted(cat, itemID, count)
}

func buildId(pre string, id string) string {
	return pre + "." + id
}
