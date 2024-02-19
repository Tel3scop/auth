package main

import (
	"context"
	"log"
	"net"

	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/handlers"
	"github.com/Tel3scop/auth/internal/services"
	"github.com/Tel3scop/auth/internal/storages"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("cannot get config: %v", err)
	}

	lis, err := net.Listen(cfg.GRPC.Protocol, cfg.GRPC.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer pool.Close()
	storageContainer := storages.New(pool)
	serviceContainer := services.New(cfg, storageContainer)

	s := grpc.NewServer()
	reflection.Register(s)

	server := handlers.Run(cfg, serviceContainer)
	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
