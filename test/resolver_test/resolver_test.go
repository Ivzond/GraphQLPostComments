package api_test

import (
	"GraphQLPostComments/api"
	"GraphQLPostComments/api/model"
	"GraphQLPostComments/internal/storage/memory"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Test the CreatePost resolver
func TestCreatePost(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	title := "New Post"
	content := "This is a new post"
	authorID := "author1"

	post, err := resolver.Mutation().CreatePost(ctx, title, content, authorID)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, title, post.Title)
	assert.Equal(t, content, post.Content)
	assert.Equal(t, authorID, post.Author.ID)
}

// Test the Posts resolver
func TestPosts(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	post1 := &model.Post{
		ID:        "1",
		Title:     "Post 1",
		Content:   "Content 1",
		Author:    &model.User{ID: "1"},
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	post2 := &model.Post{
		ID:        "2",
		Title:     "Post 2",
		Content:   "Content 2",
		Author:    &model.User{ID: "2"},
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	store.CreatePost(post1)
	store.CreatePost(post2)

	posts, err := resolver.Query().Posts(ctx)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}

// Test the Post resolver
func TestPost(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	post := &model.Post{
		ID:        "1",
		Title:     "Post 1",
		Content:   "Content 1",
		Author:    &model.User{ID: "1"},
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	store.Posts[post.ID] = post

	fetchedPost, err := resolver.Query().Post(ctx, "1")
	assert.NoError(t, err)
	assert.NotNil(t, fetchedPost)
	assert.Equal(t, post.ID, fetchedPost.ID)
}

// Test the UpdatePost resolver
func TestUpdatePost(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	post := &model.Post{
		ID:        "1",
		Title:     "Post 1",
		Content:   "Content 1",
		Author:    &model.User{ID: "1"},
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	store.Posts[post.ID] = post

	newTitle := "Updated Post 1"
	updatedPost, err := resolver.Mutation().UpdatePost(ctx, "1", &newTitle, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, updatedPost)
	assert.Equal(t, newTitle, updatedPost.Title)
}

// Test the DeletePost resolver
func TestDeletePost(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	post := &model.Post{
		ID:        "1",
		Title:     "Post 1",
		Content:   "Content 1",
		Author:    &model.User{ID: "1"},
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	store.Posts[post.ID] = post

	deleted, err := resolver.Mutation().DeletePost(ctx, "1")
	assert.NoError(t, err)
	assert.True(t, *deleted)

	fetchedPost, err := resolver.Query().Post(ctx, "1")
	assert.Error(t, err)
	assert.Nil(t, fetchedPost)
}

// Test the CreateComment resolver
func TestCreateComment(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	post := &model.Post{
		ID:              "1",
		Title:           "Post 1",
		Content:         "Content 1",
		Author:          &model.User{ID: "1"},
		CreatedAt:       time.Now().Format(time.RFC3339),
		CommentsEnabled: true,
	}
	store.Posts[post.ID] = post

	content := "This is a comment"
	authorID := "author1"

	comment, err := resolver.Mutation().CreateComment(ctx, "1", content, authorID, "")
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, content, comment.Content)
	assert.Equal(t, authorID, comment.Author.ID)
}

// Test the Comments resolver
func TestComments(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	post := &model.Post{
		ID:              "1",
		Title:           "Post 1",
		Content:         "Content 1",
		Author:          &model.User{ID: "1"},
		CreatedAt:       time.Now().Format(time.RFC3339),
		CommentsEnabled: true,
	}
	store.Posts[post.ID] = post

	comment1 := &model.Comment{
		Content:   "Comment 1",
		Author:    &model.User{ID: "1"},
		Post:      post,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	comment2 := &model.Comment{
		Content:   "Comment 2",
		Author:    &model.User{ID: "2"},
		Post:      post,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	store.CreateComment(comment1)
	store.CreateComment(comment2)

	page, limit := 1, 10
	commentsPage, err := resolver.Query().Comments(ctx, "1", &page, &limit)
	assert.NoError(t, err)
	assert.NotNil(t, commentsPage)
	assert.Len(t, commentsPage.Comments, 2)
}

// Test the DeleteComment resolver
func TestDeleteComment(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	comment := &model.Comment{
		Content:   "Comment 1",
		Author:    &model.User{ID: "1"},
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	createdComment, err := store.CreateComment(comment)

	deleted, err := resolver.Mutation().DeleteComment(ctx, createdComment.ID)
	assert.NoError(t, err)
	assert.True(t, *deleted)

	fetchedComment, err := resolver.Store.GetCommentByID(createdComment.ID)
	assert.Error(t, err)
	assert.Nil(t, fetchedComment)
}

// Test the CommentAdded resolver
func TestCommentAdded(t *testing.T) {
	store := memory.NewStorage()
	resolver := api.NewResolver(store)
	ctx := context.Background()

	post := &model.Post{
		ID:              "1",
		Title:           "Post 1",
		Content:         "Content 1",
		Author:          &model.User{ID: "1"},
		CreatedAt:       time.Now().Format(time.RFC3339),
		CommentsEnabled: true,
	}
	store.Posts[post.ID] = post

	commentsChan, err := resolver.Subscription().CommentAdded(ctx, "1")
	assert.NoError(t, err)
	assert.NotNil(t, commentsChan)

	content := "This is a comment"
	authorID := "author1"
	go resolver.Mutation().CreateComment(ctx, post.ID, content, authorID, "")

	select {
	case comment := <-commentsChan:
		assert.Equal(t, content, comment.Content)
		assert.Equal(t, authorID, comment.Author.ID)
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for comment added")
	}
}
