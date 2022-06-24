package auconfigenv

var configValues map[string]string

func Get(key string) string {
	// TODO deal with unknown keys better
	return configValues[key]
}
