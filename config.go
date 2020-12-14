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

func ServiceProviderConfigDirectory(configDirectory string) *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "config_directory",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return configDirectory, nil
		},
	}
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
	mainConfigFolderPath := v.getMainPath()
	if _, err := os.Stat(mainConfigFolderPath + "/../.env.local"); !os.IsNotExist(err) {
		err := godotenv.Load(mainConfigFolderPath + "/../.env.local")
		if err != nil {
			return err
		}
	}

	for _, key := range v.Viper.AllKeys() {
		val, ok := v.Viper.Get(key).(string)
		if ok && strings.HasPrefix(v.Viper.Get(key).(string), "ENV[") {
			envKey := strings.TrimLeft(strings.TrimRight(val, "]"), "ENV[")
			_, has := os.LookupEnv(envKey)
			if !has {
				return errors.New("missing value for ENV variable " + envKey)
			}
			s := strings.Split(key, ".")
			if len(s) > 1 {
				subKey := strings.Join(s[0:len(s)-1], ".")
				subVal := v.Viper.Get(subKey).(map[string]interface{})
				subVal[s[len(s)-1]] = os.Getenv(envKey)
			} else {
				v.Viper.Set(key, os.Getenv(envKey))
			}
		}
	}

	return nil
}

func (v *Config) getMainPath() string {
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

func serviceConfig() *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configDirectory := ctn.Get("config_directory").(string)
			return newViperConfig(DIC().App().Name(), configDirectory)
		},
	}
}
