package app

import (
	myLogger "github.com/cntrkilril/go-chat-common/pkg/logger"
	"github.com/cntrkilril/go-chat-common/pkg/postgres"
	v1 "user-service/internal/controller/grpc/v1"
	"user-service/internal/gateway"
	"user-service/internal/service"
	"user-service/pkg/hasher"

	gen "github.com/cntrkilril/go-chat-common/pb/gen/user_service"
	"github.com/cntrkilril/go-chat-common/pkg/govalidator"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	// db driver
	gogrpc "google.golang.org/grpc"
)

func Run() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	l := myLogger.NewLogger(myLogger.Config{
		Logger: myLogger.Logger{
			Level: *cfg.Logger.Level,
		},
	})

	defer myLogger.DeferLogger(l)

	l.Infof("Logger initialized successfully")

	db, err := postgres.InitPsqlDB(postgres.Config{
		Postgres: postgres.Postgres{
			ConnString:      cfg.Postgres.ConnString,
			MaxOpenConns:    cfg.Postgres.MaxOpenConns,
			ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.Postgres.ConnMaxIdleTime,
			MaxIdleConns:    cfg.Postgres.MaxOpenConns,
			MigrationsPath:  cfg.Postgres.MigrationsPath,
			DBName:          cfg.Postgres.DBName,
			AutoMigrate:     cfg.Postgres.AutoMigrate,
		},
	}, l)

	l.Debug("Connected to PostgreSQL")

	//postgres.DeferClose(db, l)

	val := govalidator.New()
	h := hasher.NewHasher()

	// repos
	userRepo := gateway.NewUserRepository(db)

	// services
	userService := service.NewUserService(userRepo, h)

	// handlers
	userController := v1.NewUserController(userService, val)

	grpcServer := gogrpc.NewServer()
	defer grpcServer.GracefulStop()

	gen.RegisterUserServiceServer(grpcServer, userController)

	go func() {
		lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port))
		if err != nil {
			log.Fatalf("tcp sock: %s", err.Error())
		}
		defer func(lis net.Listener) {
			err = lis.Close()
			if err != nil {
				l.Error(err)
				return
			}
		}(lis)

		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("GRPC server: %s", err.Error())
		}
	}()

	l.Debug("Started GRPC server")

	l.Debug("Application has started")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	l.Info("Application has been shut down")

	return
}
