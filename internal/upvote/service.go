package upvote

type Repository interface {
	GetVote(userID int64, catItemID string) (*Upvote, error)
	Upvote(userID int64, catItemID string, increment int) error
	Downvote(userID int64, catItemID string, increment int) error
	RemoveVote(userID int64, catItemID string, increment int) error
	GetVotes(userID int64, itemIDs []string) (map[string]UpvoteType, error)
	GetVoteCount(catItemID string) (int64, error)
}

type Publisher interface {
	UserUpvoted(cat string, itemID string, votes int64) error
}

type Service struct {
	repository Repository
	publisher  Publisher
}

func NewService(r Repository, p Publisher) *Service {
	return &Service{
		repository: r,
		publisher:  p,
	}
}
