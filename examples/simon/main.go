package main

import (
	"fmt"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/kegliz/qcm/qc/builder"
	"github.com/kegliz/qcm/qc/simulator"
	_ "github.com/kegliz/qcm/qc/simulator/qsim"
)

func main() {
	shots := 1024 * 10

	fmt.Println("\n--- Simon's Algorithm Demonstrations ---")

	// First demonstrate oracle mappings
	demonstrateOracleMappings()

	// Then run the algorithm
	simonDemo(shots)
}

func simonDemo(shots int) {
	fmt.Println("\n=== 2-Qubit Simon's Algorithm ===")
	fmt.Println("\n1. Testing function with no secret (s=\"00\"):")
	simonAlgorithm2Qubit(shots, "00")
	fmt.Println("\n2. Testing function with secret string s = \"01\":")
	simonAlgorithm2Qubit(shots, "01")
	fmt.Println("\n3. Testing function with secret string s = \"10\":")
	simonAlgorithm2Qubit(shots, "10")
	fmt.Println("\n4. Testing function with secret string s = \"11\":")
	simonAlgorithm2Qubit(shots, "11")

	fmt.Println("\n=== 3-Qubit Simon's Algorithm ===")
	fmt.Println("\n5. Testing function with no secret (s=\"000\"):")
	simonAlgorithm3Qubit(shots, "000")
	fmt.Println("\n6. Testing function with secret string s = \"110\":")
	simonAlgorithm3Qubit(shots, "110")
	fmt.Println("\n7. Testing function with secret string s = \"101\":")
	simonAlgorithm3Qubit(shots, "101")
	fmt.Println("\n8. Testing function with secret string s = \"011\":")
	simonAlgorithm3Qubit(shots, "011")
}

// simonAlgorithm2Qubit runs Simon's algorithm for 2 qubits with the given secret string
// Uses 4 qubits: 2 input qubits + 2 ancilla qubits
// secretString is the hidden string s in big-endian format
func simonAlgorithm2Qubit(shots int, secretString string) {
	b := builder.New(builder.Q(4), builder.C(2))
	b.H(0).H(1)
	applySimonOracle2Qubit(b, secretString)
	b.H(0).H(1)
	b.Measure(0, 0).Measure(1, 1)
	c, _ := b.BuildCircuit()
	sim, _ := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	hist, _ := sim.Run(c)
	analyzeSimonResults(2, hist, shots, secretString)
}

// simonAlgorithm3Qubit runs Simon's algorithm for 3 qubits with the given secret string
// Uses 6 qubits: 3 input qubits + 3 ancilla qubits
// secretString is the hidden string s in big-endian format
func simonAlgorithm3Qubit(shots int, secretString string) {
	b := builder.New(builder.Q(6), builder.C(3))
	b.H(0).H(1).H(2)
	applySimonOracle3Qubit(b, secretString)
	b.H(0).H(1).H(2)
	b.Measure(0, 0).Measure(1, 1).Measure(2, 2)
	c, _ := b.BuildCircuit()
	sim, _ := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	hist, _ := sim.Run(c)
	analyzeSimonResults(3, hist, shots, secretString)
}

// applySimonOracle2Qubit applies the Simon oracle for 2 qubits
// For s ≠ 0: f(x) = f(y) ⟺ x ⊕ y ∈ {0, s}
func applySimonOracle2Qubit(b builder.Builder, secretString string) {
	switch secretString {
	case "00": // f(x) = x, one-to-one
		b.CNOT(0, 2).CNOT(1, 3)
	case "01": // s_BE="01" -> s_LE="10". f(x) = (x[1], 0) - only x[1] matters
		b.CNOT(1, 2)
	case "10": // s_BE="10" -> s_LE="01". f(x) = (x[0], 0) - only x[0] matters
		b.CNOT(0, 2)
	case "11": // s_BE="11" -> s_LE="11". f(x) = (x[0] ⊕ x[1], 0) - XOR matters
		b.CNOT(0, 2).CNOT(1, 2)
	}
}

