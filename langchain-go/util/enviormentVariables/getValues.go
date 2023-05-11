package enviormentVariables

import (
	"errors"
	"os"
)

func GetFromDictOrEnv(data map[string]interface{}, key string, envKey string, defaultValue *string) (string, error) {
	if val, ok := data[key]; ok && val != nil {
		return val.(string), nil
	} else if val, ok := os.LookupEnv(envKey); ok {
		return val, nil
	} else if defaultValue != nil {
		return *defaultValue, nil
	} else {
		return "", errors.New("Did not find " + key + ", please add an environment variable " + envKey + " which contains it, or pass " + key + " as a named parameter.")
	}
}
