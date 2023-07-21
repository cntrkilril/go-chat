package v1

import (
	"context"
	gen "github.com/cntrkilril/go-chat-common/pb/gen/session_service"
	"github.com/cntrkilril/go-chat-common/pkg/errors"
	"github.com/cntrkilril/go-chat-common/pkg/govalidator"
	"google.golang.org/protobuf/types/known/emptypb"
	"session-service/internal/entity"
	"session-service/internal/service"
)

type SessionController struct {
	gen.UnimplementedSessionServiceServer
	sessionService service.SessionServiceInteractor
	val            *govalidator.Validator
}

func (c SessionController) CreateSession(ctx context.Context, req *gen.CreateSessionRequest) (*gen.Session, error) {
	params := entity.CreateSessionParams{
		ID: req.GetId(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.Session{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	result, err := c.sessionService.Create(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	return &gen.Session{
		Id:    result.ID,
		Token: result.Token,
	}, nil

}

func (c SessionController) GetSessionByToken(ctx context.Context, req *gen.GetSessionByTokenRequest) (*gen.Session, error) {
	params := entity.GetSessionByTokenParams{
		Token: req.GetToken(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.Session{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	result, err := c.sessionService.GetByToken(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	return &gen.Session{
		Id:    result.ID,
		Token: result.Token,
	}, nil

}

func (c SessionController) DeleteSession(ctx context.Context, req *gen.Session) (*emptypb.Empty, error) {
	params := entity.DeleteSessionParams{
		ID:    req.GetId(),
		Token: req.GetToken(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &emptypb.Empty{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	err = c.sessionService.Delete(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func NewSessionController(sessionService service.SessionServiceInteractor, val *govalidator.Validator) *SessionController {
	return &SessionController{
		sessionService: sessionService,
		val:            val,
	}
}
