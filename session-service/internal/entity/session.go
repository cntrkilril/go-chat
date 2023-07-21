package entity

import "time"

type (
	Session struct {
		ID    int64
		Token string
	}

	DeleteSessionParams struct {
		ID    int64  `validate:"required,gte=1"`
		Token string `validate:"required,gte=1"`
	}

	CreateSessionParams struct {
		ID        int64 `validate:"required,gte=1"`
		Token     string
		ExpiresIn time.Duration
	}

	GetSessionByTokenParams struct {
		Token string `validate:"required,gte=1"`
	}
)
