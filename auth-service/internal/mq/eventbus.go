package mq

import (
	"context"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/schema"
)

// EventBus é a interface que define os métodos para publicar eventos no RabbitMQ.
type EventBus interface {
	PublishUserCreated(ctx context.Context, evt schema.UsuarioCreated) error
	Close() error
}
