package tooling

import "os"

// function 'EnvOrKey' returns the value of the environment variable with the given key,
// if the environment variable is unset, it returns the key name itself
func EnvOrKey(key string) string {
	val := os.Getenv(key)
	if val == "" {
		return "$" + key
	}
	return val
}

// function 'EnvOrDefault' returns the value of the environment variable with the given key,
// if the environment variable is unset, it returns the default value
func EnvOrDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
