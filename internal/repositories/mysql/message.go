package mysql

import (
	"database/sql"

	"github.com/ArtemNehoda/golang-hello-world/internal/domain/message"
	port "github.com/ArtemNehoda/golang-hello-world/internal/ports"
)

type messageRepository struct {
	db  *sql.DB
	log port.Logger
}

// NewMessageRepository creates a repository backed by the given *sql.DB and logger.
func NewMessageRepository(db *sql.DB, log port.Logger) *messageRepository {
	return &messageRepository{db: db, log: log}
}

// GetAllMessages retrieves all messages from the database.
func (r *messageRepository) GetAllMessages() ([]message.Entity, error) {
	rows, err := r.db.Query("SELECT id, content, author, created_at FROM messages ORDER BY id ASC")
	if err != nil {
		r.log.Printf("query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []message.Entity
	for rows.Next() {
		var m message.Entity
		if err := rows.Scan(&m.ID, &m.Content, &m.Author, &m.CreatedAt); err != nil {
			r.log.Printf("row scan error: %v", err)
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *messageRepository) GetByTag(tagName string) ([]message.Entity, error) {
	rows, err := r.db.Query(`SELECT m.id, m.content, m.author, m.created_at FROM messages as m
		LEFT JOIN messages_tags_rels as r on m.id = r.message_id
		LEFT JOIN tags as t on r.tag_id = t.id
		WHERE t.name = ?
		ORDER BY id ASC`, tagName)
	if err != nil {
		r.log.Printf("query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []message.Entity
	for rows.Next() {
		var m message.Entity
		if err := rows.Scan(&m.ID, &m.Content, &m.Author, &m.CreatedAt); err != nil {
			r.log.Printf("row scan error: %v", err)
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}
