package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ximura/gowallet/api"
	grpcCtrl "github.com/ximura/gowallet/internal/controller/grpc"
	"github.com/ximura/gowallet/internal/core/server/grpc"
	"github.com/ximura/gowallet/internal/core/service"
	"github.com/ximura/gowallet/internal/repository"
	googleGrpc "google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connectString := "host=postgres port=5432 user=postgres password=postgres dbname=wallet sslmode=disable"
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create DB repo %w", err))
	}
	defer db.Close()

	repo := repository.NewWalletRepo(db)
	walletService := service.NewWalletService(&repo)
	walletController := grpcCtrl.NewWalletController(&walletService)

	grpcService := grpc.NewGRPCService(50051)
	defer grpcService.Close()

	grpcService.Register(func(server *googleGrpc.Server) {
		api.RegisterWalletServiceServer(server, walletController)
	})
	go func() {
		grpcService.Run(ctx)
	}()

	AddShutdownHook(grpcService, &repo)
}

func AddShutdownHook(closers ...io.Closer) {
	log.Println("listening signals...")
	c := make(chan os.Signal, 1)
	signal.Notify(
		c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)

	<-c
	log.Println("graceful shutdown...")

	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			log.Println(fmt.Errorf("failed to stop closer %w", err))
		}
	}

	log.Println("completed graceful shutdown")
}
