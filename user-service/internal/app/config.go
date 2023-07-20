package app

import (
	"github.com/cntrkilril/go-chat-common/pkg/postgres"
	"github.com/go-playground/validator/v10"
)

type (
	Config struct {
		Logger   Logger            `validate:"required"`
		GRPC     GRPC              `validate:"required"`
		Postgres postgres.Postgres `validate:"required"`
	}

	Logger struct {
		Level *int8 `validate:"required"`
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
			Port: "8065",
		},
		Logger: Logger{
			Level: &defaultLogLevel,
		},
		Postgres: postgres.Postgres{
			ConnString:      "postgresql://root:pass@127.0.0.1:5432/users?sslmode=disable&application_name=user-service",
			MaxOpenConns:    10,
			ConnMaxLifetime: 20,
			MaxIdleConns:    15,
			ConnMaxIdleTime: 30,
			AutoMigrate:     true,
			DBName:          "users",
			MigrationsPath:  "db/migration",
		},
	}

	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
