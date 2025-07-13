# QCM - Quantum Circuit Manager

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/kegliz/qcm)](https://goreportcard.com/report/github.com/kegliz/qcm)

**QCM** is a modern, high-performance quantum circuit manager and simulator written in Go. It provides an intuitive fluent API for building quantum circuits, supports multiple quantum simulation backends through a plugin architecture, and includes advanced features like circuit visualization and performance optimization.

## Features

- 🔧 **Fluent Circuit Builder**: Intuitive API for constructing quantum circuits
- 🚀 **Multiple Simulation Backends**: Plugin architecture supporting different quantum simulators
- 📊 **Circuit Visualization**: Built-in PNG rendering for quantum circuits
- ⚡ **High Performance**: Optimized parallel execution with configurable worker pools
- 🔍 **Comprehensive Gate Set**: Support for common single and multi-qubit gates
- 📈 **Performance Analysis**: Built-in benchmarking and metrics collection
- 🧪 **Testing Framework**: Comprehensive test utilities for quantum circuit validation

## Quick Start

### Installation

```bash
go get github.com/kegliz/qcm@v0.2.0
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/kegliz/qcm/qc/builder"
    "github.com/kegliz/qcm/qc/simulator"
    
    // Import quantum backend plugins
    _ "github.com/kegliz/qcm/qc/simulator/itsu"
)

func main() {
    // Create a Bell state circuit
    circuit := builder.New(builder.Q(2), builder.C(2)).
        H(0).           // Apply Hadamard to qubit 0
        CNOT(0, 1).     // Apply CNOT with control=0, target=1
        Measure(0, 0).  // Measure qubit 0 to classical bit 0
        Measure(1, 1)   // Measure qubit 1 to classical bit 1

    // Build the circuit
    circ, err := circuit.BuildCircuit()
    if err != nil {
        log.Fatal(err)
    }

    // Create and run simulator
    sim, err := simulator.NewSimulatorWithDefaults("itsu")
    if err != nil {
        log.Fatal(err)
    }
    
    // Override default settings if needed
    sim.Shots = 1024

    results, err := sim.Run(circ)
    if err != nil {
        log.Fatal(err)
    }

    // Display results
    for state, count := range results {
        fmt.Printf("|%s⟩: %d shots (%.2f%%)\n", 
            state, count, float64(count)/1024*100)
    }
}
```

## Supported Gates

### Single-Qubit Gates
- **H** - Hadamard gate
- **X** - Pauli-X (NOT) gate  
- **Y** - Pauli-Y gate
- **Z** - Pauli-Z gate
- **S** - S gate (√Z)

### Multi-Qubit Gates
- **CNOT** - Controlled-NOT gate
- **CZ** - Controlled-Z gate
- **SWAP** - Swap gate
- **Toffoli** - Three-qubit controlled-controlled-NOT
- **Fredkin** - Controlled-SWAP gate

### Measurement
- **Measure** - Quantum measurement to classical bits

## Architecture

### Plugin System

QCM features a flexible plugin architecture that allows different quantum simulation backends to be registered and used interchangeably:

```go
// Available backends
_ "github.com/kegliz/qcm/qc/simulator/itsu"  // itsubaki/q backend
_ "github.com/kegliz/qcm/qc/simulator/qsim"  // Custom optimized backend
```

### Circuit Visualization

Generate PNG visualizations of your quantum circuits:

```go
import "github.com/kegliz/qcm/qc/renderer"

// Render circuit to PNG
err := renderer.SaveToPNG(circuit, "my_circuit.png", 
    renderer.WithTitle("Bell State Circuit"))
```

### Performance Optimization

QCM includes several performance optimization features:

- **Parallel Execution**: Configurable worker pools for shot-based simulations
- **Memory Optimization**: Efficient state vector management
- **Backend Selection**: Choose optimal backend for your specific use case

## Examples

### Grover's Algorithm (2-qubit)

```go
// Search for |11⟩ state using 2-qubit Grover's algorithm
circuit := builder.New(builder.Q(2), builder.C(2)).
    H(0).H(1).              // Initialize superposition
    Z(0).Z(1).CNOT(0,1).Z(1).CNOT(0,1). // Oracle for |11⟩
    H(0).H(1).Z(0).Z(1).CNOT(0,1).Z(1).CNOT(0,1).H(0).H(1). // Diffusion
    Measure(0, 0).Measure(1, 1)

circ, _ := circuit.BuildCircuit()
sim, _ := simulator.NewSimulatorWithDefaults("itsu")
sim.Shots = 1024
results, _ := sim.Run(circ)
```

### Z-Gate Demonstration

```go
// Demonstrate Z-gate effect with Hadamard sandwich
circuit := builder.New(builder.Q(1), builder.C(1)).
    H(0).           // Create |+⟩ state
    Z(0).           // Apply Z-gate  
    H(0).           // Convert back to computational basis
    Measure(0, 0)   // Should show bit flip compared to H-H

circ, _ := circuit.BuildCircuit()
```

## Performance Comparison

QCM includes built-in performance benchmarking tools to compare different backends:

```bash
go run cmd/perf-comp/performance-comparison.go
```

Example benchmark results:
```
Backend Performance Comparison (1000 shots each)
================================================
Circuit: 4-qubit Grover
- itsu backend:    245.2ms ± 12.3ms
- qsim backend:    189.7ms ± 8.9ms  (22.6% faster)

Circuit: Bell State  
- itsu backend:     45.1ms ± 2.1ms
- qsim backend:     38.4ms ± 1.8ms  (14.9% faster)
```

## CLI Tools

QCM provides several command-line tools:

```bash
# Run example simulations
go run cmd/cli/main.go

# Performance comparison between backends  
go run cmd/perf-comp/performance-comparison.go

# Plugin demonstration
go run cmd/plugin-demo/main.go
```

## Testing

Run the complete test suite:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Documentation

- [Plugin Architecture](docs/plugin-architecture.md) - Detailed plugin system documentation
- [Performance Optimization](docs/qsim-performance-optimization.md) - Performance tuning guide  
- [Backend Comparison](docs/qsim-vs-itsubaki-performance-comparison.md) - Benchmark analysis

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- 🔧 **v0.3.0**: 

## Citation

If you use QCM in your research, please cite:

```bibtex
@software{qcm2025,
  author = {kegliz},
  title = {QCM: Quantum Circuit Manager},
  url = {https://github.com/kegliz/qcm},
  version = {0.2.0},
  year = {2025}
}
```

---

**Made with ❤️ for the quantum computing community**

