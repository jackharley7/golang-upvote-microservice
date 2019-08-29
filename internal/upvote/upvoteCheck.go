package upvote

type Check struct {
	ItemID string
	Type   UpvoteType
}

func (s *Service) CheckUserVotes(userID int64, cat string, itemIDs []string) ([]Check, error) {
	// check if user has voted

	checks := []Check{}
	catItemIDs := []string{}

	for _, id := range itemIDs {
		catItemIDs = append(catItemIDs, buildId(cat, id))
	}

	votes, err := s.repository.GetVotes(userID, catItemIDs)
	if err != nil {
		return checks, err
	}

	for _, id := range itemIDs {
		ut, ok := votes[buildId(cat, id)]
		if !ok {
			ut = NONE
		}
		checks = append(checks, Check{
			ItemID: id,
			Type:   ut,
		})
	}

	return checks, nil
}
