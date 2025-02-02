package utils

import (
	"fmt"

	"github.com/zalando/go-keyring"
	"github.com/zalando/go-keyring/backend/secretservice"
)

// Ensure Daemon reads from the same "login" collection as CLI
func init() {
	keyring.DefaultKeyring = secretservice.NewWithOptions(&secretservice.Options{
		CollectionName: "login", // Must match CLI
	})
}
func GetPassword(service string, registry string, user string) (pwd string, err error) {
	key := fmt.Sprintf("%s:%s", registry, user)
	fmt.Print(key)
	pwd, err = keyring.Get(service, key)
	if err != nil {
		return "", err
	}
	return pwd, nil
}
