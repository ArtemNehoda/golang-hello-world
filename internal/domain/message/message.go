package message

import (
	"time"
)

// Entity represents a message domain entity
type Entity struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
}

// New creates a new Entity with the given content and author,
// setting CreatedAt to the current UTC time.
func New(content, author string) *Entity {
	return &Entity{
		Content:   content,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	}
}
