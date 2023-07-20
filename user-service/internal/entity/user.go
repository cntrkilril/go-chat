package entity

type (
	User struct {
		ID       int64  `db:"id"`
		Username string `db:"username"`
		Password string `db:"password"`
	}

	UserArray struct {
		Users []User
		Count int64
	}

	PaginationRequest struct {
		Limit  int64 `validate:"required"`
		Offset int64 `validate:"gte=0"`
	}

	CreateUserParams struct {
		Username string `validate:"required,gte=1"`
		Password string `validate:"required,gte=8"`
	}

	GetUsersByUsernameParams struct {
		Username string `validate:"required,gte=1"`
		PaginationRequest
	}

	UpdatePasswordParams struct {
		ID          int64  `validate:"required,gte=1"`
		OldPassword string `validate:"required,gte=8"`
		NewPassword string `validate:"required,gte=8"`
	}

	GetUserByIDParams struct {
		ID int64 `validate:"required,gte=1"`
	}

	DeleteByIDParams struct {
		ID int64 `validate:"required,gte=1"`
	}
)
