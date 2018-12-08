package env

import (
	"os"
	"strconv"
)

// Bool return a boolean env variable if present, else always false and nil error
func Bool(key string) (bool, error) {
	value, found := os.LookupEnv(key)
	if !found {
		return found, nil
	}

	return strconv.ParseBool(value)
}
