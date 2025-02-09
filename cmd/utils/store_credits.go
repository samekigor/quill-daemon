package utils

import (
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	creditsFilePath = "/etc/quill/"
	creditsFileName = "credentials.yml"
)

var (
	creditsFullPath = fmt.Sprintf("%s%s", creditsFilePath, creditsFileName)
)

type RegistryEntry struct {
	Registry string `yaml:"registry"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (r *RegistryEntry) EncodePassword() error {
	encoded := base64.StdEncoding.EncodeToString([]byte(r.Password))
	if encoded == "" {
		return fmt.Errorf("failed to encode password")
	}
	r.Password = encoded
	return nil
}

func (r *RegistryEntry) DecodePassword() error {
	decoded, err := base64.StdEncoding.DecodeString(r.Password)
	if err != nil {
		return err
	}
	r.Password = string(decoded)
	return nil
}

func saveToYAML(entries []RegistryEntry) error {
	data, err := yaml.Marshal(entries)
	if err != nil {
		return err
	}
	return os.WriteFile(creditsFullPath, data, 0666)
}
func loadFromYAML() ([]RegistryEntry, error) {
	data, err := os.ReadFile(creditsFullPath)
	if err != nil {
		return nil, err
	}

	var entries []RegistryEntry
	err = yaml.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *RegistryEntry) AddRegistryEntry() error {
	entries, err := loadFromYAML()
	if err != nil {
		return err
	}

	entries = append(entries, *r)
	return saveToYAML(entries)
}

func (r *RegistryEntry) RemoveRegistryEntry() error {
	entries, err := loadFromYAML()
	if err != nil {
		return err
	}

	for i, entry := range entries {
		if entry.Registry == r.Registry {
			entries = append(entries[:i], entries[i+1:]...)
			break
		}
	}

	return saveToYAML(entries)
}

func GetRegistryEntryByName(registryName string) (*RegistryEntry, error) {
	entries, err := loadFromYAML()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("registry entry not found: %s", registryName)
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.Registry == registryName {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("registry entry not found: %s", registryName)
}
