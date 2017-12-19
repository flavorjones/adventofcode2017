package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CpuTrampolineMaze struct {
	instructions []int
	addr         int
	steps        int
}

func NewCpuTrampolineMaze(instruction_list string) *CpuTrampolineMaze {
	var instructions []int
	for _, instruction_entry := range strings.Fields(instruction_list) {
		instruction, err := strconv.Atoi(instruction_entry)
		if err != nil {
			panic(fmt.Sprintf("cannot parse '%s' as an int", instruction_entry))
		}
		instructions = append(instructions, instruction)
	}
	return &CpuTrampolineMaze{instructions: instructions}
}

func (ctm *CpuTrampolineMaze) tick() {
	jump := ctm.instructions[ctm.addr]
	ctm.steps += 1
	ctm.instructions[ctm.addr] += 1
	ctm.addr += jump
	if ctm.addr >= len(ctm.instructions) {
		ctm.addr = -1
	}
}

func (ctm *CpuTrampolineMaze) run() {
	for ctm.addr >= 0 {
		ctm.tick()
	}
}

func (ctm *CpuTrampolineMaze) tick2() {
	jump := ctm.instructions[ctm.addr]
	ctm.steps += 1
	if ctm.instructions[ctm.addr] >= 3 {
		ctm.instructions[ctm.addr] -= 1
	} else {
		ctm.instructions[ctm.addr] += 1
	}
	ctm.addr += jump
	if ctm.addr >= len(ctm.instructions) {
		ctm.addr = -1
	}
}

func (ctm *CpuTrampolineMaze) run2() {
	for ctm.addr >= 0 {
		ctm.tick2()
	}
}

var _ = Describe("Day5", func() {
	Describe("CpuTrampolineMaze", func() {
		var ctm *CpuTrampolineMaze

		BeforeEach(func() {
			ctm = NewCpuTrampolineMaze("0 3 0 1 -3")
		})

		Describe("tick/run", func() {
			It("can be stepped through and examined", func() {
				Expect(ctm.instructions).To(Equal([]int{0, 3, 0, 1, -3}))
				Expect(ctm.addr).To(Equal(0))

				ctm.tick()
				Expect(ctm.instructions).To(Equal([]int{1, 3, 0, 1, -3}))
				Expect(ctm.addr).To(Equal(0))
				Expect(ctm.steps).To(Equal(1))

				ctm.tick()
				Expect(ctm.instructions).To(Equal([]int{2, 3, 0, 1, -3}))
				Expect(ctm.addr).To(Equal(1))
				Expect(ctm.steps).To(Equal(2))

				ctm.tick()
				Expect(ctm.instructions).To(Equal([]int{2, 4, 0, 1, -3}))
				Expect(ctm.addr).To(Equal(4))
				Expect(ctm.steps).To(Equal(3))

				ctm.tick()
				Expect(ctm.instructions).To(Equal([]int{2, 4, 0, 1, -2}))
				Expect(ctm.addr).To(Equal(1))
				Expect(ctm.steps).To(Equal(4))

				ctm.tick()
				Expect(ctm.instructions).To(Equal([]int{2, 5, 0, 1, -2}))
				Expect(ctm.addr).To(Equal(-1))
				Expect(ctm.steps).To(Equal(5))
			})

			It("can be run to completion", func() {
				Expect(ctm.instructions).To(Equal([]int{0, 3, 0, 1, -3}))
				Expect(ctm.addr).To(Equal(0))

				ctm.run()
				Expect(ctm.instructions).To(Equal([]int{2, 5, 0, 1, -2}))
				Expect(ctm.addr).To(Equal(-1))
				Expect(ctm.steps).To(Equal(5))
			})
		})

		Describe("tick2/run2", func() {
			It("decreases instruction by 1 if greater than 3", func() {
				ctm.run2()
				Expect(ctm.instructions).To(Equal([]int{2, 3, 2, 3, -1}))
				Expect(ctm.steps).To(Equal(10))
			})
		})
	})

	Describe("puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day5.txt")
		instruction_list := string(raw_data)

		It("answers star 1", func() {
			ctm := NewCpuTrampolineMaze(instruction_list)
			ctm.run()
			fmt.Printf("d5 s1: exiting took %d steps\n", ctm.steps)
		})

		It("answers star 2", func() {
			ctm := NewCpuTrampolineMaze(instruction_list)
			ctm.run2()
			fmt.Printf("d5 s2: exiting took %d steps\n", ctm.steps)
		})
	})
})
