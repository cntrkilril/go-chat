package gateway

import (
	"context"
	"session-service/internal/entity"
)

type SessionGateway interface {
	Save(context.Context, entity.CreateSessionParams) error
	GetByToken(context.Context, string) (entity.Session, error)
	Delete(context.Context, string) error
}
