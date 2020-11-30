package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const configName = "config"

type ViperConfig struct {
	*viper.Viper
}

func NewViperConfig(localConfigFile string) (*ViperConfig, error) {
	// Consider using ConfigMap from kubernetes instead of file:
	// https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume
	// TODO: add hot reload for configMap
	// https://medium.com/@xcoulon/kubernetes-configmap-hot-reload-in-action-with-viper-d413128a1c9a

	viper.SetConfigName(configName)

	configFile, hasConfigFile := os.LookupEnv("SPRING_CONFIG_FILE")
	if hasConfigFile {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile(localConfigFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &ViperConfig{
		viper.GetViper(),
	}, nil
}

func (v *ViperConfig) GetMainPath() string {
	fileUsed := v.ConfigFileUsed()
	abs, _ := filepath.Abs(fileUsed)

	pathFragments := strings.Split(strings.TrimLeft(abs, "/"), "/")

	lastIndex := 0
	for i, fragment := range pathFragments {
		if fragment == "config" {
			lastIndex = i
			break
		}
	}

	var path string
	for j, fragment := range pathFragments {
		path += "/" + fragment
		if lastIndex == j {
			break
		}
	}

	return path
}
