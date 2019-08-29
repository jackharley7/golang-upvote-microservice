package mongoservice

import (
	"context"
	"time"

	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *repository) GetVoteCount(catItemID string) (int64, error) {
	var c upvote.UpvoteCount
	err := r.upvoteCount.FindOne(context.Background(), bson.M{"cat_item_id": catItemID}).Decode(&c)
	if err != nil {
		return 0, err
	}
	return c.Count, nil
}

func (r *repository) GetVote(userID int64, catItemID string) (*upvote.Upvote, error) {
	var c upvote.Upvote
	err := r.upvotes.FindOne(context.Background(), bson.M{"cat_item_id": catItemID, "user_id": userID}).Decode(&c)
	if err != nil {
		return nil, nil
	}
	return &c, nil
}

func (r *repository) Upvote(userID int64, catItemID string, increment int) error {
	// start a transaction
	ctx := context.Background()
	return r.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {

		// get itemCount
		// if no itemCount, create one
		var uc upvote.UpvoteCount
		filter := bson.D{{"cat_item_id", catItemID}}
		if err := r.upvoteCount.FindOne(ctx, filter).Decode(&uc); err != nil && err.Error() != "mongo: no documents in result" {
			return err
		}

		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		if uc.CatItemID == "" {
			update := bson.M{"$set": bson.M{"cat_item_id": catItemID, "count": 0}}
			_, err := r.upvoteCount.UpdateOne(sessionContext, bson.M{"cat_item_id": catItemID}, update, options.Update().SetUpsert(true))
			if err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}
		}

		// increment comments upvotes
		_, err := r.upvoteCount.UpdateOne(sessionContext, bson.D{{"cat_item_id", catItemID}}, bson.D{{"$inc", bson.D{{"count", increment}}}})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		update := bson.M{"$set": bson.M{"cat_item_id": catItemID, "user_id": userID, "created_at": time.Now(), "type": upvote.UP}}
		_, err = r.upvotes.UpdateOne(sessionContext, bson.M{"cat_item_id": catItemID, "user_id": userID}, update, options.Update().SetUpsert(true))

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		} else {
			sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})
}

func (r *repository) Downvote(userID int64, catItemID string, increment int) error {
	ctx := context.Background()
	return r.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var uc upvote.UpvoteCount
		filter := bson.D{{"cat_item_id", catItemID}}
		if err := r.upvoteCount.FindOne(ctx, filter).Decode(&uc); err != nil && err.Error() != "mongo: no documents in result" {
			return err
		}

		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		if uc.CatItemID == "" {
			update := bson.M{"$set": bson.M{"cat_item_id": catItemID, "count": 0}}
			_, err := r.upvoteCount.UpdateOne(sessionContext, bson.M{"cat_item_id": catItemID}, update, options.Update().SetUpsert(true))
			if err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}
		}

		// increment comments upvotes
		_, err := r.upvoteCount.UpdateOne(sessionContext, bson.D{{"cat_item_id", catItemID}}, bson.D{{"$inc", bson.D{{"count", increment}}}})
		if err != nil {
			return err
		}

		update := bson.M{"$set": bson.M{"cat_item_id": catItemID, "user_id": userID, "created_at": time.Now(), "type": upvote.DOWN}}
		_, err = r.upvotes.UpdateOne(sessionContext, bson.M{"cat_item_id": catItemID, "user_id": userID}, update, options.Update().SetUpsert(true))
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		} else {
			sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})
}

func (r *repository) RemoveVote(userID int64, catItemID string, increment int) error {
	ctx := context.Background()
	return r.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		// increment comments upvotes
		_, err = r.upvoteCount.UpdateOne(sessionContext, bson.D{{"cat_item_id", catItemID}}, bson.D{{"$inc", bson.D{{"count", increment}}}})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		_, err = r.upvotes.DeleteOne(sessionContext, bson.M{"cat_item_id": catItemID, "user_id": userID})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		} else {
			sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})
}

func (r *repository) GetVotes(userID int64, catItemIDs []string) (map[string]upvote.UpvoteType, error) {

	findOptions := options.Find()
	limit := len(catItemIDs)
	findOptions.SetLimit(int64(limit))

	results := make(map[string]upvote.UpvoteType)

	bsonCommentArray := bson.A{}

	for _, cID := range catItemIDs {
		bsonCommentArray = append(bsonCommentArray, cID)
	}

	query := bson.D{
		{"user_id", userID},
		{"cat_item_id", bson.D{{"$in", bsonCommentArray}}},
	}

	ctx := context.Background()

	cur, err := r.upvotes.Find(context.TODO(), query, findOptions)
	if err != nil {
		return results, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var item upvote.Upvote
		err := cur.Decode(&item)
		if err != nil {
			return results, err
		}
		results[item.CatItemID] = item.Type
	}

	if err := cur.Err(); err != nil {
		return results, err
	}

	return results, nil

}
