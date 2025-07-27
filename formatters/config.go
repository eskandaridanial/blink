package formatters

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration options for the formatters package.
// It provides comprehensive control over buffer sizes, formatting options,
// and safety features to optimize performance for different use cases.
//
// All configuration values can be set via environment variables with
// the BLINK_ prefix, or programmatically using the ConfigBuilder.
//
// Buffer size recommendations:
// - WriterBufferSize: 4KB for general use, 8KB+ for high-throughput
// - FormatBufferSize: 1KB for typical logs, 2KB+ for complex field structures
// - PayloadPoolSize: 8 for moderate concurrency, 16+ for high concurrency
type Config struct {
	// WriterBufferSize controls the size of the internal buffered writers.
	// Larger buffers reduce system calls but use more memory.
	// Environment variable: BLINK_WRITER_BUFFER_SIZE
	// Default: 4096 bytes (4KB)
	// Minimum: 1 byte
	WriterBufferSize int

	// FormatBufferSize sets the initial capacity for formatting buffers.
	// This should accommodate typical log message sizes to avoid reallocations.
	// Environment variable: BLINK_FORMAT_BUFFER_SIZE
	// Default: 1024 bytes (1KB)
	// Minimum: 1 byte
	FormatBufferSize int

	// PayloadPoolSize determines the initial capacity of Json payload maps.
	// Higher values reduce map reallocations for logs with many fields.
	// Environment variable: BLINK_PAYLOAD_POOL_SIZE
	// Default: 8 fields
	// Minimum: 1 field
	PayloadPoolSize int

	// TimeFormat specifies the Go time layout for timestamp formatting.
	// Uses Go's reference time: Mon Jan 2 15:04:05 MST 2006 (Unix: 1136239445)
	// Environment variable: BLINK_TIME_FORMAT
	// Default: "2006-01-02T15:04:05.000Z07:00" (ISO 8601 with milliseconds)
	// Must not be empty
	TimeFormat string

	// TextPattern defines the output format for text formatter using placeholders.
	// Available placeholders:
	//   %{timestamp} - Formatted timestamp
	//   %{level}     - Log level (DEBUG, INFO, WARN, ERROR, etc.)
	//   %{referenceId} - Optional reference/correlation ID
	//   %{caller}    - Optional caller information (file:line)
	//   %{fields}    - Key-value pairs formatted as "key=value key2=value2"
	//   %{message}   - The actual log message
	// Environment variable: BLINK_TEXT_PATTERN
	// Default: "%{timestamp} %{level} %{referenceId} %{caller} [%{fields}] - %{message}"
	// Must not be empty
	TextPattern string

	// JsonEscapeHtml controls whether Json strings are Html-escaped.
	// When true, characters like <, >, & are escaped as \u003c, \u003e, \u0026.
	// Set to true if logs will be embedded in Html contexts.
	// Environment variable: BLINK_JSON_ESCAPE_HTML
	// Default: false (no Html escaping)
	JsonEscapeHtml bool

	// JsonCompact determines Json output formatting.
	// When true, produces single-line compact Json.
	// When false, produces pretty-printed Json with indentation.
	// Environment variable: BLINK_JSON_COMPACT
	// Default: true (compact output)
	JsonCompact bool

	// IncludeEmptyFields controls whether empty/zero-value fields are included.
	// When true, fields like empty strings or nil values are included in output.
	// When false, only non-empty fields are included, reducing output size.
	// Environment variable: BLINK_INCLUDE_EMPTY_FIELDS
	// Default: false (exclude empty fields)
	IncludeEmptyFields bool

	// EnablePanicRecovery enables automatic recovery from panics during formatting.
	// When true, panics are caught and a fallback format is used.
	// When false, panics propagate up to the caller.
	// Disable only in development/testing environments.
	// Environment variable: BLINK_ENABLE_PANIC_RECOVERY
	// Default: true (panic recovery enabled)
	EnablePanicRecovery bool
}

// ConfigBuilder provides a fluent interface for constructing configuration objects.
// It implements the builder pattern to allow method chaining and provides
// validation to ensure only valid values are set.
//
// Example usage:
//
//	config := NewConfigBuilder().
//	    WriterBufferSize(8192).
//	    TimeFormat("2006-01-02 15:04:05").
//	    JsonCompact(false).
//	    Build()
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new ConfigBuilder initialized with default values.
// This ensures that all fields have sensible defaults even if not explicitly set.
//
// Returns:
//
//	*ConfigBuilder: A new builder instance ready for configuration
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: LoadConfig(),
	}
}

// WriterBufferSize sets the buffer size for internal buffered writers.
// Invalid values (≤ 0) are ignored to prevent runtime errors.
//
// Parameters:
//
//	size int: Buffer size in bytes, must be > 0
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) WriterBufferSize(size int) *ConfigBuilder {
	if size > 0 {
		b.config.WriterBufferSize = size
	}
	return b
}

