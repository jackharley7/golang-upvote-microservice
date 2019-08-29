package transport

import (
	"context"
	"strconv"

	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UpvoteGrpcService struct {
	service *upvote.Service
}

func NewUpvoteGrpcService(s *upvote.Service) *UpvoteGrpcService {
	return &UpvoteGrpcService{
		service: s,
	}
}

func getUserIDFromToken(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, nil
	}
	userID, err := strconv.ParseInt(md.Get("authuserid")[0], 10, 64)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, "Failed to convert userId to int64")
	}
	return userID, nil
}
