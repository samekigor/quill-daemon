package registry

import (
	"context"
	"fmt"

	"github.com/samekigor/quill-daemon/cmd/utils"
	"oras.land/oras-go/v2/registry/remote"

	// "github.com/oras-project/oras-go/v2"
	// "github.com/oras-project/oras-go/v2/content/file"
	// "github.com/oras-project/oras-go/v2/registry/remote"
	orasAuth "oras.land/oras-go/v2/registry/remote/auth"
)

type UserDetails struct {
	Username string
	Password string
}

type RegistryDetails struct {
	Registry   string
	Repository string
	Tag        string
	User       UserDetails
	KeyInStore string
}

func NewRegistryDetails(registry, repository, tag, username, password string) *RegistryDetails {
	return &RegistryDetails{
		Registry:   registry,
		Repository: repository,
		Tag:        tag,
		User: UserDetails{
			Username: username,
			Password: password,
		},
		KeyInStore: fmt.Sprintf("%s:%s", registry, username),
	}
}

// IsPingRegistry checks whether the registry is reachable and returns the status.
func (rd *RegistryDetails) IsPingRegistry(ctx context.Context) (pinged bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			utils.ErrorLogger.Printf("Panic detected in IsPingRegistry: %v", r)
			err = fmt.Errorf("unexpected error: %v", r)
			pinged = false
		}
	}()

	utils.InfoLogger.Printf("Attempting to ping registry: %s", rd.Registry)

	// Create remote registry instance
	remoteRegistry, err := remote.NewRegistry(rd.Registry)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create remote registry: %v", err)
		return false, err
	}

	// Configure authentication for the registry client
	remoteRegistry.Client = &orasAuth.Client{
		Credential: func(ctx context.Context, registry string) (orasAuth.Credential, error) {
			return orasAuth.Credential{
				Username: rd.User.Username,
				Password: rd.User.Password,
			}, nil
		},
		Cache: orasAuth.NewCache(),
	}

	// Perform registry ping check
	err = remoteRegistry.Ping(ctx)
	if err != nil {
		utils.WarnLogger.Printf("Failed to ping registry %s: %v", rd.Registry, err)
		return false, fmt.Errorf("registry ping failed: %v", err) // Return an error but don't panic
	}

	utils.InfoLogger.Printf("Successfully pinged registry: %s", rd.Registry)
	return true, nil
}
