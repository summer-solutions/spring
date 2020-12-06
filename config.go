package spring

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
)

const configName = "config"

type Config struct {
	*viper.Viper
}

func newViperConfig(appName, localConfigFolder string) (*Config, error) {
	viper.SetConfigName(configName)

	configFolder, hasConfigFolder := os.LookupEnv("SPRING_CONFIG_FOLDER")
	if !hasConfigFolder {
		configFolder = localConfigFolder
	}

	viper.SetConfigFile(configFolder + "/" + appName + "/config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viperConfig := &Config{
		viper.GetViper(),
	}

	err = viperConfig.loadEnvConfig()
	if err != nil {
		return nil, err
	}

	return viperConfig, nil
}

func (v *Config) loadEnvConfig() error {
	mainConfigFolderPath := v.GetMainPath()
	if _, err := os.Stat(mainConfigFolderPath + "/../.env.local"); os.IsNotExist(err) {
		return nil
	}

	err := godotenv.Load(mainConfigFolderPath + "/../.env.local")
	if err != nil {
		return err
	}

	for _, key := range v.Viper.AllKeys() {
		val, ok := v.Viper.Get(key).(string)
		if ok && strings.HasPrefix(v.Viper.Get(key).(string), "ENV[") {
			envKey := strings.TrimLeft(strings.TrimRight(val, "]"), "ENV[")
			envVal := os.Getenv(envKey)
			if envVal == "" {
				return errors.New("missing value for ENV variable " + envKey)
			}
			v.Viper.Set(key, os.Getenv(envKey))
		}
	}

	return nil
}

func (v *Config) GetMainPath() string {
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

func ServiceConfigDirectory(configDirectory string) *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "config_directory",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return configDirectory, nil
		},
	}
}

func serviceConfig() *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configDirectory := ctn.Get("config_directory").(string)
			return newViperConfig(ctn.Get("app").(*AppDefinition).Name, configDirectory)
		},
	}
}
