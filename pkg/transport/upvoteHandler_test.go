package transport

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackharley7/golang-upvote-microservice/internal/upvote"
	mock_upvote "github.com/jackharley7/golang-upvote-microservice/internal/upvote/mock"

	pb "github.com/jackharley7/discussproto"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpvoteGrpcService_Upvote(t *testing.T) {
	client, db, _ := initUpvoteTestSetup(t)

	type fields struct {
		service *upvote.Service
	}

	type args struct {
		ctx context.Context
		req *pb.UpvoteRequest
	}

	type results struct {
		upvote      *upvote.Upvote
		upvoteCount *upvote.UpvoteCount
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *results
		wantErr error
	}{
		{
			name: "Upvote - create new upvote & inc upvoteCount",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "somecategory", "1001", int64(1)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(1)),
				req: &pb.UpvoteRequest{
					ItemId:   "1001",
					Category: "somecategory",
				},
			},
			want: &results{
				upvote: &upvote.Upvote{
					CatItemID: "somecategory.1001",
					UserID:    1,
					Type:      upvote.UP,
				},
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "somecategory.1001",
					Count:     1,
				},
			},
			wantErr: nil,
		},
		{
			name: "Upvote - upvote already exists, remove upvote and decrement",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "cat", "1", int64(0)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(12)),
				req: &pb.UpvoteRequest{
					ItemId:   "1",
					Category: "cat",
				},
			},
			want: &results{
				upvote: nil,
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "cat.1",
					Count:     0,
				},
			},
			wantErr: nil,
		},
		{
			name: "Upvote - downvote already exists, change to upvote and increment by 2",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "cat", "11", int64(1)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(12)),
				req: &pb.UpvoteRequest{
					ItemId:   "11",
					Category: "cat",
				},
			},
			want: &results{
				// upvote: &upvote.Upvote{
				// 	ID:     data.s[11].ObjectID,
				// 	UserID: 12,
				// 	Type:   upvote.UP,
				// },
				upvote: &upvote.Upvote{
					CatItemID: "cat.11",
					UserID:    12,
					Type:      upvote.UP,
				},
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "cat.11",
					Count:     1,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UpvoteGrpcService{
				service: tt.fields.service,
			}
			_, err := c.Upvote(tt.args.ctx, tt.args.req)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpvoteGrpcService.Upvote() error = %v, wantErr %v", err.Error(), tt.wantErr.Error())
				return
			}
			if tt.want != nil {
				ctx := context.Background()

				if tt.want.upvote != nil {
					var result upvote.Upvote
					err := db.Upvotes.FindOne(ctx, bson.M{"cat_item_id": tt.want.upvote.CatItemID, "user_id": tt.want.upvote.UserID}).Decode(&result)
					if tt.want.upvote.Type != 0 && err != nil {
						t.Errorf("UpvoteGrpcService.Upvote() error finding upvote %v", err)
					}
					if tt.want.upvote.Type == 0 && err == nil {
						t.Errorf("UpvoteGrpcService.Upvote() should be no upvote, instead got %v", result)
					}
					if result.Type != tt.want.upvote.Type {
						t.Errorf("UpvoteGrpcService.Upvote() upvote type not correct")
					}
				}

				if tt.want.upvoteCount != nil {
					var upvoteResult upvote.UpvoteCount
					err := db.UpvoteCounts.FindOne(ctx, bson.M{"cat_item_id": tt.want.upvoteCount.CatItemID}).Decode(&upvoteResult)
					if err != nil {
						t.Errorf("UpvoteGrpcService.Upvote() error finding comment %v", err)
					}
					if upvoteResult.Count != tt.want.upvoteCount.Count {
						t.Errorf("UpvoteGrpcService.Upvote() upvote not created")
					}
				}
			}
		})
	}
}

