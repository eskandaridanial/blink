package formatters

import (
	"fmt"
	"strconv"

	"github.com/eskandaridanial/blink/models"
)

// patternParser handles text pattern parsing and placeholder substitution.
// It provides efficient parsing of format patterns with placeholder replacement,
// optimized for minimal allocations and high performance.
//
// The parser supports the following placeholders:
//
//	%{timestamp}   - Formatted timestamp using configured time format
//	%{level}       - Log level as string (DEBUG, INFO, WARN, ERROR, etc.)
//	%{referenceId} - Optional correlation/trace ID
//	%{caller}      - Optional caller information (file:line)
//	%{message}     - The actual log message content
//	%{fields}      - Formatted key-value pairs as "key=value key2=value2"
//
// Pattern parsing algorithm:
//   - Scans pattern character by character
//   - Recognizes %{...} placeholder syntax
//   - Handles malformed patterns gracefully (treats as literal text)
//   - Ensures output always ends with newline
//
// Performance optimizations:
//   - Direct buffer operations without string allocations
//   - Type-optimized value serialization for common types
//   - Minimal pattern parsing overhead
//   - Reuses buffer capacity across invocations
type patternParser struct {
	config Config
}

// newPatternParser creates a new pattern parser with the specified configuration.
// The parser uses the configuration to access the text pattern and field inclusion settings.
//
// Parameters:
//
//	config Config: Configuration containing pattern and field inclusion settings
//
// Returns:
//
//	*patternParser: New pattern parser instance
func newPatternParser(config Config) *patternParser {
	return &patternParser{config: config}
}

// format applies the configured pattern to format a log record into the provided buffer.
// This method performs placeholder substitution and builds the complete formatted output.
// The output is guaranteed to end with a newline character.
//
// Pattern parsing algorithm:
//  1. Scan pattern character by character
//  2. When %{ is found, look for closing }
//  3. Extract placeholder name and substitute with record data
//  4. Handle malformed placeholders as literal text
//  5. Ensure final output ends with newline
//
// Parameters:
//
//	buf []byte: Target buffer to append formatted output to (should have sufficient capacity)
//	r models.Record: Source log record containing data for placeholder substitution
//
// Returns:
//
//	[]byte: Buffer with formatted output appended (may have grown if capacity was insufficient)
//
// Performance Notes:
//   - Direct buffer operations avoid string allocations
//   - Type-optimized serialization for field values
//   - Minimal parsing overhead with single-pass algorithm
func (p *patternParser) format(buf []byte, r models.Record) []byte {
	pattern := p.config.TextPattern
	i := 0

	for i < len(pattern) {
		if i < len(pattern)-1 && pattern[i] == '%' && pattern[i+1] == '{' {
			j := i + 2
			for j < len(pattern) && pattern[j] != '}' {
				j++
			}

			if j < len(pattern) {
				fieldName := pattern[i+2 : j]
				buf = p.appendField(buf, fieldName, r)
				i = j + 1
			} else {
				buf = append(buf, pattern[i])
				i++
			}
		} else {
			buf = append(buf, pattern[i])
			i++
		}
	}

	if len(buf) == 0 || buf[len(buf)-1] != '\n' {
		buf = append(buf, '\n')
	}

	return buf
}

// appendField appends the value of a specified field placeholder to the buffer.
// This method handles all supported placeholder types and respects configuration
// settings for empty field inclusion.
//
// Supported placeholders:
//
//	timestamp   - Formatted timestamp using configured time format
//	level       - Log level as string
//	referenceId - Correlation/trace ID (if present or IncludeEmptyFields is true)
//	caller      - Caller information (if present or IncludeEmptyFields is true)
//	message     - Log message content
//	fields      - All structured fields formatted as key-value pairs
//	<unknown>   - Unknown placeholders are preserved as literal text
//
// Parameters:
//
//	buf []byte: Target buffer to append field value to
//	fieldName string: Name of the placeholder field to append
//	r models.Record: Source log record containing field data
//
// Returns:
//
//	[]byte: Buffer with field value appended (may have grown)
func (p *patternParser) appendField(buf []byte, fieldName string, r models.Record) []byte {
	switch fieldName {
	case "timestamp":
		return append(buf, r.Timestamp.Format(p.config.TimeFormat)...)

	case "level":
		return append(buf, r.Level.String()...)

	case "referenceId":
		return append(buf, r.ReferenceId...)

	case "caller":
		return append(buf, r.Caller...)

	case "message":
		return append(buf, r.Message...)

	case "fields":
		return p.appendFields(buf, r.Fields)

	default:
		return append(buf, "%{"+fieldName+"}"...)
	}
}

