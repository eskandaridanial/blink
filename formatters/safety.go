package formatters

import (
	"io"

	"github.com/eskandaridanial/blink/models"
)

// SafetyWrapper provides panic recovery functionality for formatters.
// It implements a safety layer that can catch and handle panics during
// the formatting process, ensuring that logging never brings down the application.
//
// When panic recovery is enabled, any panic during formatting is caught and
// a fallback format is used. This ensures logging always produces some output,
// even if the primary formatting logic fails.
type SafetyWrapper struct {
	Config Config
}

// NewSafetyWrapper creates a new safety wrapper with the provided configuration.
// The wrapper uses the configuration to determine whether panic recovery is enabled
// and what time format to use for fallback formatting.
//
// Parameters:
//
//	config Config: Configuration controlling panic recovery behavior
//
// Returns:
//
//	*safetyWrapper: New safety wrapper instance
func NewSafetyWrapper(config Config) *SafetyWrapper {
	return &SafetyWrapper{Config: config}
}

// WithRecovery executes a formatting function with optional panic recovery.
// If panic recovery is disabled in the configuration, the function executes normally.
// If panic recovery is enabled, any panics are caught and a fallback format is written.
//
// The fallback format ensures that some meaningful output is always produced,
// preventing the loss of log information even when formatting fails.
//
// Parameters:
//
//	w io.Writer: Destination writer for output
//	r models.Record: Log record being formatted (used for fallback)
//	fn func() (int, error): Formatting function to execute safely
//
// Returns:
//
//	n int: Number of bytes written (from function or fallback)
//	err error: Error from function or fallback write operation
//
// Panic Handling:
//   - If panic recovery is disabled: panics propagate to caller
//   - If panic recovery is enabled: panics are caught and fallback format is used
//   - Fallback format: "{timestamp} {level} - FORMATTING_ERROR: {message}\n"
func (s *SafetyWrapper) WithRecovery(w io.Writer, r models.Record, fn func() (int, error)) (n int, err error) {
	if !s.Config.EnablePanicRecovery {
		return fn()
	}

	defer func() {
		if rec := recover(); rec != nil {
			fallback := r.Timestamp.Format(s.Config.TimeFormat) +
				" " + r.Level.String() +
				" - FORMATTING_ERROR: " + r.Message + "\n"
			n, err = w.Write([]byte(fallback))
		}
	}()

	return fn()
}