// FormatBufferSize sets the initial capacity for formatting buffers.
// Invalid values (≤ 0) are ignored to prevent runtime errors.
//
// Parameters:
//
//	size int: Buffer capacity in bytes, must be > 0
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) FormatBufferSize(size int) *ConfigBuilder {
	if size > 0 {
		b.config.FormatBufferSize = size
	}
	return b
}

// PayloadPoolSize sets the initial capacity for Json payload maps.
// Invalid values (≤ 0) are ignored to prevent runtime errors.
//
// Parameters:
//
//	size int: Map capacity in number of fields, must be > 0
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) PayloadPoolSize(size int) *ConfigBuilder {
	if size > 0 {
		b.config.PayloadPoolSize = size
	}
	return b
}

// TimeFormat sets the Go time layout string for timestamp formatting.
// Empty strings are ignored to prevent invalid time formats.
//
// Parameters:
//
//	format string: Go time layout string, must not be empty
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) TimeFormat(format string) *ConfigBuilder {
	if format != "" {
		b.config.TimeFormat = format
	}
	return b
}

// TextPattern sets the formatting pattern for text output.
// Empty strings are ignored to prevent invalid patterns.
//
// Parameters:
//
//	pattern string: Text formatting pattern with placeholders, must not be empty
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) TextPattern(pattern string) *ConfigBuilder {
	if pattern != "" {
		b.config.TextPattern = pattern
	}
	return b
}

// JsonEscapeHtml configures Html escaping for Json output.
// This setting affects all string values in Json output.
//
// Parameters:
//
//	escape bool: true to enable Html escaping, false to disable
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) JsonEscapeHtml(escape bool) *ConfigBuilder {
	b.config.JsonEscapeHtml = escape
	return b
}

// JsonCompact configures Json output formatting style.
// Compact mode produces single-line output, non-compact produces pretty-printed output.
//
// Parameters:
//
//	compact bool: true for single-line output, false for pretty-printed
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) JsonCompact(compact bool) *ConfigBuilder {
	b.config.JsonCompact = compact
	return b
}

// IncludeEmptyFields configures whether empty/zero-value fields are included.
// This affects both Json and text output formats.
//
// Parameters:
//
//	include bool: true to include empty fields, false to exclude them
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) IncludeEmptyFields(include bool) *ConfigBuilder {
	b.config.IncludeEmptyFields = include
	return b
}

// EnablePanicRecovery configures automatic panic recovery during formatting.
// When enabled, panics are caught and a fallback format is used.
//
// Parameters:
//
//	enable bool: true to enable panic recovery, false to disable
//
// Returns:
//
//	*ConfigBuilder: The builder instance for method chaining
func (b *ConfigBuilder) EnablePanicRecovery(enable bool) *ConfigBuilder {
	b.config.EnablePanicRecovery = enable
	return b
}

// Build finalizes the configuration and returns the built Config struct.
// The returned configuration is a copy and can be safely modified.
//
// Returns:
//
//	Config: The finalized configuration struct
func (b *ConfigBuilder) Build() Config {
	b.config.validate()
	return b.config
}

// DefaultConfig returns a configuration with all default values.
// This is equivalent to creating a new ConfigBuilder and calling Build().
//
// Returns:
//
//	Config: Configuration struct with default values
func DefaultConfig() Config {
	return LoadConfig()
}

// LoadConfig loads configuration from environment variables with fallback to defaults.
// Environment variables are parsed with type validation and invalid values fall back to defaults.
// This function is safe to call multiple times and will always return a valid configuration.
//
// Environment variables read:
//
//	BLINK_WRITER_BUFFER_SIZE    - Writer buffer size in bytes
//	BLINK_FORMAT_BUFFER_SIZE    - Format buffer size in bytes
//	BLINK_PAYLOAD_POOL_SIZE     - Payload pool initial capacity
//	BLINK_TIME_FORMAT           - Go time layout string
//	BLINK_TEXT_PATTERN          - Text formatter pattern
//	BLINK_JSON_ESCAPE_HTML      - Boolean for Html escaping
//	BLINK_JSON_COMPACT          - Boolean for compact Json
//	BLINK_INCLUDE_EMPTY_FIELDS  - Boolean for empty field inclusion
//	BLINK_ENABLE_PANIC_RECOVERY - Boolean for panic recovery
//
// Returns:
//
//	Config: Configuration loaded from environment with defaults as fallback
func LoadConfig() Config {
	return Config{
		WriterBufferSize:    getEnvInt("BLINK_WRITER_BUFFER_SIZE", 4096),
		FormatBufferSize:    getEnvInt("BLINK_FORMAT_BUFFER_SIZE", 1024),
		PayloadPoolSize:     getEnvInt("BLINK_PAYLOAD_POOL_SIZE", 8),
		TimeFormat:          getEnvString("BLINK_TIME_FORMAT", "2006-01-02T15:04:05.000Z07:00"),
		TextPattern:         getEnvString("BLINK_TEXT_PATTERN", "%{timestamp} %{level} %{referenceId} %{caller} [%{fields}] - %{message}"),
		JsonEscapeHtml:      getEnvBool("BLINK_JSON_ESCAPE_HTML", false),
		JsonCompact:         getEnvBool("BLINK_JSON_COMPACT", true),
		IncludeEmptyFields:  getEnvBool("BLINK_INCLUDE_EMPTY_FIELDS", false),
		EnablePanicRecovery: getEnvBool("BLINK_ENABLE_PANIC_RECOVERY", true),
	}
}

