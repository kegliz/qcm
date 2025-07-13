// Package qcm provides a modern, high-performance quantum circuit manager and simulator for Go.
//
// QCM offers an intuitive fluent API for building quantum circuits, supports multiple
// quantum simulation backends through a plugin architecture, and includes advanced
// features like circuit visualization and performance optimization.
//
// # Quick Start
//
// Build a Bell state circuit and simulate it:
//
//	import (
//	    "github.com/kegliz/qcm/qc/builder"
//	    "github.com/kegliz/qcm/qc/simulator"
//	    _ "github.com/kegliz/qcm/qc/simulator/itsu"
//	)
//
//	// Create a Bell state circuit
//	circuit := builder.New(builder.Q(2), builder.C(2)).
//	    H(0).           // Apply Hadamard to qubit 0
//	    CNOT(0, 1).     // Apply CNOT with control=0, target=1
//	    Measure(0, 0).  // Measure qubit 0 to classical bit 0
//	    Measure(1, 1)   // Measure qubit 1 to classical bit 1
//
//	// Build and simulate
//	circ, _ := circuit.BuildCircuit()
//	sim, _ := simulator.NewSimulatorWithDefaults("itsu")
//	results, _ := sim.Run(circ)
//
// # Architecture
//
// QCM consists of several key packages:
//
//   - builder: Fluent API for constructing quantum circuits
//   - circuit: Immutable circuit representation with optimized layout
//   - simulator: Multi-backend simulation engine with parallel execution
//   - gate: Comprehensive quantum gate library
//   - renderer: PNG visualization for quantum circuits
//   - dag: Directed Acyclic Graph for circuit dependency management
//
// # Plugin System
//
// QCM supports multiple quantum simulation backends through plugins:
//
//   - itsu: Based on github.com/itsubaki/q
//   - qsim: Custom optimized backend
//
// Import the desired backend plugins to register them:
//
//	import _ "github.com/kegliz/qcm/qc/simulator/itsu"
//	import _ "github.com/kegliz/qcm/qc/simulator/qsim"
//
// # Supported Gates
//
// Single-qubit gates: H, X, Y, Z, S
// Multi-qubit gates: CNOT, CZ, SWAP, Toffoli, Fredkin
// Measurement: Measure quantum states to classical bits
//
// # Performance
//
// QCM is optimized for high performance with features like:
//
//   - Parallel shot execution with configurable worker pools
//   - Efficient memory management for large circuits
//   - Backend selection for optimal performance
//   - Built-in benchmarking and performance analysis tools
package qcm