func TestUpvoteGrpcService_Downvote(t *testing.T) {
	client, db, _ := initUpvoteTestSetup(t)

	type fields struct {
		service *upvote.Service
	}

	type args struct {
		ctx context.Context
		req *pb.DownvoteRequest
	}

	type results struct {
		upvote      *upvote.Upvote
		upvoteCount *upvote.UpvoteCount
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *results
		wantErr error
	}{
		{
			name: "Downvote - create new downvote & decrement upvoteCount",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "somecategory", "1001", int64(-1)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(1)),
				req: &pb.DownvoteRequest{
					ItemId:   "1001",
					Category: "somecategory",
				},
			},
			want: &results{
				upvote: &upvote.Upvote{
					CatItemID: "somecategory.1001",
					UserID:    1,
					Type:      upvote.DOWN,
				},
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "somecategory.1001",
					Count:     -1,
				},
			},
			wantErr: nil,
		},
		{
			name: "Downvote - downvote already exists, remove downvote and increment",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "cat", "11", int64(0)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(12)),
				req: &pb.DownvoteRequest{
					ItemId:   "11",
					Category: "cat",
				},
			},
			want: &results{
				upvote: nil,
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "cat.11",
					Count:     0,
				},
			},
			wantErr: nil,
		},
		{
			name: "Downvote - upvote already exists, change to downvote and decrement by 2",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "cat", "1", int64(-1)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(12)),
				req: &pb.DownvoteRequest{
					ItemId:   "1",
					Category: "cat",
				},
			},
			want: &results{
				upvote: &upvote.Upvote{
					CatItemID: "cat.1",
					UserID:    12,
					Type:      upvote.DOWN,
				},
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "cat.1",
					Count:     -1,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UpvoteGrpcService{
				service: tt.fields.service,
			}
			_, err := c.Downvote(tt.args.ctx, tt.args.req)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpvoteGrpcService.Upvote() error = %v, wantErr %v", err.Error(), tt.wantErr.Error())
				return
			}
			if tt.want != nil {
				ctx := context.Background()

				if tt.want.upvote != nil {
					var result upvote.Upvote
					err := db.Upvotes.FindOne(ctx, bson.M{"cat_item_id": tt.want.upvote.CatItemID, "user_id": tt.want.upvote.UserID}).Decode(&result)
					if tt.want.upvote.Type != 0 && err != nil {
						t.Errorf("UpvoteGrpcService.Upvote() error finding upvote %v", err)
					}
					if tt.want.upvote.Type == 0 && err == nil {
						t.Errorf("UpvoteGrpcService.Upvote() should be no upvote, instead got %v", result)
					}
					if result.Type != tt.want.upvote.Type {
						t.Errorf("UpvoteGrpcService.Upvote() upvote type not correct")
					}
				}

				if tt.want.upvoteCount != nil {
					var upvoteResult upvote.UpvoteCount
					err := db.UpvoteCounts.FindOne(ctx, bson.M{"cat_item_id": tt.want.upvoteCount.CatItemID}).Decode(&upvoteResult)
					if err != nil {
						t.Errorf("UpvoteGrpcService.Upvote() error finding comment %v", err)
					}
					if upvoteResult.Count != tt.want.upvoteCount.Count {
						t.Errorf("UpvoteGrpcService.Upvote() upvote not created")
					}
				}
			}
		})
	}
}

func TestUpvoteGrpcService_RemoveVote(t *testing.T) {
	client, db, _ := initUpvoteTestSetup(t)
	// _, services := initUpvoteGrpcService(db, client, Mocks{})

	type fields struct {
		service *upvote.Service
	}

	type args struct {
		ctx context.Context
		req *pb.RemoveVoteRequest
	}

	type results struct {
		upvote      *upvote.Upvote
		upvoteCount *upvote.UpvoteCount
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *results
		wantErr error
	}{
		{
			name: "Removevote - remove a upvote & increment",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "cat", "1", int64(0)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(12)),
				req: &pb.RemoveVoteRequest{
					ItemId:   "1",
					Category: "cat",
				},
			},
			want: &results{
				upvote: nil,
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "cat.1",
					Count:     0,
				},
			},
			wantErr: nil,
		},
		{
			name: "RemoveVote - remove a downvote & deccrement comment",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "cat", "11", int64(0)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(12)),
				req: &pb.RemoveVoteRequest{
					ItemId:   "11",
					Category: "cat",
				},
			},
			want: &results{
				upvote: nil,
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "cat.11",
					Count:     0,
				},
			},
			wantErr: nil,
		},
		{
			name: "RemoveVote - no upvote, do nothing",
			fields: fields{
				service: func() *upvote.Service {
					mocks := Mocks{
						Publisher: new(mock_upvote.Publisher),
					}
					mocks.Publisher.On("UserUpvoted", "cat", "3", int64(1)).Return(nil)
					_, services := initUpvoteGrpcService(db, client, mocks)
					return services.upvoteService
				}(),
			},
			args: args{
				ctx: createContext(int64(10)),
				req: &pb.RemoveVoteRequest{
					ItemId:   "3",
					Category: "cat",
				},
			},
			want: &results{
				upvote: nil,
				upvoteCount: &upvote.UpvoteCount{
					CatItemID: "cat.3",
					Count:     1,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UpvoteGrpcService{
				service: tt.fields.service,
			}
			_, err := c.RemoveVote(tt.args.ctx, tt.args.req)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpvoteGrpcService.Upvote() error = %v, wantErr %v", err.Error(), tt.wantErr.Error())
				return
			}
			if tt.want != nil {
				ctx := context.Background()

				if tt.want.upvote != nil {
					var result upvote.Upvote
					err := db.Upvotes.FindOne(ctx, bson.M{"cat_item_id": tt.want.upvote.CatItemID, "user_id": tt.want.upvote.UserID}).Decode(&result)
					if tt.want.upvote.Type != 0 && err != nil {
						t.Errorf("UpvoteGrpcService.Upvote() error finding upvote %v", err)
					}
					if tt.want.upvote.Type == 0 && err == nil {
						t.Errorf("UpvoteGrpcService.Upvote() should be no upvote, instead got %v", result)
					}
					if result.Type != tt.want.upvote.Type {
						t.Errorf("UpvoteGrpcService.Upvote() upvote type not correct")
					}
				}

				if tt.want.upvoteCount != nil {
					var upvoteResult upvote.UpvoteCount
					err := db.UpvoteCounts.FindOne(ctx, bson.M{"cat_item_id": tt.want.upvoteCount.CatItemID}).Decode(&upvoteResult)
					if err != nil {
						t.Errorf("UpvoteGrpcService.Upvote() error finding comment %v", err)
					}
					if upvoteResult.Count != tt.want.upvoteCount.Count {
						t.Errorf("UpvoteGrpcService.Upvote() upvote not created")
					}
				}
			}
		})
	}
}

