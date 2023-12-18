package sternx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fibonachyy/sternx/config"
	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/repository"
	"github.com/fibonachyy/sternx/internal/service"
	"github.com/fibonachyy/sternx/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main(cfg config.Config) {

	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	// Set up a signal handler to gracefully shut down the server on interrupt or termination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Set up the database and run migrations
	dbLogger := logger.NewDatabaseLogger()
	ps := repository.NewPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB, dbLogger)
	err := ps.Migrate(cfg.Postgres.MigrationsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run database migrations:")
	}

	creds, err := credentials.NewServerTLSFromFile(cfg.Tls.Cert, cfg.Tls.Key)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load TLS credentials:")

	}

	// Set up the gRPC server
	grpcServer, err := setupGRPCServer(cfg, creds, ps)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set up gRPC server:")
	}

	// Start the gRPC server in a separate goroutine
	go func() {
		portStr := fmt.Sprintf(":%s", cfg.Grpc.Port)
		listener, err := net.Listen("tcp", portStr)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to listen:")
		}
		log.Info().Msgf("gRPC server is listening on %s", portStr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC")
		}
	}()

	mux := http.NewServeMux()

	// Serve Swagger JSON definition for testing the gRPC gateway. This endpoint is also compatible
	// with Swagger UI extensions, providing a better user experience.
	// Access the Swagger JSON file at: http://serveradd:getWayport/swagger/
	const swaggerJSONPath = "/swagger/"
	swaggerJSONHandler := http.StripPrefix(swaggerJSONPath, http.FileServer(http.Dir("doc/swagger")))
	mux.Handle(swaggerJSONPath, swaggerJSONHandler)

	// Set up the gRPC gateway
	gatewayMux, err := setupGRPCGateway(fmt.Sprintf("127.0.0.1:%s", cfg.Grpc.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set up gRPC gateway:")
	}

	// Combine the ServeMux instances
	mux.Handle("/", gatewayMux)

	// Set up the HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Grpc.GetwayPort),
		Handler: mux,
	}

	go func() {
		gatewayAddr := fmt.Sprintf(":%s", cfg.Grpc.GetwayPort)
		log.Info().Msgf("gRPC gateway is listening on %s", gatewayAddr)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC gateway:")
		}
	}()

	// Wait for signals to gracefully shut down the server
	<-stop
	log.Info().Msg("Shutting down the server gracefully...")

	// Gracefully stop the gRPC server
	grpcServer.GracefulStop()
	// Gracefully stop the HTTP server
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("Failed to shut down HTTP server gracefully:")
	}
	log.Info().Msg("Server gracefully stopped")
}

func setupGRPCServer(cfg config.Config, creds credentials.TransportCredentials, ps repository.IRepository) (*grpc.Server, error) {
	log.Info().Msg("Setting up gRPC server...")

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(service.UnaryInterceptor),
		// grpc.Creds(creds),
		// Note: TLS (Transport Layer Security) is currently disabled for the service to facilitate development purposes.
		// Enabling TLS requires valid certificate files. Without them, testing the application becomes restricted.
		// Please be aware that handling TLS, including acquiring and configuring proper certificate files, has been considered and implemented.😊
	}
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)

	// Register your gRPC service implementation
	conf := service.Config{JWTDuration: time.Minute * time.Duration(cfg.Jwt.ExpireMin), TokenSymmetricKey: cfg.Jwt.TokenSymmetricKey}
	userServiceServer, err := service.NewUserServiceServer(ps, conf)
	if err != nil {
		return nil, err
	}
	userpb.RegisterUserServiceServer(grpcServer, userServiceServer)
	return grpcServer, nil
}

func setupGRPCGateway(serverAddr string) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %v", err)
	}
	err = userpb.RegisterUserServiceHandler(context.Background(), mux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gRPC gateway: %v", err)
	}

	return mux, nil
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