// applySimonOracle3Qubit applies the Simon oracle for 3 qubits
// For s ≠ 0: f(x) = f(y) ⟺ x ⊕ y ∈ {0, s}
func applySimonOracle3Qubit(b builder.Builder, secretString string) {
	switch secretString {
	case "000": // f(x) = x, one-to-one
		b.CNOT(0, 3).CNOT(1, 4).CNOT(2, 5)
	case "110": // s_BE="110" -> s_LE="011". f(x) = (x[0], x[1] ⊕ x[2], 0)
		b.CNOT(0, 3)            // f[0] = x[0]
		b.CNOT(1, 4).CNOT(2, 4) // f[1] = x[1] ⊕ x[2]
		// f[2] = 0 (no CNOT to output qubit 5)
	case "101": // s_BE="101" -> s_LE="101". f(x) = (x[0] ⊕ x[2], x[1], 0)
		b.CNOT(0, 3).CNOT(2, 3) // f[0] = x[0] ⊕ x[2]
		b.CNOT(1, 4)            // f[1] = x[1]
		// f[2] = 0 (no CNOT to output qubit 5)
	case "011": // s_BE="011" -> s_LE="110". f(x) = (x[0] ⊕ x[1], 0, x[2])
		b.CNOT(0, 3).CNOT(1, 3) // f[0] = x[0] ⊕ x[1]
		// f[1] = 0 (no CNOT to output qubit 4)
		b.CNOT(2, 5) // f[2] = x[2]
	}
}

func analyzeSimonResults(n int, hist map[string]int, shots int, secretString string) {
	fmt.Printf("Results for secret string \"%s\":\n", secretString)
	for state, count := range hist {
		percentage := float64(count) / float64(shots) * 100
		fmt.Printf("  |%s⟩: %d counts (%.2f%%)\n", state, count, percentage)
	}

	s_val, _ := strconv.ParseInt(secretString, 2, 64)
	if s_val == 0 {
		fmt.Printf("  → Function is one-to-one (no secret string)\n")
	} else {
		fmt.Printf("  → Function has secret string \"%s\"\n", secretString)
		var expectedStates []string
		// secretString is BE, y in the loop is BE. Dot product needs consistent endianness.
		for i := 0; i < (1 << n); i++ {
			y_str := fmt.Sprintf(fmt.Sprintf("%%0%db", n), i)
			dot := 0
			for j := range n {
				if y_str[j] == '1' && secretString[j] == '1' {
					dot++
				}
			}
			if dot%2 == 0 {
				expectedStates = append(expectedStates, y_str)
			}
		}
		fmt.Printf("  → Expected states (y where y·s = 0 mod 2): %v\n", expectedStates)
	}
}

// demonstrateOracleMappings shows explicit f(x) mappings for all Simon oracles
func demonstrateOracleMappings() {
	fmt.Println("\n=== Oracle Function Mappings ===")

	// 2-Qubit Oracles
	fmt.Println("\n--- 2-Qubit Oracles ---")
	secretStrings2 := []string{"00", "01", "10", "11"}
	for _, secretString := range secretStrings2 {
		demonstrateAndAnalyzeOracle(2, secretString)
	}

	// 3-Qubit Oracles
	fmt.Println("\n--- 3-Qubit Oracles ---")
	secretStrings3 := []string{"000", "110", "101", "011"}
	for _, secretString := range secretStrings3 {
		demonstrateAndAnalyzeOracle(3, secretString)
	}
}

// demonstrateAndAnalyzeOracle generates all mappings and analyzes properties in one pass
func demonstrateAndAnalyzeOracle(n int, secretString string) {
	fmt.Printf("\nOracle for secret string s = \"%s\":\n", secretString)
	fmt.Printf("Function f(x) mappings:\n")

	numInputs := 1 << n
	outputs := make(map[string][]string) // For analyzing two-to-one property

	// Generate all mappings and collect outputs for analysis
	for input := range numInputs {
		inputStr := fmt.Sprintf(fmt.Sprintf("%%0%db", n), input)
		output := getOracleMapping(n, secretString, inputStr)
		fmt.Printf("  |%s⟩ → |%s⟩\n", inputStr, output)

		// Collect for two-to-one analysis
		outputs[output] = append(outputs[output], inputStr)
	}

	// Analyze two-to-one property using collected data
	analyzeTwoToOnePropertyFromMappings(n, secretString, outputs)
}

