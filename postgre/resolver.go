package postgre

import (
	"context"
	"github.com/Futturi/ozst"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db *sqlx.DB
}

const (
	maxCommentLength = 2000
	pageLen          = 5
)

func (r *Resolver) recur(id string) ([]*ozst.Comment, error) {
	var msgs []ozst.Comment
	newms := make([]*ozst.Comment, 0)
	query := "SELECT id, postId, parentId, content FROM comments WHERE parentId = $1"
	if err := r.Db.Select(&msgs, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		slog.Error("error with query", err)
		return nil, err
	}
	for _, msg := range msgs {
		a, err := r.recur(msg.ID)
		if err != nil {
			return nil, err
		}
		newms = append(newms, &ozst.Comment{
			ID:       msg.ID,
			Content:  msg.Content,
			PostID:   msg.PostID,
			ParentID: msg.ParentID,
			Children: a,
		})
	}
	return newms, nil
}

//go:generate mockgen -source=resolver.go -destination=mocks/mock.go
type Query interface {
	Posts(ctx context.Context, page int) ([]*ozst.Post, error)
	Post(ctx context.Context, id string, page int) (*ozst.Post, error)
}

type Mutation interface {
	CreatePost(ctx context.Context, title string, content string, commentsAllowed bool) (*ozst.Post, error)
	CreateComment(ctx context.Context, postID string, parentID *string, content string) (*ozst.Comment, error)
}
