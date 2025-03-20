package helpers

import "os"

func GetEnvVar(key string, default_value string) string {
	// If key exists in environment variables, return the value of it
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	// Else return the given default_value
	return default_value
}
