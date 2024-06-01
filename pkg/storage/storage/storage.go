package storage

import "GraphQLPostComments/pkg/models"

type Storage interface {
	AddPost(post *models.Post) error
	GetPost(id string) (*models.Post, error)
	ListPosts() ([]*models.Post, error)

	AddComment(comment *models.Comment) error
	GetCommentsByPost(id string) ([]*models.Comment, error)
}
