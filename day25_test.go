package adventofcode2017_test

import (
	. "adventofcode2017"
	"fmt"
	"io/ioutil"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Day25", func() {
	Describe("TuringMachine", func() {
		var testInput = heredoc.Doc(`
			Begin in state A.
			Perform a diagnostic checksum after 6 steps.
			
			In state A:
			  If the current value is 0:
			    - Write the value 1.
			    - Move one slot to the right.
			    - Continue with state B.
			  If the current value is 1:
			    - Write the value 0.
			    - Move one slot to the left.
			    - Continue with state B.
			
			In state B:
			  If the current value is 0:
			    - Write the value 1.
			    - Move one slot to the left.
			    - Continue with state A.
			  If the current value is 1:
			    - Write the value 1.
			    - Move one slot to the right.
			    - Continue with state A.
		`)

		Describe("NewTuringMachine", func() {
			It("parses the input into starting state and instructions", func() {
				tm := NewTuringMachine(testInput)
				Expect(tm.NextState()).To(Equal("A"))
				Expect(tm.Position()).To(Equal(0))
				Expect(tm.StepsRemaining()).To(Equal(6))

				Expect(tm.State("A")).To(Equal(TuringMachineState{[2]TuringMachineInstruction{
					TuringMachineInstruction{Write: 1, Move: TmRight, NextState: TuringMachineStateName("B")},
					TuringMachineInstruction{Write: 0, Move: TmLeft, NextState: TuringMachineStateName("B")},
				}}))
				Expect(tm.State("B")).To(Equal(TuringMachineState{[2]TuringMachineInstruction{
					TuringMachineInstruction{Write: 1, Move: TmLeft, NextState: TuringMachineStateName("A")},
					TuringMachineInstruction{Write: 1, Move: TmRight, NextState: TuringMachineStateName("A")},
				}}))
			})
		})

		Describe("Step()", func() {
			It("moves through the current state into the next state", func() {
				tm := NewTuringMachine(testInput)

				tm.Step()
				Expect(tm.NextState()).To(Equal("B"))
				Expect(tm.Position()).To(Equal(1))
				Expect(tm.TapeAt(0)).To(Equal(1))
				Expect(tm.StepsRemaining()).To(Equal(5))

				tm.Step()
				Expect(tm.NextState()).To(Equal("A"))
				Expect(tm.Position()).To(Equal(0))
				Expect(tm.TapeAt(1)).To(Equal(1))
				Expect(tm.StepsRemaining()).To(Equal(4))

				tm.Step()
				Expect(tm.NextState()).To(Equal("B"))
				Expect(tm.Position()).To(Equal(-1))
				Expect(tm.TapeAt(0)).To(Equal(0))
				Expect(tm.StepsRemaining()).To(Equal(3))

				tm.Step()
				Expect(tm.NextState()).To(Equal("A"))
				Expect(tm.Position()).To(Equal(-2))
				Expect(tm.TapeAt(-1)).To(Equal(1))
				Expect(tm.StepsRemaining()).To(Equal(2))

				tm.Step()
				Expect(tm.NextState()).To(Equal("B"))
				Expect(tm.Position()).To(Equal(-1))
				Expect(tm.TapeAt(-2)).To(Equal(1))
				Expect(tm.StepsRemaining()).To(Equal(1))

				tm.Step()
				Expect(tm.NextState()).To(Equal("A"))
				Expect(tm.Position()).To(Equal(0))
				Expect(tm.TapeAt(-1)).To(Equal(1))
				Expect(tm.StepsRemaining()).To(Equal(0))
			})
		})

		Describe("Checksum()", func() {
			It("counts the number of 1s on tape", func() {
				tm := NewTuringMachine(testInput)
				Expect(tm.Checksum()).To(Equal(0))
				tm.Step()
				Expect(tm.Checksum()).To(Equal(1))
				tm.Step()
				Expect(tm.Checksum()).To(Equal(2))
				tm.Step()
				Expect(tm.Checksum()).To(Equal(1))
				tm.Step()
				Expect(tm.Checksum()).To(Equal(2))
				tm.Step()
				Expect(tm.Checksum()).To(Equal(3))
				tm.Step()
				Expect(tm.Checksum()).To(Equal(3))
			})
		})

		Describe("Run()", func() {
			It("runs Step() until steps remaining is zero", func() {
				tm := NewTuringMachine(testInput)
				tm.Run()
				Expect(tm.StepsRemaining()).To(Equal(0))
				Expect(tm.Position()).To(Equal(0))
				Expect(tm.Checksum()).To(Equal(3))
			})
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day25.txt")
		blueprint := string(rawData)

		It("solves star 1", func() {
			tm := NewTuringMachine(blueprint)
			tm.Run()
			checksum := tm.Checksum()
			fmt.Printf("d25 s1: checksum is %d\n", checksum)
		})
	})
})
