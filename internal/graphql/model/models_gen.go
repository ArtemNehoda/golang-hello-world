package model

// Message is the GraphQL model for a message.
type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
}
