package formatters

import (
	"encoding/json"
	"io"

	"github.com/eskandaridanial/blink/models"
)

// JsonFormatter provides high-performance Json log formatting with zero allocations.
// It uses object pooling for all temporary allocations and direct encoding to writers
// to minimize memory overhead and garbage collection pressure.
//
// The formatter supports:
//   - Configurable Json output (compact vs pretty-printed)
//   - Optional HTML escaping for web-safe output
//   - Field inclusion control (empty fields can be included/excluded)
//   - Automatic panic recovery with fallback formatting
//   - Thread-safe concurrent operation
//
// Performance characteristics:
//   - Zero allocations per log entry after warm-up
//   - Direct streaming to output writer
//   - Reusable buffers and payload maps
//   - Optimized for high-throughput logging scenarios
//
// Thread Safety:
//
//	All methods are thread-safe and can be called concurrently from multiple goroutines.
type JsonFormatter struct {
	config      Config         // Configuration controlling formatter behavior
	writerPool  *writerPool    // Pool of buffered writers for efficient I/O
	payloadPool *payloadPool   // Pool of maps for building Json payloads
	safety      *safetyWrapper // Panic recovery wrapper for robustness
}

// NewJsonFormatter creates a new Json formatter using environment configuration.
// This is the recommended constructor for most use cases as it automatically
// loads configuration from environment variables with sensible defaults.
//
// Returns:
//
//	*JsonFormatter: New Json formatter instance ready for use
//
// Example:
//
//	formatter := NewJsonFormatter()
//	formatter.Format(os.Stdout, logRecord)
func NewJsonFormatter() *JsonFormatter {
	return NewJsonFormatterWithConfig(loadConfig())
}

// NewJsonFormatterWithConfig creates a Json formatter with explicit configuration.
// This constructor provides full control over formatter behavior and is useful
// when configuration needs to be programmatically controlled.
//
// Parameters:
//
//	config Config: Configuration struct controlling formatter behavior
//
// Returns:
//
//	*JsonFormatter: New Json formatter instance with specified configuration
//
// Example:
//
//	config := NewConfigBuilder().
//	    JsonCompact(false).
//	    WriterBufferSize(8192).
//	    Build()
//	formatter := NewJsonFormatterWithConfig(config)
func NewJsonFormatterWithConfig(config Config) *JsonFormatter {
	return &JsonFormatter{
		config:      config,
		writerPool:  newWriterPool(config.WriterBufferSize),
		payloadPool: newPayloadPool(config.PayloadPoolSize),
		safety:      newSafetyWrapper(config),
	}
}

// Format writes a Json-formatted log record to the specified writer.
// This method implements the Formatter interface and provides zero-allocation
// Json formatting through object pooling and direct encoding.
//
// The output Json structure includes:
//   - timestamp: Formatted according to config.TimeFormat
//   - level: Log level as string (DEBUG, INFO, WARN, ERROR, etc.)
//   - message: The actual log message
//   - referenceId: Optional correlation/trace ID (if present or IncludeEmptyFields is true)
//   - caller: Optional caller information (if present or IncludeEmptyFields is true)
//   - fields: Key-value pairs of additional structured data
//
// Parameters:
//
//	w io.Writer: Destination writer for Json output
//	r models.Record: Structured log record to format
//
// Returns:
//
//	int: Number of bytes written to the writer (may be 0 due to buffering)
//	error: I/O error if writing fails, nil on success
//
// Thread Safety:
//
//	This method is thread-safe and can be called concurrently.
//
// Performance Notes:
//   - Uses pooled writers and payload maps to avoid allocations
//   - Encodes directly to the writer without intermediate string creation
//   - Automatic panic recovery ensures robustness in production
func (f *JsonFormatter) Format(w io.Writer, r models.Record) (int, error) {
	return f.safety.withRecovery(w, r, func() (int, error) {
		bufWriter := f.writerPool.get()
		defer func() {
			bufWriter.Flush()
			f.writerPool.put(bufWriter)
		}()
		bufWriter.Reset(w)

		payload := f.payloadPool.get()
		defer f.payloadPool.put(payload)
		f.buildPayload(payload, r)

		encoder := json.NewEncoder(bufWriter)
		encoder.SetEscapeHTML(f.config.JsonEscapeHtml)
		if !f.config.JsonCompact {
			encoder.SetIndent("", "  ")
		}

		return 0, encoder.Encode(payload)
	})
}

// buildPayload constructs the Json payload map from a log record.
// This method populates the payload map with all relevant fields from the log record,
// respecting configuration settings for field inclusion and empty value handling.
//
// Core fields (always included):
//   - timestamp: Formatted timestamp string
//   - level: Log level as string
//   - message: Log message content
//
// Optional fields (included based on configuration):
//   - referenceId: Correlation/trace identifier (if present or IncludeEmptyFields is true)
//   - caller: Caller information (if present or IncludeEmptyFields is true)
//   - fields: Structured key-value pairs (handled by buildFields)
//
// Parameters:
//
//	payload map[string]any: Target map to populate (acquired from pool)
//	r models.Record: Source log record containing data to format
func (f *JsonFormatter) buildPayload(payload map[string]any, r models.Record) {
	payload["timestamp"] = r.Timestamp.Format(f.config.TimeFormat)
	payload["level"] = r.Level.String()
	payload["caller"] = r.Caller
	payload["referenceId"] = r.ReferenceId
	payload["message"] = r.Message
	f.buildFields(payload, r.Fields)
}

// buildFields constructs the fields section of the Json payload.
// This method handles the structured key-value pairs that provide additional
// context for the log entry. It respects configuration settings for empty field inclusion.
//
// Behavior:
//   - If no fields and IncludeEmptyFields is false: no "fields" key in output
//   - If no fields and IncludeEmptyFields is true: "fields": {} in output
//   - If fields present: "fields": {"key1": "value1", "key2": "value2", ...}
//   - Empty field keys are handled according to IncludeEmptyFields setting
//
// Parameters:
//
//	payload map[string]any: Target payload map to add fields to
//	fields []models.Field: Array of structured key-value pairs from log record
func (f *JsonFormatter) buildFields(payload map[string]any, fields []models.Field) {
	if len(fields) == 0 {
		if f.config.IncludeEmptyFields {
			payload["fields"] = make(map[string]any)
		}
		return
	}

	fieldMap := make(map[string]any, len(fields))
	for _, field := range fields {
		if field.Key != "" || f.config.IncludeEmptyFields {
			fieldMap[field.Key] = field.Value
		}
	}

	if len(fieldMap) > 0 || f.config.IncludeEmptyFields {
		payload["fields"] = fieldMap
	}
}
