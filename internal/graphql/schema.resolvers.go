package graphql

import (
	"context"
	"fmt"

	"github.com/ArtemNehoda/golang-hello-world/internal/graphql/model"
)

// QueryResolver defines the interface for Query field resolvers.
type QueryResolver interface {
	Messages(ctx context.Context) ([]*model.Message, error)
}

type queryResolver struct{ *Resolver }

// Query returns the QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Messages resolves the `messages` query field.
func (r *queryResolver) Messages(ctx context.Context) ([]*model.Message, error) {
	entities, err := r.Service.GetAllMessages()
	if err != nil {
		r.Logger.Printf("Failed to get messages: %v", err)
		return nil, fmt.Errorf("failed to retrieve messages")
	}

	msgs := make([]*model.Message, len(entities))
	for i, e := range entities {
		msgs[i] = &model.Message{
			ID:        fmt.Sprintf("%d", e.ID),
			Content:   e.Content,
			Author:    e.Author,
			CreatedAt: e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return msgs, nil
}
