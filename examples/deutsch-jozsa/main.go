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
	fmt.Println("\n--- Bernstein-Vazirani Algorithm Demonstrations ---")
	bernsteinVaziraniDemo(shots)
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

// bernsteinVaziraniDemo demonstrates the Bernstein-Vazirani algorithm with different hidden strings
func bernsteinVaziraniDemo(shots int) {
	fmt.Println("\n=== 2-Qubit Bernstein-Vazirani Algorithm ===")

	// Test hidden string "0"
	fmt.Println("\n1. Finding hidden string s = \"0\":")
	bernsteinVazirani2Qubit(shots, "0")

	// Test hidden string "1"
	fmt.Println("\n2. Finding hidden string s = \"1\":")
	bernsteinVazirani2Qubit(shots, "1")

	fmt.Println("\n=== 3-Qubit Bernstein-Vazirani Algorithm ===")

	// Test hidden string "00"
	fmt.Println("\n3. Finding hidden string s = \"00\":")
	bernsteinVazirani3Qubit(shots, "00")

	// Test hidden string "01"
	fmt.Println("\n4. Finding hidden string s = \"01\":")
	bernsteinVazirani3Qubit(shots, "01")

	// Test hidden string "10"
	fmt.Println("\n5. Finding hidden string s = \"10\":")
	bernsteinVazirani3Qubit(shots, "10")

	// Test hidden string "11"
	fmt.Println("\n6. Finding hidden string s = \"11\":")
	bernsteinVazirani3Qubit(shots, "11")

	fmt.Println("\n=== 4-Qubit Bernstein-Vazirani Algorithm ===")

	// Test hidden string "101"
	fmt.Println("\n7. Finding hidden string s = \"101\":")
	bernsteinVazirani4Qubit(shots, "101")

	// Test hidden string "110"
	fmt.Println("\n8. Finding hidden string s = \"110\":")
	bernsteinVazirani4Qubit(shots, "110")
}

// bernsteinVazirani2Qubit implements the 2-qubit Bernstein-Vazirani algorithm
// Uses 2 qubits: 1 input qubit + 1 ancilla qubit for the oracle
func bernsteinVazirani2Qubit(shots int, hiddenString string) {
	// 2 qubits: qubit 0 (input), qubit 1 (ancilla)
	// 1 classical bit to measure the result
	b := builder.New(builder.Q(2), builder.C(1))

	// Initialize ancilla qubit in |1⟩ state
	b.X(1)

	// Apply Hadamard to both qubits
	b.H(0).H(1)

	// Apply oracle function based on hidden string
	applyBVOracle2Qubit(b, hiddenString)

	// Apply Hadamard to input qubit
	b.H(0)

	// Measure input qubit
	b.Measure(0, 0)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building Bernstein-Vazirani circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running Bernstein-Vazirani simulation: %v\n", err)
		return
	}

	// Analyze results
	analyzeBVResults2Qubit(hist, shots, hiddenString)
}

// bernsteinVazirani3Qubit implements the 3-qubit Bernstein-Vazirani algorithm
// Uses 3 qubits: 2 input qubits + 1 ancilla qubit
func bernsteinVazirani3Qubit(shots int, hiddenString string) {
	// 3 qubits: qubit 0,1 (input), qubit 2 (ancilla)
	// 2 classical bits to measure the results
	b := builder.New(builder.Q(3), builder.C(2))

	// Initialize ancilla qubit in |1⟩ state
	b.X(2)

	// Apply Hadamard to all qubits
	b.H(0).H(1).H(2)

	// Apply oracle function based on hidden string
	applyBVOracle3Qubit(b, hiddenString)

	// Apply Hadamard to input qubits
	b.H(0).H(1)

	// Measure input qubits
	b.Measure(0, 0).Measure(1, 1)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building 3-qubit Bernstein-Vazirani circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running 3-qubit Bernstein-Vazirani simulation: %v\n", err)
		return
	}

	// Analyze results
	analyzeBVResults3Qubit(hist, shots, hiddenString)
}

// bernsteinVazirani4Qubit implements the 4-qubit Bernstein-Vazirani algorithm
// Uses 4 qubits: 3 input qubits + 1 ancilla qubit
func bernsteinVazirani4Qubit(shots int, hiddenString string) {
	// 4 qubits: qubit 0,1,2 (input), qubit 3 (ancilla)
	// 3 classical bits to measure the results
	b := builder.New(builder.Q(4), builder.C(3))

	// Initialize ancilla qubit in |1⟩ state
	b.X(3)

	// Apply Hadamard to all qubits
	b.H(0).H(1).H(2).H(3)

	// Apply oracle function based on hidden string
	applyBVOracle4Qubit(b, hiddenString)

	// Apply Hadamard to input qubits
	b.H(0).H(1).H(2)

	// Measure input qubits
	b.Measure(0, 0).Measure(1, 1).Measure(2, 2)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building 4-qubit Bernstein-Vazirani circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running 4-qubit Bernstein-Vazirani simulation: %v\n", err)
		return
	}

	// Analyze results
	analyzeBVResults4Qubit(hist, shots, hiddenString)
}

