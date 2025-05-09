# Go Worker Pool

A parallel processing system implemented in Go that demonstrates the efficient use of goroutines and channels to create worker pools.

## Features

- Configurable worker pool for concurrent processing
- Control of the maximum number of simultaneous goroutines
- Elegant error handling during processing
- Job queue for task distribution
- Result collection through channels
- Practical implementation with file processor example

## Project Structure

```
worker-pool/
├── cmd/
│   └── main.go             # Application entry point
├── internal/
│   ├── pool/
│   │   └── pool.go         # Worker pool implementation
│   └── processor/
│       └── processor.go    # File processor implementation
├── pkg/
│   └── result/
│       └── result.go       # Result handling definitions
└── README.md               # This file
```

## Installation

```bash
# Clone the repository
git clone https://github.com/MrXetwy21/worker-pool.git

# Enter the directory
cd worker-pool

# Build
go build -o workerpool ./cmd/main.go
```

## Usage

```bash
# Run with default parameters
./workerpool

# Specify directory to process and number of workers
./workerpool -dir="/path/to/directory" -workers=10
```

## Code Example

```go
package main

import (
    "github.com/MrXetwy21/worker-pool/internal/pool"
    "github.com/MrXetwy21/worker-pool/internal/processor"
)

func main() {
    // Create a new pool with 5 workers
    p := pool.New(5)
    
    // Start processing
    p.Process(processor.ProcessFiles("/path/to/files"))
    
    // Get results
    results := p.Wait()
    
    // Do something with the results
    for _, result := range results {
        // Handle each result
    }
}
```
