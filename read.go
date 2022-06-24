package auconfigenv

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"regexp"
	"strings"
)

var LocalConfigFileName = "local-config.yaml"

// Read should be called in your code after Setup().
func Read() error {
	if err := ReadYaml(LocalConfigFileName); err != nil {
		return err
	}
	if err := readEnv(); err != nil {
		return err
	}
	return nil
}

// ReadYaml is exposed for direct use in testing code, allows you to ignore environment variables.
func ReadYaml(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// this is NOT an error
		return nil
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading local configuration yaml file %s: %s", filename, err.Error())
	}

	var stringValues = make(map[string]string, 0)
	err = yaml.UnmarshalStrict(yamlFile, &stringValues)
	if err != nil {
		return fmt.Errorf("error parsing local configuration flat yaml file %s (both keys and values must be strings): %s", filename, err.Error())
	}

	for _, it := range configItems {
		fileValue, ok := stringValues[it.Key]
		if ok {
			configValues[it.Key] = fileValue
			delete(stringValues, it.Key)
		}
	}

	for k, _ := range stringValues {
		return fmt.Errorf("local configuration file contained setting for unknown configuration key %s, bailing out", k)
	}

	return nil
}

func readEnv() error {
	re := regexp.MustCompile(`[^a-z0-9]`)
	for _, it := range configItems {
		// simply fill in EnvName if unset
		if it.EnvName == "" {
			it.EnvName = "CONFIG_" + strings.ToUpper(re.ReplaceAllString(it.Key, "_"))
		}

		envValue, ok := os.LookupEnv(it.EnvName)
		if ok {
			configValues[it.Key] = envValue
		}
	}

	return nil
}
