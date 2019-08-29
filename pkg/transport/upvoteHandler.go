package transport

import (
	"context"
	"strings"

	pb "github.com/jackharley7/discussproto"
	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
)

func (c *UpvoteGrpcService) Upvote(ctx context.Context, req *pb.UpvoteRequest) (*pb.UpvoteResponse, error) {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return nil, err
	}

	cat := req.GetCategory()
	itemID := req.GetItemId()

	err = c.service.Upvote(userID, cat, itemID)
	if err != nil {
		return nil, err
	}

	err = c.service.PublishVote(cat, itemID)
	if err != nil {
		return nil, err
	}

	return &pb.UpvoteResponse{
		Done: true,
	}, nil
}

func (c *UpvoteGrpcService) Downvote(ctx context.Context, req *pb.DownvoteRequest) (*pb.DownvoteResponse, error) {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return nil, err
	}

	cat := req.GetCategory()
	itemID := req.GetItemId()

	err = c.service.Downvote(userID, cat, itemID)
	if err != nil {
		return nil, err
	}

	err = c.service.PublishVote(cat, itemID)
	if err != nil {
		return nil, err
	}

	return &pb.DownvoteResponse{
		Done: true,
	}, nil
}

func (c *UpvoteGrpcService) RemoveVote(ctx context.Context, req *pb.RemoveVoteRequest) (*pb.RemoveVoteResponse, error) {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return nil, err
	}

	cat := req.GetCategory()
	itemID := req.GetItemId()

	err = c.service.RemoveVote(userID, cat, itemID)
	if err != nil {
		return nil, err
	}

	err = c.service.PublishVote(cat, itemID)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveVoteResponse{
		Done: true,
	}, nil
}

func (c *UpvoteGrpcService) CheckUserVotes(ctx context.Context, req *pb.CheckUserVotesRequest) (*pb.CheckUserVotesResponse, error) {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return nil, err
	}

	itemIds := strings.Split(req.GetItemIds(), ",")

	checks, err := c.service.CheckUserVotes(userID, req.GetCategory(), itemIds)
	if err != nil {
		return nil, err
	}

	return &pb.CheckUserVotesResponse{
		ItemIds:  req.GetItemIds(),
		Category: req.GetCategory(),
		Votes:    mapUserVoteChecksOut(checks),
	}, nil
}

func mapUserVoteChecksOut(ucs []upvote.Check) []*pb.UpvoteCheck {
	res := []*pb.UpvoteCheck{}
	for _, uc := range ucs {
		res = append(res, mapUserVoteCheckOut(uc))
	}
	return res
}

func mapUserVoteCheckOut(uc upvote.Check) *pb.UpvoteCheck {
	return &pb.UpvoteCheck{
		ItemId: uc.ItemID,
		Type:   pb.UpvoteType(uc.Type),
	}
}