// LoadConfigWithOverrides loads configuration from environment variables and applies
// programmatic overrides using a callback function. This allows for hybrid configuration
// where most settings come from environment variables but specific values are overridden in code.
//
// Parameters:
//
//	overrides func(*ConfigBuilder): Callback function to modify the loaded configuration
//
// Returns:
//
//	Config: Final configuration with environment values and overrides applied
//
// Example:
//
//	config := LoadConfigWithOverrides(func(b *ConfigBuilder) {
//	    b.WriterBufferSize(8192).JsonCompact(false)
//	})
func LoadConfigWithOverrides(overrides func(*ConfigBuilder)) Config {
	envConfig := LoadConfig()
	builder := &ConfigBuilder{config: envConfig}
	overrides(builder)
	return builder.Build()
}

// Clone creates a deep copy of the configuration.
// This is useful when you need to modify a configuration without affecting the original.
// Since Config contains only value types, a simple struct copy provides deep cloning.
//
// Returns:
//
//	Config: A complete copy of the configuration
func (c Config) Clone() Config {
	return Config{
		WriterBufferSize:    c.WriterBufferSize,
		FormatBufferSize:    c.FormatBufferSize,
		PayloadPoolSize:     c.PayloadPoolSize,
		TimeFormat:          c.TimeFormat,
		TextPattern:         c.TextPattern,
		JsonEscapeHtml:      c.JsonEscapeHtml,
		JsonCompact:         c.JsonCompact,
		IncludeEmptyFields:  c.IncludeEmptyFields,
		EnablePanicRecovery: c.EnablePanicRecovery,
	}
}

// validate checks if the configuration values are valid and consistent.
// It performs comprehensive validation of all configuration fields and
// returns detailed error messages for any invalid values.
//
// Validation rules:
//   - All buffer sizes must be positive integers
//   - TimeFormat and TextPattern must not be empty strings
//   - No additional validation is performed on format strings (Go's time package handles validation)
//
// Returns:
//
//	error: nil if configuration is valid, descriptive error otherwise
func (c Config) validate() error {
	if c.WriterBufferSize <= 0 {
		return fmt.Errorf("WriterBufferSize must be positive, got %d", c.WriterBufferSize)
	}
	if c.FormatBufferSize <= 0 {
		return fmt.Errorf("FormatBufferSize must be positive, got %d", c.FormatBufferSize)
	}
	if c.PayloadPoolSize <= 0 {
		return fmt.Errorf("PayloadPoolSize must be positive, got %d", c.PayloadPoolSize)
	}
	if c.TimeFormat == "" {
		return fmt.Errorf("TimeFormat cannot be empty")
	}
	if c.TextPattern == "" {
		return fmt.Errorf("TextPattern cannot be empty")
	}
	return nil
}

// getEnvString retrieves a string value from environment variables with fallback.
// Empty environment variables are treated as unset and fall back to the default value.
//
// Parameters:
//
//	key string: Environment variable name
//	defaultValue string: Value to return if environment variable is unset or empty
//
// Returns:
//
//	string: Environment variable value or default if not found/empty
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an integer value from environment variables with validation and fallback.
// Invalid integers (non-numeric, negative, or zero) fall back to the default value.
// This ensures that configuration always has valid positive integers for buffer sizes.
//
// Parameters:
//
//	key string: Environment variable name
//	defaultValue int: Value to return if environment variable is invalid
//
// Returns:
//
//	int: Parsed environment variable value or default if invalid/unset
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultValue
}

// getEnvBool retrieves a boolean value from environment variables with flexible parsing.
// Supports various common boolean representations (case-insensitive):
//
//	true:  "true", "1", "yes", "on", "enabled"
//	false: "false", "0", "no", "off", "disabled"
//
// Invalid or unrecognized values fall back to the default.
//
// Parameters:
//
//	key string: Environment variable name
//	defaultValue bool: Value to return if environment variable is invalid
//
// Returns:
//
//	bool: Parsed environment variable value or default if invalid/unset
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch strings.ToLower(value) {
		case "true", "1", "yes", "on", "enabled":
			return true
		case "false", "0", "no", "off", "disabled":
			return false
		}
	}
	return defaultValue
}
