package memory

import (
	"GraphQLPostComments/pkg/models"
	"sync"
)

type Storage struct {
	posts        map[string]*models.Post
	comments     map[string]*models.Comment
	postMutex    sync.RWMutex
	commentMutex sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

func (s *Storage) AddPost(post *models.Post) {
	s.postMutex.Lock()
	defer s.postMutex.Unlock()
	s.posts[post.ID] = post
}

func (s *Storage) GetPost(id string) (*models.Post, bool) {
	s.postMutex.RLock()
	defer s.postMutex.RUnlock()
	post, ok := s.posts[id]
	return post, ok
}

func (s *Storage) ListPosts() []*models.Post {
	s.postMutex.RLock()
	defer s.postMutex.RUnlock()
	posts := make([]*models.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts
}

func (s *Storage) AddComment(comment *models.Comment) {
	s.commentMutex.Lock()
	defer s.commentMutex.Unlock()
	s.comments[comment.ID] = comment
	if comment.ParentID != nil {
		parent := s.comments[*comment.ParentID]
		parent.Children = append(parent.Children, comment)
	}
}

func (s *Storage) GetCommentsByPost(postID string) []*models.Comment {
	s.commentMutex.RLock()
	defer s.commentMutex.RUnlock()
	var comments []*models.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID && comment.ParentID != nil {
			comments = append(comments, comment)
		}
	}
	return comments
}
