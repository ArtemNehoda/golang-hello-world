package graphql

import (
	"github.com/ArtemNehoda/golang-hello-world/internal/domain/message"
	"github.com/ArtemNehoda/golang-hello-world/internal/ports"
)

// MessageService is the interface the resolver depends on.
type MessageService interface {
	GetAllMessages() ([]message.Entity, error)
}

// Resolver holds the application's dependencies.
type Resolver struct {
	Service MessageService
	Logger  ports.Logger
}
