package main

import (
	"fmt"

	"github.com/kegliz/qcm/qc/builder"
	"github.com/kegliz/qcm/qc/simulator"
	_ "github.com/kegliz/qcm/qc/simulator/qsim"
)

func main() {
	shots := 1024

	fmt.Println("\n--- Deutsch-Jozsa Algorithm Demonstrations ---")
	deutschJozsaDemo(shots)
}

// deutschJozsaDemo demonstrates the Deutsch-Jozsa algorithm with different oracle functions
func deutschJozsaDemo(shots int) {
	fmt.Println("\n=== 2-Qubit Deutsch-Jozsa Algorithm ===")

	// Test constant function f(x) = 0
	fmt.Println("\n1. Testing constant function f(x) = 0:")
	deutschJozsa2Qubit(shots, "constant_0")

	// Test constant function f(x) = 1
	fmt.Println("\n2. Testing constant function f(x) = 1:")
	deutschJozsa2Qubit(shots, "constant_1")

	// Test balanced function f(0) = 0, f(1) = 1 (identity)
	fmt.Println("\n3. Testing balanced function f(x) = x:")
	deutschJozsa2Qubit(shots, "balanced_identity")

	// Test balanced function f(0) = 1, f(1) = 0 (NOT)
	fmt.Println("\n4. Testing balanced function f(x) = NOT x:")
	deutschJozsa2Qubit(shots, "balanced_not")

	fmt.Println("\n=== 3-Qubit Deutsch-Jozsa Algorithm ===")

	// Test constant function f(x) = 0
	fmt.Println("\n5. Testing constant function f(x) = 0 (3-qubit):")
	deutschJozsa3Qubit(shots, "constant_0")

	// Test balanced function - XOR of first two bits
	fmt.Println("\n6. Testing balanced function f(x1,x2) = x1 ⊕ x2 (3-qubit):")
	deutschJozsa3Qubit(shots, "balanced_xor")
}

// deutschJozsa2Qubit implements the 2-qubit Deutsch-Jozsa algorithm
// Uses 2 qubits: 1 input qubit + 1 ancilla qubit for the oracle
func deutschJozsa2Qubit(shots int, oracleType string) {
	// 2 qubits: qubit 0 (input), qubit 1 (ancilla)
	// 1 classical bit to measure the result
	b := builder.New(builder.Q(2), builder.C(1))

	// Initialize ancilla qubit in |1⟩ state
	b.X(1)

	// Apply Hadamard to both qubits
	b.H(0).H(1)

	// Apply oracle function based on type
	applyOracle2Qubit(b, oracleType)

	// Apply Hadamard to input qubit
	b.H(0)

	// Measure input qubit
	b.Measure(0, 0)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building Deutsch-Jozsa circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running Deutsch-Jozsa simulation: %v\n", err)
		return
	}

	// Analyze results
	analyzeResults(hist, shots, oracleType)
}

// deutschJozsa3Qubit implements the 3-qubit Deutsch-Jozsa algorithm
// Uses 3 qubits: 2 input qubits + 1 ancilla qubit
func deutschJozsa3Qubit(shots int, oracleType string) {
	// 3 qubits: qubit 0,1 (input), qubit 2 (ancilla)
	// 2 classical bits to measure the results
	b := builder.New(builder.Q(3), builder.C(2))

	// Initialize ancilla qubit in |1⟩ state
	b.X(2)

	// Apply Hadamard to all qubits
	b.H(0).H(1).H(2)

	// Apply oracle function based on type
	applyOracle3Qubit(b, oracleType)

	// Apply Hadamard to input qubits
	b.H(0).H(1)

	// Measure input qubits
	b.Measure(0, 0).Measure(1, 1)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building 3-qubit Deutsch-Jozsa circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running 3-qubit Deutsch-Jozsa simulation: %v\n", err)
		return
	}

	// Analyze results
	analyzeResults3Qubit(hist, shots, oracleType)
}

// applyOracle2Qubit applies the oracle function for 2-qubit Deutsch-Jozsa
func applyOracle2Qubit(b builder.Builder, oracleType string) {
	switch oracleType {
	case "constant_0":
		// f(x) = 0: Do nothing (oracle outputs 0 for all inputs)
		// No gates needed
	case "constant_1":
		// f(x) = 1: Flip the ancilla qubit for all inputs
		b.X(1)
	case "balanced_identity":
		// f(x) = x: Apply CNOT(input, ancilla)
		b.CNOT(0, 1)
	case "balanced_not":
		// f(x) = NOT x: Apply X to input, then CNOT, then X again
		b.X(0).CNOT(0, 1).X(0)
	}
}

// applyOracle3Qubit applies the oracle function for 3-qubit Deutsch-Jozsa
func applyOracle3Qubit(b builder.Builder, oracleType string) {
	switch oracleType {
	case "constant_0":
		// f(x1,x2) = 0: Do nothing
		// No gates needed
	case "constant_1":
		// f(x1,x2) = 1: Flip the ancilla qubit for all inputs
		b.X(2)
	case "balanced_xor":
		// f(x1,x2) = x1 ⊕ x2: Apply CNOT(0,2) then CNOT(1,2)
		b.CNOT(0, 2).CNOT(1, 2)
	}
}

// analyzeResults analyzes and displays the results for 2-qubit Deutsch-Jozsa
func analyzeResults(hist map[string]int, shots int, oracleType string) {
	zeroCount := hist["0"]
	oneCount := hist["1"]

	fmt.Printf("Results for %s:\n", oracleType)
	fmt.Printf("  |0⟩: %d counts (%.2f%%)\n", zeroCount, float64(zeroCount)/float64(shots)*100)
	fmt.Printf("  |1⟩: %d counts (%.2f%%)\n", oneCount, float64(oneCount)/float64(shots)*100)

	// Determine if function is constant or balanced
	threshold := int(float64(shots) * 0.9)
	if zeroCount > threshold {
		fmt.Printf("  → Function is CONSTANT (measured |0⟩)\n")
	} else if oneCount > threshold {
		fmt.Printf("  → Function is BALANCED (measured |1⟩)\n")
	} else {
		fmt.Printf("  → Inconclusive result (noise or error)\n")
	}
}

// analyzeResults3Qubit analyzes and displays the results for 3-qubit Deutsch-Jozsa
func analyzeResults3Qubit(hist map[string]int, shots int, oracleType string) {
	zeroZeroCount := hist["00"]
	otherCount := shots - zeroZeroCount

	fmt.Printf("Results for %s:\n", oracleType)
	fmt.Printf("  |00⟩: %d counts (%.2f%%)\n", zeroZeroCount, float64(zeroZeroCount)/float64(shots)*100)
	fmt.Printf("  Other states: %d counts (%.2f%%)\n", otherCount, float64(otherCount)/float64(shots)*100)

	// Determine if function is constant or balanced
	threshold := int(float64(shots) * 0.9)
	if zeroZeroCount > threshold {
		fmt.Printf("  → Function is CONSTANT (measured |00⟩)\n")
	} else if otherCount > threshold {
		fmt.Printf("  → Function is BALANCED (measured non-|00⟩)\n")
	} else {
		fmt.Printf("  → Inconclusive result (noise or error)\n")
	}
}

// reverseString reverses a string to handle bit ordering
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
