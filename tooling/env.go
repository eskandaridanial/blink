package tooling

import (
	"os"
)

func EnvOrKey(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return "$" + key
}

func EnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
