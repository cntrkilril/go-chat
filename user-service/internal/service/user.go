package service

import (
	"context"
	"github.com/cntrkilril/go-chat-common/pkg/errors"
	"user-service/internal/entity"
	"user-service/internal/gateway"
	"user-service/pkg/hasher"
)

type (
	UserService struct {
		userGateway gateway.UserGateway
		hasher      hasher.Interactor
	}

	UserServiceInteractor interface {
		Create(context.Context, entity.CreateUserParams) (entity.User, error)
		UpdatePassword(context.Context, entity.UpdatePasswordParams) (entity.User, error)
		GetByID(context.Context, entity.GetUserByIDParams) (entity.User, error)
		GetByIncludedUsername(context.Context, entity.GetUsersByUsernameParams) (entity.UserArray, error)
		Delete(context.Context, entity.DeleteByIDParams) error
	}
)

func (s UserService) Create(ctx context.Context, params entity.CreateUserParams) (entity.User, error) {
	user, err := s.userGateway.GetByUsername(ctx, params.Username)
	if err == nil {
		if err != errors.ErrUserNotFound {
			if user.ID != 0 {
				return entity.User{}, errors.HandleServiceError(errors.ErrUserAlreadyExist)
			}
			return entity.User{}, errors.HandleServiceError(err)
		}
	}

	params.Password, err = s.hasher.HashPassword(params.Password)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	result, err := s.userGateway.Save(ctx, params)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	return result, nil
}

func (s UserService) UpdatePassword(ctx context.Context, params entity.UpdatePasswordParams) (entity.User, error) {
	user, err := s.userGateway.GetByID(ctx, params.ID)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	if !s.hasher.CompareAndHash(user.Password, params.OldPassword) {
		return entity.User{}, errors.HandleServiceError(errors.ErrValidationError)
	}

	params.NewPassword, err = s.hasher.HashPassword(params.NewPassword)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	result, err := s.userGateway.UpdatePassword(ctx, params)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	return result, nil
}

func (s UserService) GetByID(ctx context.Context, params entity.GetUserByIDParams) (entity.User, error) {
	result, err := s.userGateway.GetByID(ctx, params.ID)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	return result, nil
}

func (s UserService) GetByIncludedUsername(ctx context.Context, params entity.GetUsersByUsernameParams) (entity.UserArray, error) {
	users, err := s.userGateway.GetByIncludedUsername(ctx, params)
	if err != nil {
		return entity.UserArray{}, errors.HandleServiceError(err)
	}

	count, err := s.userGateway.CountByIncludedUsername(ctx, params.Username)
	if err != nil {
		return entity.UserArray{}, errors.HandleServiceError(err)
	}

	return entity.UserArray{
		Users: users,
		Count: count,
	}, nil
}

func (s UserService) Delete(ctx context.Context, params entity.DeleteByIDParams) error {
	_, err := s.userGateway.GetByID(ctx, params.ID)
	if err != nil {
		return errors.HandleServiceError(err)
	}

	err = s.userGateway.Delete(ctx, params.ID)
	if err != nil {
		return errors.HandleServiceError(err)
	}

	return nil
}

var _ UserServiceInteractor = (*UserService)(nil)

func NewUserService(userGateway gateway.UserGateway, hasher hasher.Interactor) *UserService {
	return &UserService{
		userGateway: userGateway,
		hasher:      hasher,
	}
}
