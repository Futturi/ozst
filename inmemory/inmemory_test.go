package inmemory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCreatePost(t *testing.T) {
	store := NewStore()
	mutation := mutationResolver{Store: store}

	post, err := mutation.CreatePost(context.Background(), "Title", "Content", true)
	assert.NoError(t, err)
	assert.Equal(t, "Title", post.Title)
	assert.Equal(t, "Content", post.Content)
	assert.True(t, post.CommentsAllowed)
}

func TestCreateComment(t *testing.T) {
	store := NewStore()
	mutation := mutationResolver{Store: store}

	post, err := mutation.CreatePost(context.Background(), "Title", "Content", true)
	assert.NoError(t, err)

	comment, err := mutation.CreateComment(context.Background(), post.ID, nil, "Comment Content")
	assert.NoError(t, err)
	assert.Equal(t, "Comment Content", comment.Content)
	assert.Equal(t, post.ID, comment.PostID)
}

func TestPostsPagination(t *testing.T) {
	store := NewStore()
	mutation := mutationResolver{Store: store}
	query := queryResolver{Store: store}

	for i := 1; i <= 12; i++ {
		_, err := mutation.CreatePost(context.Background(), "Title "+strconv.Itoa(i), "Content "+strconv.Itoa(i), true)
		assert.NoError(t, err)
	}

	posts, err := query.Posts(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, posts, 5)

	posts, err = query.Posts(context.Background(), 3)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}
