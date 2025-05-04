package tweet

import (
	"time"
)

type Tweet struct {
	Id        int
	Timestamp time.Time
	Message   string
	Author_id int
}


type TweetInput struct {
	Message  string `json:"message" example:"Hola mundo" binding:"required"`
	AuthorID int    `json:"author_id" example:"1" binding:"required"`
}
