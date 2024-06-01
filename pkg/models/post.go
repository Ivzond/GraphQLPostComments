package models

import "time"

type Post struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Body             string    `json:"body"`
	CommentsDisabled bool      `json:"commentsDisabled"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func NewPost(id, title, body string) *Post {
	return &Post{
		ID:               id,
		Title:            title,
		Body:             body,
		CommentsDisabled: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
