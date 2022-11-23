package auconfigenv

import (
	"errors"
	"fmt"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	"regexp"
	"strconv"
)

func Validate() error {
	for _, k := range unknownConfigKeys {
		warnFunction(fmt.Sprintf("local configuration file contained setting for unknown configuration key %s", k))
	}

	var errorList = make([]error, 0)
	for _, it := range configItems {
		if it.Validate != nil {
			err := it.Validate(it.Key)
			if err != nil {
				warnFunction(fmt.Sprintf("failed to validate configuration field %s: %s", it.EnvName, err.Error()))
				errorList = append(errorList, err)
			}
		}
	}

	if len(errorList) > 0 {
		return fmt.Errorf("some configuration values failed to validate or parse. There were %d error(s). See details above", len(errorList))
	} else {
		return nil
	}
}

// --- generators for common validation functions ---

func ObtainPatternValidator(pattern string) auconfigapi.ConfigValidationFunc {
	return func(key string) error {
		value := Get(key)
		matched, err := regexp.MatchString(pattern, value)
		if err != nil {
			return err
		}

		if matched {
			return nil
		} else {
			return fmt.Errorf("must match %s", pattern)
		}
	}
}

func ObtainNotEmptyValidator() auconfigapi.ConfigValidationFunc {
	return func(key string) error {
		value := Get(key)
		if value == "" {
			return errors.New("must not be empty")
		} else {
			return nil
		}
	}
}

func ObtainUintRangeValidator(min uint, max uint) auconfigapi.ConfigValidationFunc {
	return func(key string) error {
		value := Get(key)
		vInt, err := AToUint(value)
		if err != nil {
			return err
		}

		if vInt < min || vInt > max {
			return fmt.Errorf("value %s is out of range [%d..%d]", value, min, max)
		}
		return nil
	}
}

func ObtainIntRangeValidator(min int, max int) auconfigapi.ConfigValidationFunc {
	return func(key string) error {
		value := Get(key)
		vInt, err := AToInt(value)
		if err != nil {
			return err
		}

		if vInt < min || vInt > max {
			return fmt.Errorf("value %s is out of range [%d..%d]", value, min, max)
		}
		return nil
	}
}

func ObtainIsBooleanValidator() auconfigapi.ConfigValidationFunc {
	return func(key string) error {
		value := Get(key)
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("value %s is not a valid boolean value", value)
		}
		return nil
	}
}

func ObtainIsRegexValidator() auconfigapi.ConfigValidationFunc {
	return func(key string) error {
		value := Get(key)
		if _, err := regexp.Compile(value); err != nil {
			return fmt.Errorf("value %s is not a valid regex pattern", value)
		}
		return nil
	}
}

func ObtainSingleCharacterValidator() auconfigapi.ConfigValidationFunc {
	return func(key string) error {
		value := Get(key)
		if len(value) < 1 {
			return fmt.Errorf("cannot be empty")
		} else if len(value) > 1 {
			return fmt.Errorf("cannot consist of multiple characters")
		}
		return nil
	}
}

// --- conversion helpers ---

func AToUint(s string) (uint, error) {
	vInt, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("value %s is not a valid integer: %s", s, err.Error())
	}
	if vInt < 0 {
		return 0, fmt.Errorf("value %s is negative", s)
	}
	return uint(vInt), nil
}

func AToInt(s string) (int, error) {
	vInt, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("value %s is not a valid integer: %s", s, err.Error())
	}
	return vInt, nil
}
