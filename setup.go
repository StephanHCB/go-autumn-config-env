package auconfigenv

import (
	"fmt"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	"log"
)

var warnFunction auconfigapi.ConfigWarnFunc = warn

var configItems []auconfigapi.ConfigItem

var unknownConfigKeys []string

// Setup initializes configuration with the default values - you need to call this from your code before calling
// Read() and Validate().
func Setup(items []auconfigapi.ConfigItem, warnFunc auconfigapi.ConfigWarnFunc) error {
	configItems = items
	configValues = make(map[string]string, 0)
	warnFunction = warnFunc

	for _, it := range configItems {
		defaultStr, ok := it.Default.(string)
		if !ok {
			return fmt.Errorf("error parsing default value for key %s - this library only supports strings", it.Key)
		}
		configValues[it.Key] = defaultStr
	}

	return nil
}

func warn(message string) {
	log.Print(message)
}
