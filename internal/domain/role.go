package domain

import (
	"strings"

	userpb "github.com/fibonachyy/sternx/internal/api"
)

const (
	StandardRole = "standard"
	AdminRole    = "admin"
)

func StringToRole(roleStr string) userpb.Role {
	switch strings.ToLower(roleStr) {
	case "admin":
		return userpb.Role_ADMIN
	case "standard":
		return userpb.Role_ADMIN
	default:
		return userpb.Role_ADMIN
	}
}
