package main

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/kegliz/qcm/qc/builder"
	"github.com/kegliz/qcm/qc/simulator"
	"github.com/stretchr/testify/assert"
)

// checkOracle computes the output of an oracle for a given input state.
func checkOracle(t *testing.T, nInput int, oracleFunc func(b builder.Builder), input string) string {
	t.Helper()
	nAncilla := 1
	nQubits := nInput + nAncilla
	b := builder.New(builder.Q(nQubits))

	// Prepare input state |x⟩|0⟩.
	// The input string is big-endian, but the qubits are little-endian.
	// We need to reverse the input string to match the qubit order.
	inputLE := reverseString(input)
	for i, bit := range inputLE {
		if bit == '1' {
			b.X(i) // Little-endian qubit order
		}
	}

	// Apply oracle
	oracleFunc(b)

	c, err := b.BuildCircuit()
	assert.NoError(t, err)

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{})
	assert.NoError(t, err)

	// Get statevector directly from simulator
	sv, err := sim.GetStatevector(c)
	assert.NoError(t, err)

	// Find the non-zero amplitude state
	nonZeroStates := []int{}
	for i, amp := range sv {
		if cmplx.Abs(amp) > 1e-9 {
			nonZeroStates = append(nonZeroStates, i)
		}
	}
	assert.Len(t, nonZeroStates, 1, "Expected exactly one non-zero amplitude")

	foundState := nonZeroStates[0]

	// The oracle computes |x⟩|y⟩ -> |x⟩|y ⊕ f(x)⟩.
	// We start with |x⟩|0⟩, so we get |x⟩|f(x)⟩.
	// The output f(x) is on the ancilla qubit, which is at index nInput.
	outputBit := (foundState >> nInput) & 1
	return fmt.Sprintf("%d", outputBit)
}

func TestBernsteinVaziraniOracles(t *testing.T) {
	// 2-Qubit Oracles (1 input qubit)
	testCases2Qubit := []struct {
		name         string
		hiddenString string
		input        string
		expected     string
	}{
		{"s=0, x=0", "0", "0", "0"},
		{"s=0, x=1", "0", "1", "0"},
		{"s=1, x=0", "1", "0", "0"},
		{"s=1, x=1", "1", "1", "1"},
	}

	for _, tc := range testCases2Qubit {
		t.Run(fmt.Sprintf("BV_2Qubit_%s", tc.name), func(t *testing.T) {
			oracleFunc := func(b builder.Builder) { applyBVOracle2Qubit(b, tc.hiddenString) }
			output := checkOracle(t, 1, oracleFunc, tc.input)
			assert.Equal(t, tc.expected, output)
		})
	}

	// 3-Qubit Oracles (2 input qubits)
	testCases3Qubit := []struct {
		name         string
		hiddenString string
		input        string
		expected     string
	}{
		{"s=01, x=00", "01", "00", "0"},
		{"s=01, x=01", "01", "01", "1"},
		{"s=01, x=10", "01", "10", "0"},
		{"s=01, x=11", "01", "11", "1"},
		{"s=11, x=00", "11", "00", "0"},
		{"s=11, x=01", "11", "01", "1"},
		{"s=11, x=10", "11", "10", "1"},
		{"s=11, x=11", "11", "11", "0"},
	}

	for _, tc := range testCases3Qubit {
		t.Run(fmt.Sprintf("BV_3Qubit_%s", tc.name), func(t *testing.T) {
			oracleFunc := func(b builder.Builder) { applyBVOracle3Qubit(b, tc.hiddenString) }
			output := checkOracle(t, 2, oracleFunc, tc.input)
			assert.Equal(t, tc.expected, output)
		})
	}

	// 4-Qubit Oracles (3 input qubits)
	testCases4Qubit := []struct {
		name         string
		hiddenString string
		input        string
		expected     string
	}{
		{"s=101, x=000", "101", "000", "0"},
		{"s=101, x=111", "101", "111", "0"}, // 1*1 + 0*1 + 1*1 = 0 mod 2
		{"s=101, x=101", "101", "101", "0"}, // 1*1 + 0*0 + 1*1 = 0 mod 2
		{"s=101, x=010", "101", "010", "0"},
		{"s=101, x=100", "101", "100", "1"},
		{"s=101, x=001", "101", "001", "1"},
	}

	for _, tc := range testCases4Qubit {
		t.Run(fmt.Sprintf("BV_4Qubit_%s", tc.name), func(t *testing.T) {
			oracleFunc := func(b builder.Builder) { applyBVOracle4Qubit(b, tc.hiddenString) }
			output := checkOracle(t, 3, oracleFunc, tc.input)
			assert.Equal(t, tc.expected, output)
		})
	}
}
