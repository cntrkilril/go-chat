package app

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type (
	Config struct {
		Logger  Logger  `validate:"required"`
		GRPC    GRPC    `validate:"required"`
		Redis   Redis   `validate:"required"`
		Session Session `validate:"required"`
	}

	Logger struct {
		Level *int8 `validate:"required"`
	}

	Redis struct {
		Host     string `validate:"required"`
		Port     string `validate:"required"`
		Password string
		DB       int
	}

	Session struct {
		ExpiresIn time.Duration `validate:"required"`
	}

	GRPC struct {
		Host string `validate:"required"`
		Port string `validate:"required"`
	}
)

func LoadConfig() (*Config, error) {

	defaultLogLevel := int8(-1)

	cfg := &Config{
		GRPC: GRPC{
			Host: "localhost",
			Port: "8070",
		},
		Logger: Logger{
			Level: &defaultLogLevel,
		},
		Redis: Redis{
			Host:     "127.0.0.1",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		Session: Session{
			ExpiresIn: time.Hour * 24,
		},
	}

	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
