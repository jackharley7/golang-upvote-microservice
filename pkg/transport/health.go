package transport

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type HealthGrpcService struct {
	db *mongo.Client
}

func NewHealthGrpcService(db *mongo.Client) *HealthGrpcService {
	return &HealthGrpcService{
		db: db,
	}
}

func (r *HealthGrpcService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {

	healthCtx := context.Background()
	dbErr := r.db.Ping(healthCtx, nil)
	fmt.Println("error connecting to db", dbErr)

	if dbErr != nil {
		return &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	if req.Service == "" {
		// check the server overall health status.
		return &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_SERVING,
		}, nil
	}

	return nil, status.Error(codes.NotFound, "unknown service")
}

func (r *HealthGrpcService) Watch(*pb.HealthCheckRequest, pb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

var _ pb.HealthServer = &HealthGrpcService{}
