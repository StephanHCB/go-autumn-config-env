package auconfigenv

var configValues map[string]string

func Get(key string) string {
	return configValues[key]
}
