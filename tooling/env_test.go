package tooling

import (
	"os"
	"testing"
)

func TestEnvOrKey_GivenEnvSet(t *testing.T) {
	os.Setenv("TEST_ENVORKEY", "test_envorkey")
	defer os.Unsetenv("TEST_ENVORKEY")
	if EnvOrKey("TEST_ENVORKEY") != "test_envorkey" {
		t.Errorf("func EnvOrKey returned %s, expected %s", EnvOrKey("TEST_ENVORKEY"), "test_envorkey")
	}
}

func TestEnvOrKey_GivenEnvNotSet(t *testing.T) {
	if EnvOrKey("TEST_ENVORKEY") != "$TEST_ENVORKEY" {
		t.Errorf("func EnvOrKey returned %s, expected %s", EnvOrKey("TEST_ENVORKEY"), "$TEST_ENVORKEY")
	}
}

func TestEnvOrDefault_GivenEnvSet(t *testing.T) {
	os.Setenv("TEST_ENVORKEY", "test_envorkey")
	defer os.Unsetenv("TEST_ENVORKEY")
	if EnvOrDefault("TEST_ENVORKEY", "default_value") != "test_envorkey" {
		t.Errorf("func EnvOrDefault returned %s, expected %s", EnvOrDefault("TEST_ENVORKEY", "default_value"), "test_envorkey")
	}
}

func TestEnvOrDefault_GivenEnvNotSet(t *testing.T) {
	if EnvOrDefault("TEST_ENVORKEY", "default_value") != "default_value" {
		t.Errorf("func EnvOrDefault returned %s, expected %s", EnvOrDefault("TEST_ENVORKEY", "default_value"), "default_value")
	}
}
