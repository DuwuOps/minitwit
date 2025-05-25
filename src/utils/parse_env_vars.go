package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetEnvDuration(key, def string) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid %s: %v", key, err)
		}
		return d
	}
	// fall back to default if unset
	d, err := time.ParseDuration(def)
	if err != nil {
		log.Fatalf("invalid default %s for %s: %v", def, key, err)
	}
	return d
}

func GetEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid %s: %v", key, err)
		}
		return i
	}
	return def
}