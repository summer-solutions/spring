package global

import (
	"io/ioutil"
	"os"

	"github.com/summer-solutions/spring"

	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

var ormConfig orm.ValidatedRegistry

type RegistryInitFunc func(registry *orm.Registry)

func OrmConfigGlobalService(init RegistryInitFunc) spring.InitHandler {
	return func(s *spring.Server, def *spring.Def) {
		registry, err := initOrmRegistry(s)
		if err != nil {
			panic(err)
		}

		init(registry)

		err = initOrmConfig(registry, def)
		if err != nil {
			panic(err)
		}
		if s.IsInLocalMode() {
			err = checkForAlters(ormConfig)
			if err != nil {
				panic(err)
			}
		}
	}
}

func initOrmRegistry(_ *spring.Server) (*orm.Registry, error) {
	var yamlFileData []byte
	var err error

	if os.Getenv("ORM_CONFIG_FILE") != "" {
		yamlFileData, err = ioutil.ReadFile(os.Getenv("ORM_CONFIG_FILE"))
	} else {
		yamlFileData, err = ioutil.ReadFile("../../config/orm/config.yaml")
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
	for k, v := range parsedYaml["orm"].(map[interface{}]interface{}) {
		data[k.(string)] = v
	}

	return orm.InitByYaml(data), nil
}

func initOrmConfig(registry *orm.Registry, def *spring.Def) error {
	var err error

	ormConfig, err = registry.Validate()

	if err != nil {
		return err
	}

	def.Name = "orm_config"
	def.Build = func(ctn di.Container) (interface{}, error) {
		return ormConfig, nil
	}

	return nil
}

func checkForAlters(ormConfig orm.ValidatedRegistry) error {
	engine := ormConfig.CreateEngine()

	alters := engine.GetAlters()
	for _, alter := range alters {
		if alter.Safe {
			color.Green("%s\n\n", alter.SQL)
		} else {
			color.Red("%s\n\n", alter.SQL)
		}
	}

	return nil
}
