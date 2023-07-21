package app

import (
	"context"
	gen "github.com/cntrkilril/go-chat-common/pb/gen/session_service"
	"github.com/cntrkilril/go-chat-common/pkg/govalidator"
	myLogger "github.com/cntrkilril/go-chat-common/pkg/logger"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/redis/go-redis/v9"
	v1 "session-service/internal/controller/grpc/v1"
	"session-service/internal/gateway"
	"session-service/internal/service"

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

	l.Infof("Logger initialized successfully")

	defer myLogger.DeferLogger(l)

	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil || pong != "PONG" {
		l.Fatalf("Unable to ping redis: %v", err)
	}
	l.Infof("Connected to redis successfully")

	val := govalidator.New()

	// repos
	sessionRepo := gateway.NewSessionRepository(client)

	// services
	sessionService := service.NewSessionService(sessionRepo, cfg.Session.ExpiresIn)

	// handlers
	sessionController := v1.NewSessionController(sessionService, val)

	grpcServer := gogrpc.NewServer()
	defer grpcServer.GracefulStop()

	gen.RegisterSessionServiceServer(grpcServer, sessionController)

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
