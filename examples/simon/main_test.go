package main

import (
	"fmt"
	"math/cmplx"
	"strconv"
	"strings"
	"testing"

	"github.com/kegliz/qcm/qc/builder"
	"github.com/kegliz/qcm/qc/simulator"
	"github.com/stretchr/testify/assert"
)

// getOracleOutput simulates the oracle with a given input and returns the output part of the state.
func getOracleOutput(t *testing.T, n int, oracle func(b builder.Builder), input string) (string, error) {
	t.Helper()
	// 2n qubits for input and output registers.
	b := builder.New(builder.Q(2 * n))

	// Prepare input state |input⟩|0...0⟩
	// The input string is big-endian, builder uses little-endian qubits.
	for i := range n {
		if input[n-1-i] == '1' {
			b.X(i)
		}
	}

	// Apply the oracle
	oracle(b)

	c, err := b.BuildCircuit()
	if err != nil {
		return "", fmt.Errorf("Failed to build circuit: %v", err)
	}

	// We use a statevector simulator to get the exact output state.
	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{})
	if err != nil {
		return "", fmt.Errorf("Failed to create simulator: %v", err)
	}

	sv, err := sim.GetStatevector(c)
	if err != nil {
		return "", fmt.Errorf("Failed to get statevector: %v", err)
	}

	nonZeroStates := []int{}
	for i, amp := range sv {
		if cmplx.Abs(amp) > 1e-9 {
			nonZeroStates = append(nonZeroStates, i)
		}
	}
	if len(nonZeroStates) != 1 {
		return "", fmt.Errorf("expected one non-zero amplitude, found %d for input %s", len(nonZeroStates), input)
	}
	foundState := nonZeroStates[0]

	// Verify the input part of the state was not changed by the oracle.
	inputIntLE_prepared := 0
	for i := range n {
		if input[n-1-i] == '1' {
			inputIntLE_prepared |= (1 << i)
		}
	}

	if (foundState & ((1 << n) - 1)) != inputIntLE_prepared {
		actualInput := foundState & ((1 << n) - 1)
		return "", fmt.Errorf("input state changed by oracle. Expected LE %d (%s), got LE %d", inputIntLE_prepared, input, actualInput)
	}

	outputIntLE := foundState >> n
	// output is little-endian integer. Format as big-endian string.
	var sb strings.Builder
	for i := n - 1; i >= 0; i-- {
		if (outputIntLE>>i)&1 == 1 {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	outputStrBE := sb.String()

	return outputStrBE, nil
}

// checkOracleProperties verifies that the oracle for a given secret string 's'
// satisfies f(x) = f(x ⊕ s) for all x.
func checkOracleProperties(t *testing.T, n int, secretString string, oracle func(b builder.Builder)) {
	t.Helper()
	numInputs := 1 << n
	outputs := make(map[string]string)

	// Get the oracle output for all possible inputs.
	for i := range numInputs {
		input := fmt.Sprintf("%0*b", n, i)
		output, err := getOracleOutput(t, n, oracle, input)
		if err != nil {
			t.Fatalf("Error getting oracle output for input %s: %v", input, err)
		}
		outputs[input] = output
	}

	s, _ := strconv.ParseInt(secretString, 2, n+1)

	if s == 0 {
		// For s=0, a valid Simon oracle function must be one-to-one.
		outputSet := make(map[string]bool)
		for input, output := range outputs {
			assert.False(t, outputSet[output], "For s=0, f(x) must be one-to-one, but f(%s)=%s is a duplicate output", input, output)
			outputSet[output] = true
		}
	} else {
		// For s!=0, f(x) must equal f(x ⊕ s).
		for i := range numInputs {
			x := i
			x_xor_s := i ^ int(s)
			input_x := fmt.Sprintf("%0*b", n, x)
			input_x_xor_s := fmt.Sprintf("%0*b", n, x_xor_s)
			assert.Equal(t, outputs[input_x], outputs[input_x_xor_s], "f(%s)=%s != f(%s)=%s for s=%s", input_x, outputs[input_x], input_x_xor_s, outputs[input_x_xor_s], secretString)
		}
	}
}

func TestApplySimonOracle2Qubit(t *testing.T) {
	testCases := []string{"00", "01", "10", "11"}
	for _, s := range testCases {
		t.Run(fmt.Sprintf("s=%s", s), func(t *testing.T) {
			oracle := func(b builder.Builder) {
				applySimonOracle2Qubit(b, s)
			}
			checkOracleProperties(t, 2, s, oracle)
		})
	}
}

func TestApplySimonOracle3Qubit(t *testing.T) {
	testCases := []string{"000", "110", "101", "011"}
	for _, s := range testCases {
		t.Run(fmt.Sprintf("s=%s", s), func(t *testing.T) {
			oracle := func(b builder.Builder) {
				applySimonOracle3Qubit(b, s)
			}
			checkOracleProperties(t, 3, s, oracle)
		})
	}
}
