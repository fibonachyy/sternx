package sternx

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fibonachyy/sternx/config"
	"github.com/fibonachyy/sternx/internal/repository"
	"github.com/fibonachyy/sternx/internal/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main(cfg config.Config) {

	// Set up a signal handler to gracefully shut down the server on interrupt or termination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Set up the gRPC server
	grpcServer, err := setupGRPCServer(cfg)
	if err != nil {
		log.Fatalf("failed to set up gRPC server: %v", err)
	}

	// Set up the database and run migrations
	ps := repository.NewPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB)
	err = ps.Migrate(cfg.Postgres.MigrationsPath)
	if err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	// Start the gRPC server in a separate goroutine
	go func() {
		service.NewUserServiceServer(grpcServer)

		portStr := fmt.Sprintf(":%s", cfg.Grpc.Port)
		listener, err := net.Listen("tcp", portStr)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("Server is listening on %s", portStr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for signals to gracefully shut down the server
	<-stop
	log.Println("Shutting down the server gracefully...")
	grpcServer.GracefulStop()
	log.Println("Server gracefully stopped")
}

func setupGRPCServer(cfg config.Config) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(cfg.Tls.Cert, cfg.Tls.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS credentials: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(service.UnaryInterceptor),
		grpc.Creds(creds),
	}

	return grpc.NewServer(opts...), nil
}

func Register(root *cobra.Command) {
	root.PersistentFlags().String("config", "config.yaml", "read config file")
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run server",
			Run: func(cmd *cobra.Command, args []string) {
				configPath, _ := cmd.Flags().GetString("config")
				cfg := config.ReadConfig(configPath)
				main(cfg)
			},
		},
	)
}
