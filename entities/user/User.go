package user

import "encoding/json"

type User struct {
	Id           int             `json:"id" example:"1"`
	Name         string          `json:"name" example:"John Doe"`
	Followers_id json.RawMessage `json:"followers_id" example:"[2, 3, 4]"`
	Following_id json.RawMessage `json:"following_id" example:"[5, 6]"`
	Feed         json.RawMessage `json:"feed" example:"[101, 102]"`
}

type UserInput struct {
	Name string `json:"name" example:"John Doe" binding:"required"`
}
