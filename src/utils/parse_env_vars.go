package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetEnvDuration(key, def string) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			logMsg := fmt.Sprintf("invalid %s", key)
			LogFatal(logMsg, err)
		}
		return d
	}
	// fall back to default if unset
	d, err := time.ParseDuration(def)
	if err != nil {
		logMsg := fmt.Sprintf("invalid default %s for %s", def, key)
		LogFatal(logMsg, err)
	}
	return d
}

func GetEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			logMsg := fmt.Sprintf("invalid %s", key)
			LogFatal(logMsg, err)
		}
		return i
	}
	return def
}
