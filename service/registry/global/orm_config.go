package global

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/service/config"

	"gopkg.in/yaml.v2"

	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

var ormConfig orm.ValidatedRegistry

type RegistryInitFunc func(registry *orm.Registry)

func OrmConfigGlobalService(init RegistryInitFunc) spring.InitHandler {
	return func(s *spring.Server, def *spring.Def) {
		err := initOrmConfig(s, init, def)
		if err != nil {
			panic(err)
		}
	}
}

func initOrmRegistry(_ *spring.Server, configService *config.ViperConfig) (*orm.Registry, error) {
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

	config := parsedYaml["orm"].(map[interface{}]interface{})
	loadEnvConfig(config)

	for k, v := range config {
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

func initOrmConfig(s *spring.Server, init RegistryInitFunc, def *spring.Def) error {
	def.Name = "orm_config"
	def.Build = func(ctn di.Container) (interface{}, error) {
		configService := ctn.Get("config").(*config.ViperConfig)

		registry, err := initOrmRegistry(s, configService)
		if err != nil {
			return nil, err
		}
		ormConfig, err = registry.Validate()

		init(registry)
		return ormConfig, err
	}

	return nil
}
