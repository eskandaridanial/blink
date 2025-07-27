package main

import (
	"os"
	"time"

	"github.com/eskandaridanial/blink/formatters"
	"github.com/eskandaridanial/blink/models"
)

func main() {

	record := models.Record{
		Level:       models.Debug,
		Message:     "hello world",
		Caller:      "main.main",
		ReferenceId: "1234567890",
		Fields:      []models.Field{models.Bool("bool", true), models.Int("number", 1)},
		Timestamp:   time.Now(),
	}

	// formatter
	jsonFormatter := formatters.NewJsonFormatter()
	jsonFormatter.Format(os.Stdout, record)

	jsonFormatterWithConfig := formatters.NewJsonFormatterWithConfig(formatters.Config{
		JsonCompact:         false,
		JsonEscapeHtml:      true,
		EnablePanicRecovery: true,
		IncludeEmptyFields:  true,
		PayloadPoolSize:     8,
		WriterBufferSize:    4096,
		FormatBufferSize:    1024,
		TimeFormat:          "2006-01-02 15:04:05",
	})
	jsonFormatterWithConfig.Format(os.Stdout, record)

	textFormatter := formatters.NewTextFormatter()
	textFormatter.Format(os.Stdout, record)

	textFormatterWithPattern := formatters.NewTextFormatterWithPattern("%{timestamp} %{level} - %{message}")
	textFormatterWithPattern.Format(os.Stdout, record)

	textFormatterWithConfig := formatters.NewTextFormatterWithConfig(formatters.Config{
		JsonCompact:         false,
		JsonEscapeHtml:      true,
		EnablePanicRecovery: true,
		IncludeEmptyFields:  true,
		PayloadPoolSize:     8,
		WriterBufferSize:    4096,
		FormatBufferSize:    1024,
		TimeFormat:          "2006-01-02 15:04:05",
		TextPattern:         "%{timestamp} %{level} - %{message}",
	})
	textFormatterWithConfig.Format(os.Stdout, record)

	configBuilder := formatters.NewConfigBuilder()
	config := configBuilder.
		JsonEscapeHtml(true).
		JsonCompact(false).
		IncludeEmptyFields(true).
		PayloadPoolSize(8).
		WriterBufferSize(4096).
		FormatBufferSize(1024).
		TimeFormat("2006-01-02 15:04:05").
		TextPattern("%{timestamp} %{level} - %{message}").
		EnablePanicRecovery(true).
		Build()
	jsonFormatterWithBuilderConfig := formatters.NewJsonFormatterWithConfig(config)
	jsonFormatterWithBuilderConfig.Format(os.Stdout, record)
}
