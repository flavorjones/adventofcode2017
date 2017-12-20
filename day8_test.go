package adventofcode2017_test

import (
	"fmt"
	"regexp"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var instructionRe = regexp.MustCompile(`^(\w+) (\w+) ([-\w]+)`)

type RegisterSet map[string]int // register name → register value

func NewRegisterSet() RegisterSet {
	rval := make(RegisterSet)
	return rval
}

func (rs RegisterSet) execInstruction(instruction string) {
	matches := instructionRe.FindStringSubmatch(instruction)
	if len(matches) == 0 {
		panic(fmt.Sprintf("error: could not parse instruction `%s`", instruction))
	}

	registerName := matches[1]
	operator := matches[2]
	operandStr := matches[3]

	operand, err := strconv.Atoi(operandStr)
	if err != nil {
		panic(fmt.Sprintf("error: cannot parse `%s` as an int", operandStr))
	}

	_, ok := rs[registerName]
	if !ok {
		rs[registerName] = 0
	}

	switch operator {
	case "inc":
		rs[registerName] += operand
	case "dec":
		rs[registerName] -= operand
	default:
		panic(fmt.Sprintf("error: unrecognized operator `%s`", operator))
	}

}

var _ = Describe("Day8", func() {
	Describe("RegisterSet", func() {
		var rs RegisterSet

		BeforeEach(func() {
			rs = NewRegisterSet()
		})

		Describe("execInstruction", func() {
			It("creates a register when it encounters a reference", func() {
				_, ok := rs["a"]
				Expect(ok).To(BeFalse())

				rs.execInstruction("a inc 1 if b < 5")

				_, ok = rs["a"]
				Expect(ok).To(BeTrue())
			})

			It("increments a register by a positive integer", func() {
				rs.execInstruction("a inc 11 if b < 5")
				aVal, _ := rs["a"]
				Expect(aVal).To(Equal(11))
			})

			It("increments a register by a negative integer", func() {
				rs.execInstruction("a inc -13 if b < 5")
				aVal, _ := rs["a"]
				Expect(aVal).To(Equal(-13))
			})

			It("decrements a register by a positive integer", func() {
				rs.execInstruction("a dec 15 if b < 5")
				aVal, _ := rs["a"]
				Expect(aVal).To(Equal(-15))
			})

			It("decrements a register by a negative integer", func() {
				rs.execInstruction("a dec -17 if b < 5")
				aVal, _ := rs["a"]
				Expect(aVal).To(Equal(17))
			})

			Describe("ignoring an instruction based on the conditional", func() {
				It("<", func() {

				})

				It("<=", func() {

				})

				It(">", func() {

				})

				It(">=", func() {

				})

				It("==", func() {

				})

				It("!=", func() {

				})
			})
		})
	})
})
