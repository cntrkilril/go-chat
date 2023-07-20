package gateway

import (
	"context"
	"user-service/internal/entity"
)

type UserGateway interface {
	Save(context.Context, entity.CreateUserParams) (entity.User, error)
	UpdatePassword(context.Context, entity.UpdatePasswordParams) (entity.User, error)
	GetByID(context.Context, int64) (entity.User, error)
	GetByUsername(context.Context, string) (entity.User, error)
	GetByIncludedUsername(context.Context, entity.GetUsersByUsernameParams) ([]entity.User, error)
	CountByIncludedUsername(context.Context, string) (int64, error)
	Delete(context.Context, int64) error
}
