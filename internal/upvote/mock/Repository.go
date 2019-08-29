// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import upvote "github.com/jackharley7/golang-upvote-microservice/internal/upvote"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Downvote provides a mock function with given fields: userID, catItemID, increment
func (_m *Repository) Downvote(userID int64, catItemID string, increment int) error {
	ret := _m.Called(userID, catItemID, increment)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, string, int) error); ok {
		r0 = rf(userID, catItemID, increment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetVote provides a mock function with given fields: userID, catItemID
func (_m *Repository) GetVote(userID int64, catItemID string) (*upvote.Upvote, error) {
	ret := _m.Called(userID, catItemID)

	var r0 *upvote.Upvote
	if rf, ok := ret.Get(0).(func(int64, string) *upvote.Upvote); ok {
		r0 = rf(userID, catItemID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*upvote.Upvote)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(userID, catItemID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVoteCount provides a mock function with given fields: catItemID
func (_m *Repository) GetVoteCount(catItemID string) (int64, error) {
	ret := _m.Called(catItemID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(catItemID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(catItemID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVotes provides a mock function with given fields: userID, itemIDs
func (_m *Repository) GetVotes(userID int64, itemIDs []string) (map[string]upvote.UpvoteType, error) {
	ret := _m.Called(userID, itemIDs)

	var r0 map[string]upvote.UpvoteType
	if rf, ok := ret.Get(0).(func(int64, []string) map[string]upvote.UpvoteType); ok {
		r0 = rf(userID, itemIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]upvote.UpvoteType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, []string) error); ok {
		r1 = rf(userID, itemIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveVote provides a mock function with given fields: userID, catItemID, increment
func (_m *Repository) RemoveVote(userID int64, catItemID string, increment int) error {
	ret := _m.Called(userID, catItemID, increment)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, string, int) error); ok {
		r0 = rf(userID, catItemID, increment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Upvote provides a mock function with given fields: userID, catItemID, increment
func (_m *Repository) Upvote(userID int64, catItemID string, increment int) error {
	ret := _m.Called(userID, catItemID, increment)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, string, int) error); ok {
		r0 = rf(userID, catItemID, increment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
