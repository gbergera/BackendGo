package follow

import (
	"time"
)

type Follow struct {
	FollowerID int       `json:"follower_id"`
	FollowedID int       `json:"followed_id"`
	Timestamp  time.Time `json:"timestamp"`
}

type FollowInput struct {
	FollowerID int `json:"follower_id" example:"1"`
	FollowedID int `json:"followed_id" example:"2"`
}
