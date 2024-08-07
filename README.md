# go-autumn-config-env

## About go-autumn-config-env

A library that handles configuration for enterprise microservices.

This is a **lightweight alternative** to [go-autumn-config](https://github.com/StephanHCB/go-autumn-config) that
is tailored specifically towards 15 factor microservices, especially when run in a container.

In order of precedence, configuration values come from
- environment variables
- an optional flat configuration file (intended for use in local development only)
- hardcoded defaults

There is **no support** for structured configuration values. Just like environment variables, configuration
values are strings only.

There is **no support** for command line switches.

There is **no support** for application profiles.

There is **no support** for a separate file for secrets.

On the plus side, this guides developers on the right path towards a 15 factor app, and there is
a minimal dependency and runtime footprint.

## How to use

We recommend collecting all configuration related code in a package `internal/repository/configuration`.

You configure the configuration subsystem by a call to `auconfigenv.Setup(...)`. This function takes 2 arguments:
 - a list of `auconfigapi.ConfigItem` to specify what configuration items exist 
 - a warning message handler of type `auconfigapi.ConfigWarnFunc`, which should probably log a warning
   using your preferred logging framework. `log.Print` satisfies the type requirements, but again we
   hope this is not what you'll use in production...

See [go-autumn-config-api](https://github.com/StephanHCB/go-autumn-config-api/blob/master/api.go) for 
the precise type definitions.

When you request your configuration to be loaded, which you must do yourself with a call to 
`auconfigenv.Read()`, every key is assigned its value by going through the following precedence list:
 - environment variable
 - configuration read from local-config.yaml
 - default value specified in ConfigItems

Once loaded, validate the values by calling `auconfigenv.Validate()` to validate each configuration entry.
It will use the warning message callback to notify you about individual errors, and if any validation
errors were found, return an error at the end.

Once loaded (even before validation), access configuration values by calling `auconfigenv.Get(key)`,
which returns the cached string value. 

We also export `auconfigenv.Set(key, value)`, which allows injecting configuration values at runtime, for example
after querying a secrets management solution such as OpenBAO, or during tests.

## Structured data under a key

... is not directly supported, but you can parse JSON from configuration values, and the validation function
can check validity and correctness.

## Examples:

Using this can be as simple as:

```go
package configuration

import (
	"fmt"
	"github.com/StephanHCB/go-autumn-config-api"
	"github.com/StephanHCB/go-autumn-config-env"
	"log"
    "strconv"
)

// custom validation function example
func checkValidPortNumber(key string) error {
    portStr := auconfigenv.Get(key)
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return fmt.Errorf("fatal error: configuration value for key %s is not a valid number: %s", key, err.Error())
    }
    if port < 1024 || port > 65535 {
        return fmt.Errorf("fatal error: configuration value for key %s is not in range 1024..65535", key)
    }
    return nil
}

// define what configuration items you want
var configItems = []auconfigapi.ConfigItem{
	{
		Key:         "server.address",
		Default:     "",
		Description: "ip address or hostname to listen on, can be left blank for localhost",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         "server.port",
		Default:     "8080",
		Description: "port to listen on, defaults to 8080 if not set",
		Validate:    checkValidPortNumber,
	},
}

// initialize the library.
func Setup() {
    _ = auconfigenv.Setup(configItems, log.Print)
	_ = auconfigenv.Read()
	_ = auconfigenv.Validate()
}

// provide accessor functions

func ServerAddress() string {
	return fmt.Sprintf("%s:%s", auconfigenv.Get("server.address"), auconfigenv.Get("server.port"))
}
```
