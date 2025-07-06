package uart

import (
	"bytes"
	"testing"
	"time"
)

func TestMockUART_Write(t *testing.T) {
	uart := NewMockUART(0)
	testData := []byte("UART test data")

	n, err := uart.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write to UART: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, got %d", len(testData), n)
	}

	txData := uart.GetTxBuffer()
	if !bytes.Equal(txData, testData) {
		t.Errorf("Expected tx buffer to contain %q, got %q", testData, txData)
	}
}

func TestMockUART_ReadInChunks(t *testing.T) {
	uart := NewMockUART(0)
	testData := []byte("UART received data")

	n, err := uart.SetRxBuffer(testData)
	if err != nil {
		t.Fatalf("Failed to set RX buffer: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to receive %d bytes, got %d", len(testData), n)
	}

	readBuf := make([]byte, 0, len(testData))
	buf := make([]byte, 8)
	for {
		// Read in chunks of up to 8 bytes
		n, err = uart.Read(buf)
		if err != nil {
			t.Fatalf("Failed to read from UART: %v", err)
		}
		if n > 0 {
			readBuf = append(readBuf, buf[:n]...)
		} else {
			break // No more data to read
		}
	}
	if !bytes.Equal(readBuf, testData) {
		t.Errorf("Expected read data to be %q, got %q", testData, readBuf)
	}
	if uart.Buffered() != 0 {
		t.Errorf("Expected buffer to be empty after read, got %d bytes", uart.Buffered())
	}
}

func TestMockUART_Delay(t *testing.T) {
	// Use a significant delay that we can measure
	delayMs := 50
	uart := NewMockUART(delayMs)

	// Add some data to read
	uart.SetRxBuffer([]byte("test"))

	// Measure time taken to read
	start := time.Now()
	buf := make([]byte, 4)
	_, err := uart.Read(buf)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	// We should have experienced some delay, but not more than the maximum
	// Allow some buffer for timing inaccuracies
	if elapsed < time.Millisecond || elapsed > time.Duration(delayMs+10)*time.Millisecond {
		t.Errorf("Expected delay between 1-%d ms, got %v", delayMs, elapsed)
	}
}
