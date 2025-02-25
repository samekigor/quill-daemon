package server

import (
	"context"
	"fmt"
	"log"

	"github.com/samekigor/quill-daemon/cmd/registry"
	auth "github.com/samekigor/quill-daemon/proto/auth"
)

type AuthServer struct {
	auth.UnimplementedAuthServer
}

func (s *AuthServer) LoginToRegistry(ctx context.Context, req *auth.LoginRequest) (*auth.LoginStatus, error) {

	log.Printf("Received Login request for registry: %s, username: %s", req.Registry, req.Username)
	re, err := registry.GetRegistryEntryByName(req.Registry)
	defer func() { re.Password = "" }()

	if err != nil {
		return &auth.LoginStatus{
			IsSuccess: false,
			Message:   fmt.Sprintf("Failed to get registry entry for %s", req.Registry),
		}, err
	}

	err = re.DecodePassword()
	if err != nil {
		return &auth.LoginStatus{
			IsSuccess: false,
			Message:   fmt.Sprintf("Failed to decode password %s", req.Registry),
		}, err
	}

	rd := registry.NewRegistryDetails(re.Registry, "", "", req.Username, re.Password)

	if err = rd.PingRegistry(ctx); err != nil {
		return &auth.LoginStatus{
			IsSuccess: false,
			Message:   fmt.Sprintf("Failed to ping registry %s: %v", req.Registry, err),
		}, err
	} else {
		return &auth.LoginStatus{
			IsSuccess: true,
			Message:   fmt.Sprintf("User %s logged in successfully to registry %s\n", req.Username, req.Registry),
		}, nil
	}

}

func (s *AuthServer) LogoutFromRegistry(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutStatus, error) {
	log.Printf("Received Logout request for registry: %s", req.Registry)

	re, err := registry.GetRegistryEntryByName(req.Registry)

	if err != nil {
		return &auth.LogoutStatus{
			IsSuccess: false,
			Message:   fmt.Sprintf("Failed to get registry entry for %s", req.Registry),
		}, err
	}

	err = re.RemoveRegistryEntry()

	if err != nil {
		return &auth.LogoutStatus{
			IsSuccess: false,
			Message:   fmt.Sprintf("Failed to log out from registry %s", req.Registry),
		}, err
	} else {
		return &auth.LogoutStatus{
			IsSuccess: true,
			Message:   fmt.Sprintf("Logged out successfully from registry %s", req.Registry),
		}, nil
	}

}
