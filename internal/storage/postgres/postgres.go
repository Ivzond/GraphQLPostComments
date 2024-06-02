package postgres

import (
	"GraphQLPostComments/api/model"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dataSourceName string) (*Storage, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetPosts() ([]*model.Post, error) {
	rows, err := s.db.Query("SELECT id, title, content, author_id, comments_enabled, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*model.Post{}
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author.ID, post.CommentsEnabled, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Storage) GetPostByID(id string) (*model.Post, error) {
	post := &model.Post{}
	err := s.db.QueryRow("SELECT id, title, content, author_id, comments_enabled, created_at FROM posts WHERE id = $1", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Author.ID, post.CommentsEnabled, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return post, nil
}

func (s *Storage) CreatePost(post *model.Post) (*model.Post, error) {
	post.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := s.db.Exec("INSERT INTO posts (id, title, content, author_id, comments_enabled, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		post.ID, post.Title, post.Content, post.Author.ID, post.CommentsEnabled, post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Storage) UpdatePost(post *model.Post) (*model.Post, error) {
	_, err := s.db.Exec("UPDATE posts SET title = $1, content = $2, comments_enabled = $3 WHERE id = $4",
		post.Title, post.Content, post.CommentsEnabled, post.ID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Storage) DeletePost(id string) (bool, error) {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Storage) GetComments(postID string, page int, limit int) (*model.CommentPage, error) {
	var totalComments int
	err := s.db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = $1", postID).Scan(&totalComments)
	if err != nil {
		return nil, err
	}
	start := (page - 1) * limit

	rows, err := s.db.Query("SELECT id, content, author_id, post_id, parent_id, created_at FROM comments WHERE post_id = $1 LIMIT $2 OFFSET $3",
		postID, limit, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.Comment{}
	for rows.Next() {
		comment := &model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.Author.ID, &comment.Post.ID, &comment.Parent.ID, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &model.CommentPage{
		Comments: comments,
		PageInfo: &model.PageInfo{
			CurrentPage:   page,
			TotalPages:    (totalComments + limit - 1) / limit,
			TotalComments: totalComments,
		},
	}, nil
}

func (s *Storage) CreateComment(comment *model.Comment) (*model.Comment, error) {
	comment.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := s.db.Exec("INSERT INTO comments (id, content, author_id, post_id, parent_id, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		comment.ID, comment.Content, comment.Author.ID, comment.Post.ID, comment.Parent.ID, comment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *Storage) DeleteComment(id string) (bool, error) {
	_, err := s.db.Exec("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	return true, nil
}
