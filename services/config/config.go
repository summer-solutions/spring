package config

type Config interface {
	GetInt(key string) int
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	Get(key string) interface{}
	GetStringSlice(key string) []string
	GetBool(key string) bool
	GetMainPath() string
	ConfigFileUsed() string
}
