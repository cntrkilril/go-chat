package service

import (
	"context"
	"github.com/cntrkilril/go-chat-common/pkg/errors"
	"github.com/google/uuid"
	"session-service/internal/entity"
	"session-service/internal/gateway"
	"time"
)

type (
	SessionService struct {
		sessionGateway gateway.SessionGateway
		expiresIn      time.Duration
	}

	SessionServiceInteractor interface {
		Create(context.Context, entity.CreateSessionParams) (entity.Session, error)
		GetByToken(context.Context, entity.GetSessionByTokenParams) (entity.Session, error)
		Delete(context.Context, entity.DeleteSessionParams) error
	}
)

func (s SessionService) Create(ctx context.Context, params entity.CreateSessionParams) (entity.Session, error) {
	token := uuid.New().String()
	err := s.sessionGateway.Save(ctx, entity.CreateSessionParams{
		ID:        params.ID,
		Token:     token,
		ExpiresIn: s.expiresIn,
	})
	if err != nil {
		return entity.Session{}, errors.HandleServiceError(err)
	}

	return entity.Session{ID: params.ID, Token: token}, nil
}

func (s SessionService) GetByToken(ctx context.Context, params entity.GetSessionByTokenParams) (entity.Session, error) {
	result, err := s.sessionGateway.GetByToken(ctx, params.Token)
	if err != nil {
		return entity.Session{}, errors.HandleServiceError(err)
	}

	return result, nil
}

func (s SessionService) Delete(ctx context.Context, params entity.DeleteSessionParams) error {
	_, err := s.sessionGateway.GetByToken(ctx, params.Token)
	if err != nil {
		return errors.HandleServiceError(err)
	}

	err = s.sessionGateway.Delete(ctx, params.Token)
	if err != nil {
		return errors.HandleServiceError(err)
	}

	return nil
}

var _ SessionServiceInteractor = (*SessionService)(nil)

func NewSessionService(sessionGateway gateway.SessionGateway, expiresIn time.Duration) *SessionService {
	return &SessionService{
		sessionGateway: sessionGateway,
		expiresIn:      expiresIn,
	}
}
