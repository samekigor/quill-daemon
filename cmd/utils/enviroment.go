package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

var PrefixEnviromentVariables = "quill"

func InitEnviromentVariables() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(PrefixEnviromentVariables)
}

func GetEnviromentVariable(name string) (envVar string, err error) {
	envVar = viper.GetString(name)
	if envVar == "" {
		return "", fmt.Errorf("failure with reading environment variable: %s", name)
	}
	return envVar, nil
}
