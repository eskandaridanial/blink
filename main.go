package main

import (
	"github.com/eskandaridanial/blink/formatters"
	"github.com/eskandaridanial/blink/handlers"
	"github.com/eskandaridanial/blink/hooks"
	"github.com/eskandaridanial/blink/logging"
	"github.com/eskandaridanial/blink/models"
)

func main() {
	logger := logging.NewLogger(
		logging.WithLevel(models.Debug),
		logging.WithField(models.String("key", "value")),
		logging.WithHandler(handlers.NewConsoleHandler(&formatters.TextFormatter{})),
		logging.WithHook(&hooks.DefaultHook{}),
	)

	logger.Debug("hello world", models.Bool("bool", true), models.Int("number", 1))
	logger.Info("hello world", models.Bool("bool", false), models.Int("number", 2))
}
