package storage

import "GraphQLPostComments/api/model"

type Storage interface {
	GetPosts() ([]*model.Post, error)
	GetPostByID(id string) (*model.Post, error)
	CreatePost(*model.Post) (*model.Post, error)
	UpdatePost(*model.Post) (*model.Post, error)
	DeletePost(id string) (bool, error)

	GetCommentByID(id string) (*model.Comment, error)
	GetComments(postID string, page int, limit int) (*model.CommentPage, error)
	CreateComment(comment *model.Comment) (*model.Comment, error)
	DeleteComment(id string) (bool, error)
}
