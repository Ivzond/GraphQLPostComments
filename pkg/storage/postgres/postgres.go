package postgres

import (
	"GraphQLPostComments/pkg/models"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dataSourceName string) (*Storage, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) AddPost(post *models.Post) error {
	query := `INSERT INTO posts (id, title, body, comments_disabled, created_at, update_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Exec(query, post.ID, post.Title, post.Body, post.CommentsDisabled, post.CreatedAt, post.UpdatedAt)
	return err
}

func (s *Storage) GetPost(id string) (*models.Post, error) {
	query := `SELECT id, title, body, comments_disabled, created_at, updated_at FROM posts WHERE id = $1`
	row := s.db.QueryRow(query, id)
	post := &models.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Body, &post.CommentsDisabled, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Storage) ListPosts() ([]*models.Post, error) {
	query := `SELECT id, title, body, comments_disabled, created_at, update_at FROM posts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.CommentsDisabled, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Storage) AddComment(comment *models.Comment) error {
	query := `INSERT INTO comments (id, post_id, body, parent_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Exec(query, comment.ID, comment.PostID, comment.Body, comment.ParentID, comment.CreatedAt, comment.UpdatedAt)
	return err
}

func (s *Storage) GetCommentsByPost(postID string) ([]*models.Comment, error) {
	query := `SELECT id, post_id, body, parent_id, created_at, update_at FROM comments WHERE post_id = $1 AND parent_id IS NULL`
	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.Body, &comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
