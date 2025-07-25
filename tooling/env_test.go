package tooling

import (
	"os"
	"testing"
)

func TestEnvOrKey_GivenEnvSet(t *testing.T) {
	os.Setenv("TEST_ENV", "test_env")
	defer os.Unsetenv("TEST_ENV")
	if EnvOrKey("TEST_ENV") != "test_env" {
		t.Errorf("func EnvOrKey returned %s, expected %s", EnvOrKey("TEST_ENV"), "test_env")
	}
}

func TestEnvOrKey_GivenEnvNotSet(t *testing.T) {
	if EnvOrKey("TEST_ENV") != "$TEST_ENV" {
		t.Errorf("func EnvOrKey returned %s, expected %s", EnvOrKey("TEST_ENV"), "$TEST_ENV")
	}
}

func TestEnvOrDefault_GivenEnvSet(t *testing.T) {
	os.Setenv("TEST_ENV", "test_env")
	defer os.Unsetenv("TEST_ENV")
	if EnvOrDefault("TEST_ENV", "default_value") != "test_env" {
		t.Errorf("func EnvOrDefault returned %s, expected %s", EnvOrDefault("TEST_ENV", "default_value"), "test_env")
	}
}

func TestEnvOrDefault_GivenEnvNotSet(t *testing.T) {
	if EnvOrDefault("TEST_ENV", "default_value") != "default_value" {
		t.Errorf("func EnvOrDefault returned %s, expected %s", EnvOrDefault("TEST_ENV", "default_value"), "default_value")
	}
}
