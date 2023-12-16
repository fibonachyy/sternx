// internal/service/user_test.go
package service

import (
	"context"
	"net"

	"testing"

	"github.com/fibonachyy/sternx/internal/api/user"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestCreateUserEndpoint(t *testing.T) {
	// Create a listener to simulate a gRPC connection
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	// Register your gRPC service implementation
	userService := &UserServiceServer{}
	user.RegisterUserServiceServer(server, userService)

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Fatalf("Server exited with error: %v", err)
		}
	}()

	// Create a gRPC client connection to the simulated server
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := user.NewUserServiceClient(conn)

	// Test the CreateUser endpoint
	req := &user.UserRequest{
		Name:  "John",
		Email: "john@example.com",
	}

	res, err := client.CreateUser(ctx, req)
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, res, "Expected a non-nil response")
	assert.NotEmpty(t, res.UserId, "Expected a non-empty user ID")
	assert.Equal(t, req.Name, res.Name, "Expected names to match")
	assert.Equal(t, req.Email, res.Email, "Expected emails to match")
}

// Implement similar test functions for other methods (GetUser, UpdateUser, DeleteUser)