// appendFields formats and appends all structured fields to the buffer.
// Fields are formatted as space-separated key-value pairs in the format "key=value".
// This method handles empty field arrays and respects configuration for empty field inclusion.
//
// Output format examples:
//
//	No fields: "" (empty string)
//	One field: "userId=12345"
//	Multiple fields: "userId=12345 requestId=abc-def-ghi status=200"
//
// Parameters:
//
//	buf []byte: Target buffer to append formatted fields to
//	fields []models.Field: Array of structured key-value pairs from log record
//
// Returns:
//
//	[]byte: Buffer with formatted fields appended (may have grown)
func (p *patternParser) appendFields(buf []byte, fields []models.Field) []byte {
	if len(fields) == 0 && !p.config.IncludeEmptyFields {
		return buf
	}

	first := true
	for _, field := range fields {
		if field.Key == "" && !p.config.IncludeEmptyFields {
			continue
		}

		if !first {
			buf = append(buf, ' ')
		}
		first = false

		// Format field as key=value
		buf = append(buf, field.Key...)
		buf = append(buf, '=')
		buf = p.appendFieldValue(buf, field.Value)
	}

	return buf
}

// appendFieldValue efficiently appends a field value to the buffer with type-optimized serialization.
// This method provides fast serialization for common types without allocating intermediate strings.
// It uses Go's strconv package for numeric types and handles special cases like Stringer and error interfaces.
//
// Supported types (with optimized serialization):
//
//	string           - Direct append (no conversion)
//	int/int32/int64  - Optimized integer formatting
//	uint/uint32/uint64 - Optimized unsigned integer formatting
//	float32/float64  - Optimized floating point formatting
//	bool             - Optimized boolean formatting
//	fmt.Stringer     - Uses String() method
//	error            - Uses Error() method
//	other types      - Falls back to fmt.Sprintf("%v", value)
//
// Parameters:
//
//	buf []byte: Target buffer to append serialized value to
//	value any: Value to serialize and append (can be any type)
//
// Returns:
//
//	[]byte: Buffer with serialized value appended (may have grown)
//
// Performance Notes:
//   - Type switch provides O(1) dispatch for common types
//   - Uses strconv.Append* functions to avoid string allocations
//   - Only falls back to fmt.Sprintf for uncommon types
//   - Direct buffer operations minimize memory copying
func (p *patternParser) appendFieldValue(buf []byte, value any) []byte {
	switch v := value.(type) {
	case string:
		return append(buf, v...)

	case int:
		return strconv.AppendInt(buf, int64(v), 10)

	case int32:
		return strconv.AppendInt(buf, int64(v), 10)

	case int64:
		return strconv.AppendInt(buf, v, 10)

	case uint:
		return strconv.AppendUint(buf, uint64(v), 10)

	case uint32:
		return strconv.AppendUint(buf, uint64(v), 10)

	case uint64:
		return strconv.AppendUint(buf, v, 10)

	case float32:
		return strconv.AppendFloat(buf, float64(v), 'f', -1, 32)

	case float64:
		return strconv.AppendFloat(buf, v, 'f', -1, 64)

	case bool:
		return strconv.AppendBool(buf, v)

	case fmt.Stringer:
		return append(buf, v.String()...)

	case error:
		return append(buf, v.Error()...)

	default:
		return append(buf, fmt.Sprintf("%v", v)...)
	}
}
