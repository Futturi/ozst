package inmemory

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/Futturi/ozst"
)

const maxCommentLength = 2000
const pageLen = 5

// Структура для хранения данных
type Store struct {
	mu         sync.RWMutex
	posts      []*ozst.Post
	comments   []*ozst.Comment
	commentMap map[string]*ozst.Comment
}

func NewStore() *Store {
	return &Store{
		posts:      make([]*ozst.Post, 0),
		comments:   make([]*ozst.Comment, 0),
		commentMap: make(map[string]*ozst.Comment),
	}
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, commentsAllowed bool) (*ozst.Post, error) {
	r.Store.mu.Lock()
	defer r.Store.mu.Unlock()

	a := &ozst.Post{
		ID:              fmt.Sprint(len(r.Store.posts) + 1),
		Title:           title,
		Content:         content,
		CommentsAllowed: commentsAllowed,
		Comments:        make([]*ozst.Comment, 0),
	}
	r.Store.posts = append(r.Store.posts, a)
	return a, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*ozst.Comment, error) {
	if len(content) > maxCommentLength {
		slog.Info("comment out of maxcommentlen")
		return nil, errors.New("comment content exceeds maximum length")
	}
	idInt, err := strconv.Atoi(postID)
	if err != nil {
		slog.Error("incorrect id", "error", err)
		return nil, err
	}
	if !r.posts[idInt-1].CommentsAllowed {
		slog.Info("comment are not allowed in post: ", "id", idInt)
		return nil, errors.New("comment are not allowed in this post")
	}
	r.Store.mu.Lock()
	defer r.Store.mu.Unlock()

	a := &ozst.Comment{
		ID:       fmt.Sprint(len(r.Store.comments) + 1),
		PostID:   postID,
		ParentID: parentID,
		Content:  content,
		Children: make([]*ozst.Comment, 0),
	}
	r.Store.comments = append(r.Store.comments, a)
	r.Store.commentMap[a.ID] = a

	// Добавление комментария в пост или в другой комментарий
	if parentID != nil {
		if parentComment, exists := r.Store.commentMap[*parentID]; exists { // Если у комментария есть parentid
			parentComment.Children = append(parentComment.Children, a) // добавляем этому parentId children
		}
	} else { // иначе
		for _, post := range r.Store.posts { // Добавляем его в пост
			if post.ID == postID {
				post.Comments = append(post.Comments, a)
				break
			}
		}
	}

	return a, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, page int) ([]*ozst.Post, error) {
	r.Store.mu.RLock()
	defer r.Store.mu.RUnlock()

	start := (page - 1) * pageLen // 1 запись
	end := page * pageLen

	if start >= len(r.Store.posts) {
		return []*ozst.Post{}, nil
	}

	if end > len(r.Store.posts) {
		end = len(r.Store.posts)
	}

	posts := r.Store.posts[start:end]

	for _, post := range posts {
		if end > len(post.Comments) {
			end = len(post.Comments)
		}
		if start > len(post.Comments) {
			post.Comments = []*ozst.Comment{}
		} else {
			post.Comments = post.Comments[start:end]
		}
	}
	return posts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string, page int) (*ozst.Post, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &ozst.Post{}, err
	}

	r.Store.mu.RLock()
	defer r.Store.mu.RUnlock()

	if idInt < 1 || idInt > len(r.Store.posts) {
		return nil, errors.New("post not found")
	}
	post := r.Store.posts[idInt-1]
	start := (page - 1) * pageLen // 1 запись
	end := page * pageLen
	if end > len(post.Comments) {
		end = len(post.Comments)
	}
	if start > len(post.Comments) {
		post.Comments = []*ozst.Comment{}
	} else {
		post.Comments = post.Comments[start:end]
	}
	return post, nil
}

// Mutation returns ozst.MutationResolver implementation.
func (r *Resolver) Mutation() ozst.MutationResolver { return &mutationResolver{r, r.Store} }

// Query returns ozst.QueryResolver implementation.
func (r *Resolver) Query() ozst.QueryResolver { return &queryResolver{r, r.Store} }

type mutationResolver struct {
	*Resolver
	*Store
}
type queryResolver struct {
	*Resolver
	*Store
}
