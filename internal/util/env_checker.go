package util

import (
	"errors"
	"os"
)

type Envs = map[string]string

func GetEnvs(envs []string) (Envs, error) {
	var errs []string
	result := make(Envs)

	for _, envKey := range envs {
		env, ok := os.LookupEnv(envKey)
		if !ok {
			errs = append(errs, "Environment variable not set: "+envKey)
		}
		result[envKey] = env
	}

	if len(errs) > 0 {
		var errorString string
		for _, err := range errs {
			errorString += err + "\n"
		}
		return nil, errors.New(errorString)
	}

	return result, nil
}
