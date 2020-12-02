package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sarulabs/di"
	diLocal "github.com/summer-solutions/spring/di"

	"github.com/summer-solutions/spring/services/config"

	"gopkg.in/yaml.v2"

	"github.com/summer-solutions/orm"
)

var ormConfig orm.ValidatedRegistry

type RegistryInitFunc func(registry *orm.Registry)

func OrmRegistry(init RegistryInitFunc) *diLocal.ServiceDefinition {
	return &diLocal.ServiceDefinition{
		Name:   "orm_config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService, err := ctn.SafeGet("config")
			if err != nil {
				return nil, fmt.Errorf("missing config service")
			}

			registry, err := initOrmRegistry(configService.(*config.ViperConfig))
			if err != nil {
				return nil, err
			}

			init(registry)

			ormConfig, err = registry.Validate()

			return ormConfig, err
		},
	}
}

func initOrmRegistry(configService *config.ViperConfig) (*orm.Registry, error) {
	var yamlFileData []byte
	var err error

	if os.Getenv("ORM_CONFIG_FILE") != "" {
		yamlFileData, err = ioutil.ReadFile(os.Getenv("ORM_CONFIG_FILE"))
	} else {
		yamlFileData, err = ioutil.ReadFile(configService.GetMainPath() + "/orm/config.yaml")
	}

	if err != nil {
		return nil, err
	}

	var parsedYaml map[string]interface{}

	err = yaml.Unmarshal(yamlFileData, &parsedYaml)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})

	conf := parsedYaml["orm"].(map[interface{}]interface{})
	loadEnvConfig(conf)

	for k, v := range conf {
		data[k.(string)] = v
	}

	return orm.InitByYaml(data), nil
}

func loadEnvConfig(configData map[interface{}]interface{}) {
	for k, v := range configData {
		_, isString := v.(string)
		_, isInt := v.(int)

		if !isString && !isInt {
			loadEnvConfig(v.(map[interface{}]interface{}))
		} else if isString {
			if strings.HasPrefix(v.(string), "ENV[") {
				envKey := strings.TrimLeft(strings.TrimRight(v.(string), "]"), "ENV[")
				envVal := os.Getenv(envKey)
				if envVal == "" {
					panic("missing value for ENV variable " + v.(string))
				}
				configData[k] = envVal
			}
		}
	}
}