// analyzeTwoToOnePropertyFromMappings analyzes using pre-computed mappings
func analyzeTwoToOnePropertyFromMappings(n int, secretString string, outputs map[string][]string) {
	fmt.Printf("Two-to-one property analysis for s = \"%s\":\n", secretString)

	s_val, _ := strconv.ParseInt(secretString, 2, 64)
	if s_val == 0 {
		fmt.Printf("  → s = %s: Function should be one-to-one\n", secretString)
		// Check if it's actually one-to-one
		allDistinct := true
		for output, inputs := range outputs {
			if len(inputs) > 1 {
				fmt.Printf("    ✗ f(%s) all map to %s (not one-to-one)\n", strings.Join(inputs, ", "), output)
				allDistinct = false
			}
		}
		if allDistinct {
			fmt.Printf("    ✓ Function is one-to-one\n")
		}
		return
	}

	fmt.Printf("  → s ≠ 0: Function should be exactly two-to-one\n")
	fmt.Printf("  → Checking f(x) = f(y) ⟺ x ⊕ y ∈ {0, s}:\n")

	allCorrect := true
	for output, inputs := range outputs {
		if len(inputs) != 2 {
			fmt.Printf("    ✗ Output %s has %d inputs: %v (should be exactly 2)\n", output, len(inputs), inputs)
			allCorrect = false
			continue
		}

		// Check if the two inputs differ by s
		x1, _ := strconv.ParseInt(inputs[0], 2, 64)
		x2, _ := strconv.ParseInt(inputs[1], 2, 64)
		diff := x1 ^ x2

		if diff == s_val {
			fmt.Printf("    ✓ f(%s) = f(%s) = %s, %s ⊕ %s = %s = s\n",
				inputs[0], inputs[1], output, inputs[0], inputs[1],
				fmt.Sprintf(fmt.Sprintf("%%0%db", n), diff))
		} else {
			fmt.Printf("    ✗ f(%s) = f(%s) = %s, but %s ⊕ %s = %s ≠ s\n",
				inputs[0], inputs[1], output, inputs[0], inputs[1],
				fmt.Sprintf(fmt.Sprintf("%%0%db", n), diff))
			allCorrect = false
		}
	}

	if allCorrect {
		fmt.Printf("  → ✓ Function is properly two-to-one\n")
	} else {
		fmt.Printf("  → ✗ Function does not satisfy two-to-one property\n")
	}
}

// getOracleMapping calculates the output of the oracle for a given input using the quantum circuit
func getOracleMapping(n int, secretString string, input string) string {
	b := builder.New(builder.Q(2 * n))

	// Prepare input state |input⟩|0...0⟩
	// The input string is big-endian, builder uses little-endian qubits
	for i := range n {
		if input[n-1-i] == '1' {
			b.X(i)
		}
	}

	// Apply the oracle based on secret string and number of qubits
	switch n {
	case 2:
		applySimonOracle2Qubit(b, secretString)
	case 3:
		applySimonOracle3Qubit(b, secretString)
	default:
		return "ERROR"
	}

	c, err := b.BuildCircuit()
	if err != nil {
		return "ERROR"
	}

	// Get statevector
	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{})
	if err != nil {
		return "ERROR"
	}

	sv, err := sim.GetStatevector(c)
	if err != nil {
		return "ERROR"
	}

	// Find the non-zero amplitude state
	for i, amp := range sv {
		if cmplx.Abs(amp) > 1e-9 {
			// Extract output part (bits n to 2n-1)
			outputIntLE := i >> n
			// Convert to big-endian string
			var sb strings.Builder
			for j := n - 1; j >= 0; j-- {
				if (outputIntLE>>j)&1 == 1 {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
			return sb.String()
		}
	}
	return "ERROR"
}
