package v1

import (
	"context"
	gen "github.com/cntrkilril/go-chat-common/pb/gen/user_service"
	"github.com/cntrkilril/go-chat-common/pkg/errors"
	"github.com/cntrkilril/go-chat-common/pkg/govalidator"
	"google.golang.org/protobuf/types/known/emptypb"
	"user-service/internal/entity"
	"user-service/internal/service"
)

type UserController struct {
	gen.UnimplementedUserServiceServer
	userService service.UserServiceInteractor
	val         *govalidator.Validator
}

func (c UserController) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.User, error) {
	params := entity.CreateUserParams{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.User{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	result, err := c.userService.Create(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	return &gen.User{
		Id:       result.ID,
		Username: result.Username,
		Password: result.Password,
	}, nil

}

func (c UserController) UpdatePassword(ctx context.Context, req *gen.UpdatePasswordRequest) (*gen.User, error) {
	params := entity.UpdatePasswordParams{
		ID:          req.GetId(),
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.User{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	result, err := c.userService.UpdatePassword(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	return &gen.User{
		Id:       result.ID,
		Username: result.Username,
		Password: result.Password,
	}, nil

}

func (c UserController) GetUserByID(ctx context.Context, req *gen.GetByIDRequest) (*gen.User, error) {
	params := entity.GetUserByIDParams{
		ID: req.GetId(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.User{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	result, err := c.userService.GetByID(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	return &gen.User{
		Id:       result.ID,
		Username: result.Username,
		Password: result.Password,
	}, nil

}

func (c UserController) GetUsersByUsername(ctx context.Context, req *gen.GetByUsernameRequest) (*gen.UserArray, error) {
	params := entity.GetUsersByUsernameParams{
		Username: req.GetUsername(),
		PaginationRequest: entity.PaginationRequest{
			Limit:  req.GetLimit(),
			Offset: req.GetOffset(),
		},
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &gen.UserArray{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	result, err := c.userService.GetByIncludedUsername(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	var users = gen.UserArray{Users: make([]*gen.User, 0, len(result.Users)), Count: result.Count}

	for _, v := range result.Users {
		users.Users = append(users.Users,
			&gen.User{
				Id:       v.ID,
				Password: v.Password,
				Username: v.Username,
			},
		)
	}

	return &users, nil

}

func (c UserController) DeleteUser(ctx context.Context, req *gen.DeleteUserRequest) (*emptypb.Empty, error) {
	params := entity.DeleteByIDParams{
		ID: req.GetId(),
	}

	err := c.val.Validate(ctx, &params)
	if err != nil {
		return &emptypb.Empty{}, errors.HandleGrpcError(
			errors.NewError(errors.ErrValidationError.Error(), errors.ErrCodeInvalidArgument))
	}

	err = c.userService.Delete(ctx, params)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func NewUserController(userService service.UserServiceInteractor, val *govalidator.Validator) *UserController {
	return &UserController{
		userService: userService,
		val:         val,
	}
}
