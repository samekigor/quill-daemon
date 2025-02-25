package registry

import (
	"context"
	"fmt"

	"github.com/samekigor/quill-daemon/cmd/utils"
	"oras.land/oras-go/v2/registry/remote"

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

func (rd *RegistryDetails) GetImgRef() (imgRef string) {
	return fmt.Sprintf("%s/%s:%s", rd.Registry, rd.Repository, rd.Tag)
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

func (rd *RegistryDetails) PingRegistry(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			utils.ErrorLogger.Printf("Panic detected in IsPingRegistry: %v", r)
			err = fmt.Errorf("unexpected error: %v", r)
		}
	}()

	utils.InfoLogger.Printf("Attempting to ping registry: %s", rd.Registry)

	remoteRegistry, err := remote.NewRegistry(rd.Registry)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create remote registry: %v", err)
		return err
	}

	remoteRegistry.Client = &orasAuth.Client{
		Credential: func(ctx context.Context, registry string) (orasAuth.Credential, error) {
			return orasAuth.Credential{
				Username: rd.User.Username,
				Password: rd.User.Password,
			}, nil
		},
		Cache: orasAuth.NewCache(),
	}

	err = remoteRegistry.Ping(ctx)
	if err != nil {
		utils.WarnLogger.Printf("Failed to ping registry %s: %v", rd.Registry, err)
		return fmt.Errorf("registry ping failed: %v", err) // Return an error but don't panic
	}

	utils.InfoLogger.Printf("Successfully pinged registry: %s", rd.Registry)
	return nil
}
