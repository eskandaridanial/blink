package formatters

import (
	"io"

	"github.com/eskandaridanial/blink/models"
)

// TextFormatter provides high-performance text log formatting with configurable patterns.
// It supports custom formatting patterns with placeholder substitution and uses object pooling
// for zero-allocation operation in high-throughput scenarios.
//
// The formatter supports:
//   - Configurable text patterns with placeholder substitution
//   - Multiple placeholder types (timestamp, level, message, fields, etc.)
//   - Zero-allocation formatting through buffer pooling
//   - Automatic panic recovery with fallback formatting
//   - Thread-safe concurrent operation
//   - Efficient field formatting with type-optimized serialization
//
// Pattern placeholders:
//
//	%{timestamp}   - Formatted timestamp
//	%{level}       - Log level string
//	%{referenceId} - Correlation/trace ID
//	%{caller}      - Caller information
//	%{message}     - Log message content
//	%{fields}      - Formatted key-value pairs
//
// Performance characteristics:
//   - Zero allocations per log entry after warm-up
//   - Direct buffer operations without string concatenation
//   - Type-optimized field value serialization
//   - Reusable buffers and writers
//
// Thread Safety:
//
//	All methods are thread-safe and can be called concurrently from multiple goroutines.
type TextFormatter struct {
	config     Config         // Configuration controlling formatter behavior
	writerPool *writerPool    // Pool of buffered writers for efficient I/O
	bufferPool *bufferPool    // Pool of byte buffers for building output
	safety     *safetyWrapper // Panic recovery wrapper for robustness
	parser     *patternParser // Pattern parser for placeholder substitution
}

// newTextFormatter creates a new text formatter using environment configuration.
// This constructor loads configuration from environment variables and uses the
// default text pattern for formatting output.
//
// Returns:
//
//	*TextFormatter: New text formatter instance ready for use
//
// Example:
//
//	formatter := NewTextFormatter()
//	formatter.Format(os.Stdout, logRecord)
func NewTextFormatter() *TextFormatter {
	return NewTextFormatterWithConfig(loadConfig())
}

// NewTextFormatterWithPattern creates a text formatter with a custom pattern.
// This constructor allows specifying a custom formatting pattern while using
// environment configuration for other settings.
//
// Parameters:
//
//	pattern string: Custom formatting pattern with placeholders (if empty, uses config default)
//
// Returns:
//
//	*TextFormatter: New text formatter instance with specified pattern
//
// Example:
//
//	formatter := NewTextFormatterWithPattern("[%{level}] %{timestamp} - %{message}")
//	formatter.Format(os.Stdout, logRecord)
func NewTextFormatterWithPattern(pattern string) *TextFormatter {
	config := loadConfig()
	if pattern != "" {
		config.TextPattern = pattern
	}
	return NewTextFormatterWithConfig(config)
}

// NewTextFormatterWithConfig creates a text formatter with explicit configuration.
// This constructor provides full control over all formatter settings and is useful
// when configuration needs to be programmatically controlled.
//
// Parameters:
//
//	config Config: Configuration struct controlling formatter behavior
//
// Returns:
//
//	*TextFormatter: New text formatter instance with specified configuration
//
// Example:
//
//	config := NewConfigBuilder().
//	    TextPattern("%{timestamp} [%{level}] %{message}").
//	    WriterBufferSize(8192).
//	    Build()
//	formatter := NewTextFormatterWithConfig(config)
func NewTextFormatterWithConfig(config Config) *TextFormatter {
	return &TextFormatter{
		config:     config,
		writerPool: newWriterPool(config.WriterBufferSize),
		bufferPool: newBufferPool(config.FormatBufferSize),
		safety:     newSafetyWrapper(config),
		parser:     newPatternParser(config),
	}
}

// Format writes a text-formatted log record to the specified writer.
// This method implements the Formatter interface and provides zero-allocation
// text formatting through object pooling and direct buffer operations.
//
// The output format is determined by the configured text pattern, with placeholders
// replaced by actual log record values. The formatted output always ends with a newline.
//
// Parameters:
//
//	w io.Writer: Destination writer for text output
//	r models.Record: Structured log record to format
//
// Returns:
//
//	int: Number of bytes written to the writer
//	error: I/O error if writing fails, nil on success
//
// Thread Safety:
//
//	This method is thread-safe and can be called concurrently.
//
// Performance Notes:
//   - Uses pooled writers and buffers to avoid allocations
//   - Direct buffer operations without string concatenation
//   - Type-optimized field value serialization
//   - Automatic panic recovery ensures robustness
func (f *TextFormatter) Format(w io.Writer, r models.Record) (int, error) {
	return f.safety.withRecovery(w, r, func() (int, error) {
		bufWriter := f.writerPool.get()
		defer func() {
			bufWriter.Flush()
			f.writerPool.put(bufWriter)
		}()
		bufWriter.Reset(w)

		bufPtr := f.bufferPool.get()
		buf := (*bufPtr)[:0]
		defer f.bufferPool.put(bufPtr)

		buf = f.parser.format(buf, r)

		return bufWriter.Write(buf)
	})
}
