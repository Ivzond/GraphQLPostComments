package memory_test

import (
	"GraphQLPostComments/api/model"
	"GraphQLPostComments/internal/storage/memory"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePost(t *testing.T) {
	// Создаем новое хранилище в памяти
	store := memory.NewStorage()

	// Тестовые данные для создания поста
	post := &model.Post{
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  &model.User{ID: "user1"},
	}

	// Вызываем метод создания поста
	createdPost, err := store.CreatePost(post)
	if err != nil {
		t.Errorf("Error creating post: %v", err)
		return
	}

	// Проверяем, что пост был успешно создан
	if createdPost == nil {
		t.Errorf("Expected created post to not be nil, got nil")
	}

	// Проверяем, что ID поста был установлен
	if createdPost.ID == "" {
		t.Errorf("Expected non-empty post ID, got empty")
	}
}

func TestGetPostByID(t *testing.T) {
	// Создаем новое хранилище в памяти
	store := memory.NewStorage()

	// Создаем тестовый пост
	post := &model.Post{
		ID:      "post1",
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  &model.User{ID: "user1"},
	}

	// Добавляем пост в хранилище
	store.Posts["post1"] = post

	// Вызываем метод получения поста по ID
	retrievedPost, err := store.GetPostByID("post1")
	if err != nil {
		t.Errorf("Error getting post by ID: %v", err)
		return
	}

	// Проверяем, что получен корректный пост
	if retrievedPost == nil {
		t.Errorf("Expected retrieved post to be non-nil, got nil")
	}

	// Проверяем, что данные поста совпадают с ожидаемыми
	if retrievedPost.ID != post.ID || retrievedPost.Title != post.Title || retrievedPost.Content != post.Content || retrievedPost.Author.ID != post.Author.ID {
		t.Errorf("Mismatch in retrieved post data")
	}
}

func TestUpdatePost(t *testing.T) {
	store := memory.NewStorage()

	post := &model.Post{
		ID:      "post1",
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  &model.User{ID: "user1"},
	}

	store.Posts["post1"] = post

	updatedPost := &model.Post{
		ID:      "post1",
		Title:   "Updated Post",
		Content: "This is a test post",
		Author:  &model.User{ID: "user1"},
	}

	p, err := store.UpdatePost(updatedPost)
	assert.NoError(t, err)
	assert.Equal(t, updatedPost, p)
}

func TestDeletePost(t *testing.T) {
	store := memory.NewStorage()

	post := &model.Post{
		ID:      "post1",
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  &model.User{ID: "user1"},
	}

	store.Posts["post1"] = post

	deleted, err := store.DeletePost("post1")
	assert.NoError(t, err)
	assert.True(t, deleted)

	p, err := store.GetPostByID("post1")
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestCreateComment(t *testing.T) {
	store := memory.NewStorage()

	comment := &model.Comment{
		ID:      "comment1",
		Content: "This is a test comment",
		Author:  &model.User{ID: "user1"},
		Post:    &model.Post{ID: "post1"},
	}

	createdComment, err := store.CreateComment(comment)
	assert.NoError(t, err)
	assert.Equal(t, comment, createdComment)
}

func TestGetComments(t *testing.T) {
	store := memory.NewStorage()

	comment1 := &model.Comment{
		ID:      "comment1",
		Content: "This is a test comment 1",
		Author:  &model.User{ID: "user1"},
		Post:    &model.Post{ID: "post1"},
	}

	comment2 := &model.Comment{
		ID:      "comment2",
		Content: "This is a test comment 2",
		Author:  &model.User{ID: "user1"},
		Post:    &model.Post{ID: "post1"},
	}

	store.CreateComment(comment1)
	store.CreateComment(comment2)

	commentPage, err := store.GetComments("post1", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(commentPage.Comments))
}

func TestDeleteComment(t *testing.T) {
	store := memory.NewStorage()

	comment := &model.Comment{
		Content:   "Comment 1",
		Author:    &model.User{ID: "1"},
		Post:      &model.Post{ID: "1"},
		CreatedAt: "2022-01-01T00:00:00Z",
	}

	createdComment, _ := store.CreateComment(comment)

	deleted, err := store.DeleteComment(createdComment.ID)
	fmt.Println(deleted)
	assert.NoError(t, err)
	assert.True(t, deleted)

	c, err := store.GetCommentByID("1")
	assert.Error(t, err)
	assert.Nil(t, c)
}
