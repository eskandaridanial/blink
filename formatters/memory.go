package formatters

import (
	"bufio"
	"sync"
)

// writerPool manages a pool of reusable buffered writers to minimize allocations.
// Buffered writers reduce the number of system calls by batching small writes,
// significantly improving performance for frequent logging operations.
//
// The pool ensures that buffered writers are properly reset between uses
// and maintains a consistent buffer size across all instances.
//
// Thread Safety:
//
//	All methods are thread-safe using sync.Pool's internal synchronization.
type writerPool struct {
	pool sync.Pool
}

// newWriterPool creates a new writer pool with the specified buffer size.
// All writers created by this pool will use the same buffer size,
// ensuring consistent memory usage and performance characteristics.
//
// Parameters:
//
//	bufferSize int: Buffer size in bytes for each buffered writer
//
// Returns:
//
//	*writerPool: New writer pool instance
func newWriterPool(bufferSize int) *writerPool {
	return &writerPool{
		pool: sync.Pool{
			// New function creates fresh buffered writers when the pool is empty
			// The writer is created with a nil underlying writer, which must be
			// set via Reset() before use
			New: func() any {
				return bufio.NewWriterSize(nil, bufferSize)
			},
		},
	}
}

// get retrieves a buffered writer from the pool.
// If the pool is empty, a new writer is created using the pool's New function.
// The returned writer must be reset with a destination before use.
//
// Returns:
//
//	*bufio.Writer: Buffered writer ready for reset and use
//
// Usage:
//
//	writer := pool.get()
//	defer pool.put(writer)
//	writer.Reset(destination)
func (p *writerPool) get() *bufio.Writer {
	return p.pool.Get().(*bufio.Writer)
}

// put returns a buffered writer to the pool for reuse.
// The writer should be flushed before returning to ensure all data is written.
// The writer's underlying destination is not reset, but will be overwritten
// on the next use via Reset().
//
// Parameters:
//
//	writer *bufio.Writer: Buffered writer to return to pool
//
// Best Practice:
//
//	Always flush the writer before putting it back:
//	writer.Flush()
//	pool.put(writer)
func (p *writerPool) put(writer *bufio.Writer) {
	p.pool.Put(writer)
}

// BufferPool manages a pool of reusable byte slices to minimize garbage collection.
// Byte slices are used for building formatted output before writing to the final destination.
// Pooling these buffers eliminates the allocation overhead of creating new slices for each log entry.
//
// The pool maintains buffers with a consistent initial capacity but allows them to grow
// as needed. When returned to the pool, buffers retain their capacity but have their
// length reset to zero.
//
// Thread Safety:
//
//	All methods are thread-safe using sync.Pool's internal synchronization.
type bufferPool struct {
	pool sync.Pool
}

// newBufferPool creates a new buffer pool with the specified initial capacity.
// All buffers created by this pool start with the same capacity,
// but can grow larger if needed during formatting operations.
//
// Parameters:
//
//	capacity int: Initial capacity in bytes for each buffer
//
// Returns:
//
//	*bufferPool: New buffer pool instance
func newBufferPool(capacity int) *bufferPool {
	return &bufferPool{
		pool: sync.Pool{
			// New function creates fresh byte slices when the pool is empty
			// Slices are created with zero length but specified capacity
			New: func() any {
				buf := make([]byte, 0, capacity)
				return &buf
			},
		},
	}
}

// get retrieves a byte slice buffer from the pool.
// The returned buffer has zero length but may have significant capacity.
// If the pool is empty, a new buffer is created with the configured initial capacity.
//
// Returns:
//
//	*[]byte: Pointer to byte slice buffer ready for use
//
// Usage:
//
//	bufPtr := pool.Get()
//	defer pool.Put(bufPtr)
//	buf := (*bufPtr)[:0] // Reset length while preserving capacity
func (p *bufferPool) get() *[]byte {
	return p.pool.Get().(*[]byte)
}

// put returns a byte slice buffer to the pool for reuse.
// The buffer's length is not reset here - the caller should reset it when retrieved.
// The buffer retains its capacity, allowing for efficient reuse.
//
// Parameters:
//
//	buf *[]byte: Pointer to byte slice buffer to return to pool
//
// Note:
//
//	The buffer contents are not cleared, but the length should be reset to 0
//	when the buffer is next retrieved from the pool.
func (p *bufferPool) put(buf *[]byte) {
	p.pool.Put(buf)
}

// payloadPool manages a pool of reusable map[string]any instances for JSON formatting.
// JSON formatting requires building a map structure before encoding, and these maps
// can be expensive to allocate for each log entry. Pooling eliminates this overhead.
//
// Maps are cleared when returned to the pool to prevent memory leaks from
// retained references, but they maintain their underlying capacity for efficient reuse.
//
// Thread Safety:
//
//	All methods are thread-safe using sync.Pool's internal synchronization.
type payloadPool struct {
	pool sync.Pool
}

// newPayloadPool creates a new payload pool with the specified initial capacity.
// All maps created by this pool start with the same capacity,
// reducing the need for map growth during typical formatting operations.
//
// Parameters:
//
//	capacity int: Initial capacity for each map (number of key-value pairs)
//
// Returns:
//
//	*payloadPool: New payload pool instance
func newPayloadPool(capacity int) *payloadPool {
	return &payloadPool{
		pool: sync.Pool{
			// New function creates fresh maps when the pool is empty
			// Maps are created with the specified initial capacity to minimize growth
			New: func() any {
				return make(map[string]any, capacity)
			},
		},
	}
}

// Get retrieves a payload map from the pool.
// The returned map is empty but may have significant underlying capacity.
// If the pool is empty, a new map is created with the configured initial capacity.
//
// Returns:
//
//	map[string]any: Empty map ready for population with log data
//
// Usage:
//
//	payload := pool.get()
//	defer pool.put(payload)
//	payload["timestamp"] = timestamp
//	payload["level"] = level
func (p *payloadPool) get() map[string]any {
	return p.pool.Get().(map[string]any)
}

// put returns a payload map to the pool after clearing all entries.
// The map is cleared to prevent memory leaks from retained references,
// but maintains its underlying capacity for efficient reuse.
//
// Parameters:
//
//	payload map[string]any: Map to clear and return to pool
//
// Implementation:
//
//	All key-value pairs are deleted from the map before returning it to the pool.
//	This prevents references to log data from being retained across uses.
func (p *payloadPool) put(payload map[string]any) {
	// Clear all entries from the map to prevent memory leaks
	// The map retains its capacity for efficient reuse
	for k := range payload {
		delete(payload, k)
	}
	p.pool.Put(payload)
}
