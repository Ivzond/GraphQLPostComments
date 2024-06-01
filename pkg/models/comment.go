package models

import "time"

type Comment struct {
	ID        string     `json:"id"`
	PostID    string     `json:"postId"`
	Body      string     `json:"body"`
	ParentID  *string    `json:"parentId"`
	Children  []*Comment `json:"children"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func NewComment(id, postId, body string, parentId *string) *Comment {
	return &Comment{
		ID:        id,
		PostID:    postId,
		Body:      body,
		ParentID:  parentId,
		Children:  []*Comment{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
