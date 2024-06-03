package memory

import (
	"GraphQLPostComments/api/model"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Storage struct {
	Posts        map[string]*model.Post
	comments     map[string]*model.Comment
	postMutex    sync.RWMutex
	commentMutex sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		Posts:    make(map[string]*model.Post),
		comments: make(map[string]*model.Comment),
	}
}

func (s *Storage) GetPosts() ([]*model.Post, error) {
	s.postMutex.RLock()
	defer s.postMutex.RUnlock()

	posts := make([]*model.Post, 0, len(s.Posts))
	for _, post := range s.Posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Storage) GetPostByID(id string) (*model.Post, error) {
	s.postMutex.RLock()
	defer s.postMutex.RUnlock()

	post, ok := s.Posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *Storage) CreatePost(post *model.Post) (*model.Post, error) {
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now().Format(time.RFC3339)

	s.postMutex.Lock()
	defer s.postMutex.Unlock()

	s.Posts[post.ID] = post
	return post, nil
}

func (s *Storage) UpdatePost(updatedPost *model.Post) (*model.Post, error) {
	s.postMutex.Lock()
	defer s.postMutex.Unlock()

	_, ok := s.Posts[updatedPost.ID]
	if !ok {
		return nil, errors.New("post not found")
	}

	s.Posts[updatedPost.ID] = updatedPost
	return updatedPost, nil
}

func (s *Storage) DeletePost(id string) (bool, error) {
	s.postMutex.Lock()
	defer s.postMutex.Unlock()

	_, ok := s.Posts[id]
	if !ok {
		return false, errors.New("post not found")
	}
	delete(s.Posts, id)
	return true, nil
}

func (s *Storage) GetCommentByID(id string) (*model.Comment, error) {
	s.commentMutex.RLock()
	defer s.commentMutex.RUnlock()
	comment, ok := s.comments[id]
	if !ok {
		return nil, errors.New("comment not found")
	}
	return comment, nil
}

func (s *Storage) GetComments(postID string, page int, limit int) (*model.CommentPage, error) {
	s.commentMutex.RLock()
	defer s.commentMutex.RUnlock()

	comments := make([]*model.Comment, 0)
	for _, comment := range s.comments {
		if comment.Post.ID == postID {
			comments = append(comments, comment)
		}
	}

	totalComments := len(comments)
	start := (page - 1) * limit
	end := start + limit
	if start > totalComments {
		start = totalComments
	}
	if end > totalComments {
		end = totalComments
	}

	return &model.CommentPage{
		Comments: comments[start:end],
		PageInfo: &model.PageInfo{
			CurrentPage:   page,
			TotalPages:    (totalComments + limit - 1) / limit,
			TotalComments: totalComments,
		},
	}, nil
}

func (s *Storage) CreateComment(comment *model.Comment) (*model.Comment, error) {
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now().Format(time.RFC3339)

	s.commentMutex.Lock()
	defer s.commentMutex.Unlock()

	s.comments[comment.ID] = comment
	return comment, nil
}

func (s *Storage) DeleteComment(id string) (bool, error) {
	s.commentMutex.Lock()
	defer s.commentMutex.Unlock()

	_, ok := s.comments[id]
	if !ok {
		return false, errors.New("comment not found")
	}

	delete(s.comments, id)
	return true, nil
}
