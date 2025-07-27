package formatters

import (
	"io"

	"github.com/eskandaridanial/blink/models"
)

// Formatter defines the interface that all log formatters must implement.
// This interface provides a consistent API for different formatting strategies
// while allowing for different internal implementations (JSON, text, binary, etc.).
//
// Implementations must be thread-safe and should optimize for zero allocations
// when possible. The interface is designed to support streaming output and
// allows formatters to report the number of bytes written.
//
// Contract:
//   - Format must be thread-safe and callable from multiple goroutines
//   - Format should minimize memory allocations for high-performance logging
//   - Format should handle nil or invalid records gracefully
//   - Returned byte count should reflect actual bytes written to the writer
//   - Errors should be returned for I/O failures, not formatting issues
type Formatter interface {
	// Format writes the formatted log record to the provided writer.
	// The method transforms a structured log record into its final output format
	// and writes it directly to the writer to avoid intermediate string allocations.
	//
	// Parameters:
	//   w io.Writer: Destination writer for formatted output
	//   record models.Record: Structured log record to format
	//
	// Returns:
	//   int: Number of bytes successfully written to the writer
	//   error: I/O error if writing fails, nil on success
	//
	// Thread Safety:
	//   This method must be safe to call concurrently from multiple goroutines.
	//
	// Performance Notes:
	//   - Implementations should use object pooling to minimize allocations
	//   - Buffered writers should be used for small, frequent writes
	//   - String concatenations should be avoided in favor of direct buffer operations
	Format(w io.Writer, record models.Record) (int, error)
}
