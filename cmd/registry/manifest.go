package registry

import (
	"context"
	"encoding/json"
	"fmt"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type Manifest struct {
	SchemaVersion int `json:"schemaVersion"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

func FetchManifest(rd *RegistryDetails) (*Manifest, error) {
	ctx := context.Background()

	repo, err := remote.NewRepository(imageRef)
	if err != nil {
		return nil, fmt.Errorf("connection error with registry: %w", err)
	}

	repo.Client = &auth.Client{
		Credential: func(ctx context.Context, host string) (auth.Credential, error) {
			fmt.Println("🔑 Uwierzytelnianie dla hosta:", host)
			return auth.Credential{
				Username: username,
				Password: password,
			}, nil
		},
	}

	// Pobranie manifestu
	descriptor, manifestBytes, err := oras.FetchBytes(ctx, repo, imageRef, oras.DefaultFetchOptions)
	if err != nil {
		return nil, fmt.Errorf("błąd pobierania manifestu: %w", err)
	}

	fmt.Println("✅ Pobrany manifest dla obrazu:", imageRef)
	fmt.Println("📜 MediaType manifestu:", descriptor.MediaType)

	// Parsowanie JSON do struktury
	var manifest Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("błąd parsowania manifestu: %w", err)
	}

	return &manifest, nil
}
