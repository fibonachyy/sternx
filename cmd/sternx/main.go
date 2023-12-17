package sternx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	grpcServer, err := setupGRPCServer(creds, ps)
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

	gatewayMux, err := setupGRPCGateway(fmt.Sprintf("localhost:%s", cfg.Grpc.GetwayPort), creds)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set up gRPC gateway:")
	}

	// Create a ServeMux for custom routing
	mux := http.NewServeMux()

	const (
		grpcGatewayPath = "/"
		swaggerUIPath   = "/swagger/"
		swaggerJSONPath = "/doc/swagger/"
		swaggerBaseURL  = "https://raw.githubusercontent.com/swagger-api/swagger-ui/master/dist/"
	)

	// ...

	// Register gRPC gateway handler
	mux.Handle(grpcGatewayPath, gatewayMux)

	// Serve Swagger UI from /swagger/
	swaggerHandler := http.StripPrefix(swaggerUIPath, http.FileServer(http.Dir("doc/swagger")))
	mux.Handle(swaggerUIPath, swaggerHandler)

	// Serve Swagger definition from /doc/swagger/user_service.swagger.json
	swaggerJSONHandler := http.StripPrefix(swaggerJSONPath, http.FileServer(http.Dir("doc/swagger")))
	mux.Handle(swaggerJSONPath, swaggerJSONHandler)
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

func setupGRPCServer(creds credentials.TransportCredentials, ps repository.IRepository) (*grpc.Server, error) {
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
	userpb.RegisterUserServiceServer(grpcServer, service.NewUserServiceServer(ps))
	return grpcServer, nil
}

func setupGRPCGateway(serverAddr string, creds credentials.TransportCredentials) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	fmt.Println(serverAddr)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server")
	}
	defer conn.Close()

	err = userpb.RegisterUserServiceHandler(context.Background(), mux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gRPC gateway")
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
