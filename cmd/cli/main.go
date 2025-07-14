package main

import (
	"fmt"
	"sort" // Import the sort package

	"github.com/kegliz/qcm/qc/builder"
	"github.com/kegliz/qcm/qc/simulator"
	_ "github.com/kegliz/qcm/qc/simulator/itsu"
	_ "github.com/kegliz/qcm/qc/simulator/qsim"
)

func main() {
	shots := 1024

	fmt.Println("--- Bell State Simulation ---")
	simulateBellState(shots)
	fmt.Println("\n--- 2-Qubit Grover Simulation (|11>) ---")
	simulateGrover2Qubit(shots)
	fmt.Println("\n--- 3-Qubit Grover Simulation (|111>) - 2 iterations (optimal) ---")
	simulateGrover3Qubit(shots)
	fmt.Println("\n--- 4-Qubit Grover Simulation (|1111>) - 3 iterations (optimal) ---")
	simulateGrover4Qubit(shots)
}

// simulateBellState prepares the |Φ⁺⟩ Bell state and checks ~50/50 statistics.
func simulateBellState(shots int) {
	b := builder.New(builder.Q(2), builder.C(2))
	b.H(0).CNOT(0, 1).Measure(0, 0).Measure(1, 1)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building Bell state circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("itsu", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator with itsu: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running Bell state simulation: %v\n", err)
		return
	}

	pretty(hist, shots)
}

// simulateGrover2Qubit demonstrates one Grover iteration on 2‑qubit search space
// amplifying the |11⟩ state.
func simulateGrover2Qubit(shots int) {
	b := builder.New(builder.Q(2), builder.C(2))

	// — initial superposition —
	b.H(0).H(1)

	// — oracle marks |11⟩ by phase flip (controlled‑Z) —
	b.CZ(0, 1)

	// — diffusion operator —
	b.H(0).H(1)
	b.X(0).X(1)
	b.CZ(0, 1)
	b.X(0).X(1)
	b.H(0).H(1)

	// — measurement —
	b.Measure(0, 0).Measure(1, 1)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building 2-qubit Grover circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator with qsim: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running 2-qubit Grover simulation: %v\n", err)
		return
	}

	pretty(hist, shots)
}

// simulateGrover3Qubit demonstrates optimal Grover iterations (2) on 3‑qubit search space
// amplifying the |111⟩ state.
func simulateGrover3Qubit(shots int) {
	b := builder.New(builder.Q(3), builder.C(3))

	// — initial superposition —
	b.H(0).H(1).H(2)

	// Perform 2 Grover iterations (optimal for 3 qubits: π/4 * √8 ≈ 2.22)
	for range 2 {
		// — oracle marks |111⟩ by phase flip (CCZ) —
		// Implement CCZ using H and Toffoli: H(target) Toffoli(c1, c2, target) H(target)
		b.H(2).Toffoli(0, 1, 2).H(2)

		// — diffusion operator (3 qubits) —
		// HHH - XXX - CCZ - XXX - HHH
		b.H(0).H(1).H(2)
		b.X(0).X(1).X(2)
		// CCZ
		b.H(2).Toffoli(0, 1, 2).H(2)
		b.X(0).X(1).X(2)
		b.H(0).H(1).H(2)
	}

	// — measurement —
	b.Measure(0, 0).Measure(1, 1).Measure(2, 2)

	c, err := b.BuildCircuit()

	if err != nil {
		fmt.Printf("Error building 3-qubit Grover circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator with qsim: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running 3-qubit Grover simulation: %v\n", err)
		return
	}

	pretty(hist, shots)
}

// simulateGrover4Qubit demonstrates optimal Grover iterations (3) on 4‑qubit search space
// amplifying the |1111⟩ state.
func simulateGrover4Qubit(shots int) {
	// This circuit uses a CCCZ gate
	// so we need an ancilla qubit to implement the CCCZ gate with H and Toffoli gates.

	b := builder.New(builder.Q(5), builder.C(4))

	// — initial superposition —
	b.H(0).H(1).H(2).H(3)

	// Perform 3 Grover iterations (optimal for 4 qubits: π/4 * √16 ≈ 3.14)
	for range 3 {
		// — oracle marks |1111⟩ by phase flip (CCCZ) —
		// CCCZ using H and Toffoli:  H(3) - CCCX - H(3)
		b.H(3)
		//  CCCX using ancilla qubit 4:
		//  Toffoli(0,1,4) - Toffoli(2,4,3) - Toffoli(0,1,4)
		b.Toffoli(0, 1, 4).Toffoli(2, 4, 3).Toffoli(0, 1, 4)
		b.H(3)

		// — diffusion operator (4 qubits) —
		// HHHH - XXXX - CCCZ - XXXX - HHHH
		b.H(0).H(1).H(2).H(3)
		b.X(0).X(1).X(2).X(3)
		//  CCCZ using H and Toffoli:  H(3) - CCCX - H(3)
		b.H(3)
		// CCCX using ancilla qubit 4:
		// Toffoli(0,1,4) - Toffoli(2,4,3) - Toffoli(0,1,4)
		b.Toffoli(0, 1, 4).Toffoli(2, 4, 3).Toffoli(0, 1, 4)
		b.H(3)

		b.X(0).X(1).X(2).X(3)
		b.H(0).H(1).H(2).H(3)
	}

	// — measurement —
	b.Measure(0, 0).Measure(1, 1).Measure(2, 2).Measure(3, 3)

	c, err := b.BuildCircuit()
	if err != nil {
		fmt.Printf("Error building 4-qubit Grover circuit: %v\n", err)
		return
	}

	sim, err := simulator.NewSimulatorWithRunner("qsim", simulator.SimulatorOptions{Shots: shots})
	if err != nil {
		fmt.Printf("Error creating simulator with qsim: %v\n", err)
		return
	}

	hist, err := sim.Run(c)
	if err != nil {
		fmt.Printf("Error running 4-qubit Grover simulation: %v\n", err)
		return
	}

	pretty(hist, shots)
}

// pretty prints the histogram results in a readable, sorted format
func pretty(hist map[string]int, shots int) {
	// Extract keys for sorting
	keys := make([]string, 0, len(hist))
	for k := range hist {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Sort keys alphabetically

	// Print sorted results
	for _, state := range keys {
		count := hist[state]
		probability := float64(count) / float64(shots)
		fmt.Printf("State |%s>: %d counts (%.2f%%)\n", state, count, probability*100)
	}
}
