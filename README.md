# MockHW - Go Hardware Interface Mocks

A Go library providing mock implementations of common hardware interfaces for testing embedded systems code.

## Features

- **UART Mock**: A simulated UART (Universal Asynchronous Receiver-Transmitter) interface with:
  - Configurable read/write operations
  - Simulated delays to mimic real hardware behavior
  - Buffer management for testing I/O operations
  - Hardware FIFO simulation

## Installation

```bash
go get github.com/m-s-sh/mockhw
```

## Usage

### UART Mock Example

```go
import (
    "github.com/m-s-sh/mockhw"
    "testing"
)

func TestDeviceWithUart(t *testing.T) {
    // Create a mock UART with a maximum delay of 10ms
    uart := mockhw.NewUart(10)
    
    // Inject test data to be read by the device under test
    uart.SetRxBuffer([]byte("test data"))
    
    // Pass the mock UART to your device under test
    device := NewDeviceWithUart(uart)
    
    // Perform operations with your device
    device.ProcessInput()
    
    // Verify data written by the device to the UART
    txData := uart.TxBuffer()
    // Assert on txData...
}
```

## Hardware Interfaces

### Uart

The `Uart` struct implements interfaces typically used for UART communication. It simulates:

- Random delays in reading data (configurable)
- Hardware FIFO buffer behavior
- Chunk-based reading similar to real hardware

#### Methods

- `NewUart(maxDelayMs int) *Uart` - Creates a new UART mock with specified max delay
- `Read(p []byte) (n int, err error)` - Reads data with simulated hardware behavior
- `Write(p []byte) (n int, err error)` - Writes data to the transmission buffer
- `Buffered() int` - Returns number of bytes available for reading
- `SetRxBuffer(data []byte) (n int, err error)` - Sets data to be read
- `TxBuffer() []byte` - Gets data written to the UART

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.