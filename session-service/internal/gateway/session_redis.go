package gateway

import (
	"context"
	"github.com/cntrkilril/go-chat-common/pkg/errors"
	"github.com/redis/go-redis/v9"
	"session-service/internal/entity"
	"strconv"
)

type SessionRepository struct {
	rds *redis.Client
}

func (r SessionRepository) Save(ctx context.Context, params entity.CreateSessionParams) error {
	err := r.rds.Set(ctx, params.Token, params.ID, params.ExpiresIn).Err()
	if err != nil {
		return errors.ErrUnknown
	}

	return nil
}

func (r SessionRepository) GetByToken(ctx context.Context, token string) (entity.Session, error) {
	res, err := r.rds.Get(ctx, token).Result()

	if err != nil {
		if res == "" {
			return entity.Session{}, errors.ErrSessionNotFound
		}
		return entity.Session{}, errors.ErrUnknown
	}

	id, err := strconv.Atoi(res)
	if err != nil {
		return entity.Session{}, errors.ErrUnknown
	}

	return entity.Session{ID: int64(id), Token: token}, nil
}

func (r SessionRepository) Delete(ctx context.Context, token string) error {
	err := r.rds.Del(ctx, token).Err()
	if err != nil {
		return errors.ErrUnknown
	}

	return nil
}

var _ SessionGateway = (*SessionRepository)(nil)

func NewSessionRepository(rds *redis.Client) *SessionRepository {
	return &SessionRepository{rds: rds}
}
