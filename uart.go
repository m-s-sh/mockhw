// Package mockhw provides a mock implementation of UART interfaces for testing.
package mockhw

import (
	"bytes"
	"math/rand"
	"time"
)

// UART implements the UART interface for testing purposes.
// It simulates real UART behavior with random delays during reading.
type UART struct {
	txBuffer     *bytes.Buffer // Buffer for transmitted data
	rxBuffer     *bytes.Buffer // Buffer for received data
	maxDelay     time.Duration // Maximum delay for read operations
	availableLen int           // Number of bytes reported as available/buffered
}

// NewUART creates a new UART instance.
// maxDelayMs is the maximum delay in milliseconds for read operations.
func NewUART(maxDelayMs int) *UART {
	return &UART{
		txBuffer:     bytes.NewBuffer(nil),
		rxBuffer:     bytes.NewBuffer(nil),
		maxDelay:     time.Duration(maxDelayMs) * time.Millisecond,
		availableLen: 0,
	}
}

// updateAvailableBytes recalculates the number of available bytes to read
// to simulate UART hardware FIFO buffer behavior
func (m *UART) updateAvailableBytes() {
	totalBytes := m.rxBuffer.Len()
	if totalBytes == 0 {
		m.availableLen = 0
		return
	}

	// Real UARTs typically have small hardware FIFOs (16, 32, 64 bytes)
	// Simulate this by making available only a random portion of the actual data
	const maxFifoSize = 16

	if totalBytes <= maxFifoSize {
		m.availableLen = totalBytes
	} else {
		// Return a random number between 1 and maxFifoSize or the total bytes available
		m.availableLen = rand.Intn(maxFifoSize) + 1
		if m.availableLen > totalBytes {
			m.availableLen = totalBytes
		}
	}
}

// Read implements the io.Reader interface.
// It reads up to len(p) bytes into p with a random delay and in random chunk sizes
// to simulate real UART behavior.
func (m *UART) Read(p []byte) (n int, err error) {
	// Simulate real UART delay
	if m.maxDelay > 0 {
		delay := time.Duration(rand.Int63n(int64(m.maxDelay)))
		time.Sleep(delay)
	}

	// If there's no data or empty slice, return immediately
	if m.rxBuffer.Len() == 0 || len(p) == 0 {
		return 0, nil
	}

	// Determine a random chunk size - between 1 and min(len(p), available bytes)
	maxSize := len(p)
	available := m.rxBuffer.Len()
	if available < maxSize {
		maxSize = available
	}

	// Read a random number of bytes between 1 and maxSize
	chunkSize := rand.Intn(maxSize) + 1

	// Read the random chunk
	n, err = m.rxBuffer.Read(p[:chunkSize])

	// Recalculate available bytes after each read
	m.updateAvailableBytes()

	return n, err
}

// Write implements the io.Writer interface.
// It writes len(p) bytes from p to the UART's tx buffer.
func (m *UART) Write(p []byte) (n int, err error) {
	return m.txBuffer.Write(p)
}

// Buffered returns the number of bytes that can be read from the rx buffer.
// This simulates how real UARTs report only what's in their hardware FIFO.
func (m *UART) Buffered() int {
	return m.availableLen
}

// SetRxBuffer writes data to the rx buffer to simulate data reception.
// This is useful for testing when you need to simulate incoming data.
func (m *UART) SetRxBuffer(data []byte) (n int, err error) {
	n, err = m.rxBuffer.Write(data)
	m.updateAvailableBytes()
	return n, err
}

// TxBuffer returns the contents of the transmission buffer.
// This is useful for testing to verify what data was sent.
func (m *UART) TxBuffer() []byte {
	return m.txBuffer.Bytes()
}
