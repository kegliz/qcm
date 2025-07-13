# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-07-13

### Added
- Comprehensive fluent API for quantum circuit building
- Plugin architecture supporting multiple quantum simulation backends
- Built-in support for itsu (itsubaki/q) and qsim backends
- Circuit visualization with PNG rendering capabilities
- High-performance parallel execution with configurable worker pools
- Extensive gate set including:
  - Single-qubit gates: H, X, Y, Z, S
  - Multi-qubit gates: CNOT, CZ, SWAP, Toffoli, Fredkin
  - Measurement operations
- Performance benchmarking tools and backend comparison utilities
- Comprehensive test suite with quantum circuit validation
- Command-line tools for demonstrations and performance analysis
- Complete documentation including plugin architecture guide

### Features
- **Builder Package**: Fluent DSL for constructing quantum circuits
- **Circuit Package**: Immutable circuit representation with layout optimization
- **DAG Package**: Directed Acyclic Graph for circuit dependency management
- **Gate Package**: Comprehensive quantum gate library
- **Simulator Package**: Multi-backend simulation engine with parallel execution
- **Renderer Package**: PNG visualization for quantum circuits
- **Logger Package**: Structured logging with zerolog

### Technical
- Go 1.23+ compatibility
- Zero external dependencies for core functionality
- Modular plugin system for easy backend extension
- Efficient memory management for large circuit simulations
- Thread-safe parallel execution

### Examples
- Bell state preparation and measurement
- Grover's algorithm implementations (2, 3, 4 qubits)
- Z-gate demonstration with Hadamard sandwich
- Performance comparison between backends
- Plugin usage demonstrations

## [Unreleased]

### Planned Features