// applyBVOracle2Qubit applies the Bernstein-Vazirani oracle for 2-qubit case
// The oracle computes f(x) = s·x where s is the hidden string
func applyBVOracle2Qubit(b builder.Builder, hiddenString string) {
	// For each bit position i where s[i] = '1', apply CNOT(i, ancilla)
	if len(hiddenString) >= 1 && hiddenString[0] == '1' {
		b.CNOT(0, 1)
	}
}

// applyBVOracle3Qubit applies the Bernstein-Vazirani oracle for 3-qubit case
// The oracle computes f(x1,x0) = s1*x1 ⊕ s0*x0 where s = s1s0
func applyBVOracle3Qubit(b builder.Builder, hiddenString string) {
	// hiddenString is big-endian s1s0. Qubits are little-endian q1q0.
	// Ancilla is qubit 2.
	if len(hiddenString) >= 2 && hiddenString[1] == '1' { // s0
		b.CNOT(0, 2)
	}
	if len(hiddenString) >= 1 && hiddenString[0] == '1' { // s1
		b.CNOT(1, 2)
	}
}

// applyBVOracle4Qubit applies the Bernstein-Vazirani oracle for 4-qubit case
// The oracle computes f(x2,x1,x0) = s2*x2 ⊕ s1*x1 ⊕ s0*x0 where s = s2s1s0
func applyBVOracle4Qubit(b builder.Builder, hiddenString string) {
	// hiddenString is big-endian s2s1s0. Qubits are little-endian q2q1q0.
	// Ancilla is qubit 3.
	if len(hiddenString) >= 3 && hiddenString[2] == '1' { // s0
		b.CNOT(0, 3)
	}
	if len(hiddenString) >= 2 && hiddenString[1] == '1' { // s1
		b.CNOT(1, 3)
	}
	if len(hiddenString) >= 1 && hiddenString[0] == '1' { // s2
		b.CNOT(2, 3)
	}
}

// analyzeBVResults2Qubit analyzes and displays the results for 2-qubit Bernstein-Vazirani
func analyzeBVResults2Qubit(hist map[string]int, shots int, hiddenString string) {
	fmt.Printf("Results for hidden string \"%s\":\n", hiddenString)

	for state, count := range hist {
		percentage := float64(count) / float64(shots) * 100
		fmt.Printf("  |%s⟩: %d counts (%.2f%%)\n", state, count, percentage)
	}

	// The measured string is the hidden string.
	expectedResult := hiddenString
	if count, exists := hist[expectedResult]; exists && count > int(float64(shots)*0.9) {
		fmt.Printf("  ✓ Successfully found hidden string: \"%s\"\n", expectedResult)
	} else {
		fmt.Printf("  ✗ Failed to find hidden string \"%s\"\n", hiddenString)
	}
}

// analyzeBVResults3Qubit analyzes and displays the results for 3-qubit Bernstein-Vazirani
func analyzeBVResults3Qubit(hist map[string]int, shots int, hiddenString string) {
	fmt.Printf("Results for hidden string \"%s\":\n", hiddenString)

	for state, count := range hist {
		percentage := float64(count) / float64(shots) * 100
		fmt.Printf("  |%s⟩: %d counts (%.2f%%)\n", reverseString(state), count, percentage)
	}

	// The measured string is the reverse of the hidden string due to endianness.
	expectedResult := reverseString(hiddenString)
	if count, exists := hist[expectedResult]; exists && count > int(float64(shots)*0.9) {
		fmt.Printf("  ✓ Successfully found hidden string: \"%s\" (measured as |%s⟩)\n", hiddenString, expectedResult)
	} else {
		fmt.Printf("  ✗ Failed to find hidden string \"%s\"\n", hiddenString)
	}
}

// analyzeBVResults4Qubit analyzes and displays the results for 4-qubit Bernstein-Vazirani
func analyzeBVResults4Qubit(hist map[string]int, shots int, hiddenString string) {
	fmt.Printf("Results for hidden string \"%s\":\n", hiddenString)

	for state, count := range hist {
		percentage := float64(count) / float64(shots) * 100
		fmt.Printf("  |%s⟩: %d counts (%.2f%%)\n", reverseString(state), count, percentage)
	}

	// The measured string is the reverse of the hidden string due to endianness.
	expectedResult := reverseString(hiddenString)
	if count, exists := hist[expectedResult]; exists && count > int(float64(shots)*0.9) {
		fmt.Printf("  ✓ Successfully found hidden string: \"%s\" (measured as |%s⟩)\n", hiddenString, expectedResult)
	} else {
		fmt.Printf("  ✗ Failed to find hidden string \"%s\"\n", hiddenString)
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