func TestUpvoteGrpcService_CheckUserVotes(t *testing.T) {
	client, db, _ := initUpvoteTestSetup(t)
	_, services := initUpvoteGrpcService(db, client, Mocks{})

	type fields struct {
		service *upvote.Service
	}
	type args struct {
		ctx context.Context
		req *pb.CheckUserVotesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CheckUserVotesResponse
		wantErr bool
	}{
		{
			name: "CheckUserVotes - items do not exist",
			fields: fields{
				service: services.upvoteService,
			},
			args: args{
				ctx: createContext(int64(10)),
				req: &pb.CheckUserVotesRequest{
					Category: "doesnotexist",
					ItemIds:  "1,2",
				},
			},
			want: &pb.CheckUserVotesResponse{
				Category: "doesnotexist",
				ItemIds:  "1,2",
				Votes: []*pb.UpvoteCheck{
					&pb.UpvoteCheck{
						ItemId: "1",
						Type:   pb.UpvoteType_NONE,
					},
					&pb.UpvoteCheck{
						ItemId: "2",
						Type:   pb.UpvoteType_NONE,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CheckUserVotes - returns none if user not voted",
			fields: fields{
				service: services.upvoteService,
			},
			args: args{
				ctx: createContext(int64(10)),
				req: &pb.CheckUserVotesRequest{
					Category: "cat",
					ItemIds:  "1,12,32",
				},
			},
			want: &pb.CheckUserVotesResponse{
				Category: "cat",
				ItemIds:  "1,12,32",
				Votes: []*pb.UpvoteCheck{
					&pb.UpvoteCheck{
						ItemId: "1",
						Type:   pb.UpvoteType_NONE,
					},
					&pb.UpvoteCheck{
						ItemId: "12",
						Type:   pb.UpvoteType_NONE,
					},
					&pb.UpvoteCheck{
						ItemId: "32",
						Type:   pb.UpvoteType_NONE,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CheckUserVotes - returns upvoted and downvoted",
			fields: fields{
				service: services.upvoteService,
			},
			args: args{
				ctx: createContext(int64(12)),
				req: &pb.CheckUserVotesRequest{
					Category: "cat",
					ItemIds:  "1,12,32",
				},
			},
			want: &pb.CheckUserVotesResponse{
				Category: "cat",
				ItemIds:  "1,12,32",
				Votes: []*pb.UpvoteCheck{
					&pb.UpvoteCheck{
						ItemId: "1",
						Type:   pb.UpvoteType_UP,
					},
					&pb.UpvoteCheck{
						ItemId: "12",
						Type:   pb.UpvoteType_DOWN,
					},
					&pb.UpvoteCheck{
						ItemId: "32",
						Type:   pb.UpvoteType_NONE,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UpvoteGrpcService{
				service: tt.fields.service,
			}
			got, err := c.CheckUserVotes(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpvoteGrpcService.CheckUserVotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpvoteGrpcService.CheckUserVotes() = %v, want %v", got, tt.want)
			}
		})
	}
}
