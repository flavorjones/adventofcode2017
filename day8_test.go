package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var instructionRe = regexp.MustCompile(`^(\w+) (\w+) ([-\w]+) if (\w+) (.*) ([-\w]+)$`)

type RegisterSet map[string]int // register name → register value

func NewRegisterSet() RegisterSet {
	rval := make(RegisterSet)
	return rval
}

func (rs RegisterSet) ensureRegister(registerName string) {
	_, ok := rs[registerName]
	if !ok {
		rs[registerName] = 0
	}
}

func (rs RegisterSet) execInstruction(instruction string) {
	matches := instructionRe.FindStringSubmatch(instruction)
	if len(matches) == 0 {
		panic(fmt.Sprintf("error: could not parse instruction `%s`", instruction))
	}

	registerName := matches[1]
	operator := matches[2]
	operandStr := matches[3]
	predSubject := matches[4]
	predicate := matches[5]
	predOperandStr := matches[6]

	operand, err := strconv.Atoi(operandStr)
	if err != nil {
		panic(fmt.Sprintf("error: cannot parse `%s` as an int", operandStr))
	}

	predOperand, err := strconv.Atoi(predOperandStr)
	if err != nil {
		panic(fmt.Sprintf("error: cannot parse `%s` as an int", predOperandStr))
	}

	rs.ensureRegister(registerName)
	rs.ensureRegister(predSubject)

	// test the predicate
	value := rs[predSubject]
	predVal := false
	switch predicate {
	case "<":
		if value < predOperand {
			predVal = true
		}
	case "<=":
		if value <= predOperand {
			predVal = true
		}
	case ">":
		if value > predOperand {
			predVal = true
		}
	case ">=":
		if value >= predOperand {
			predVal = true
		}
	case "==":
		if value == predOperand {
			predVal = true
		}
	case "!=":
		if value != predOperand {
			predVal = true
		}
	default:
		panic(fmt.Sprintf("error: unrecognized operator `%s`", predicate))
	}
	if !predVal {
		return
	}

	// execute the instruction
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
			It("creates a register when it encounters a reference in the instruction", func() {
				_, ok := rs["a"]
				Expect(ok).To(BeFalse())

				rs.execInstruction("a inc 1 if b < 5")

				_, ok = rs["a"]
				Expect(ok).To(BeTrue())
			})

			It("increments a register by a positive integer", func() {
				rs.execInstruction("a inc 11 if b < 5")
				Expect(rs["a"]).To(Equal(11))
			})

			It("increments a register by a negative integer", func() {
				rs.execInstruction("a inc -13 if b < 5")
				Expect(rs["a"]).To(Equal(-13))
			})

			It("decrements a register by a positive integer", func() {
				rs.execInstruction("a dec 15 if b < 5")
				Expect(rs["a"]).To(Equal(-15))
			})

			It("decrements a register by a negative integer", func() {
				rs.execInstruction("a dec -17 if b < 5")
				Expect(rs["a"]).To(Equal(17))
			})

			It("creates a register when it encounters a reference in the predicate", func() {
				_, ok := rs["b"]
				Expect(ok).To(BeFalse())

				rs.execInstruction("a inc 1 if b < 5")

				_, ok = rs["b"]
				Expect(ok).To(BeTrue())
			})

			Describe("ignoring an instruction based on the conditional", func() {
				It("<", func() {
					rs.execInstruction("execute inc 1 if c < 5") // true
					rs.execInstruction("ignore inc 1 if d < -5") // false

					Expect(rs["execute"]).To(Equal(1))
					Expect(rs["ignore"]).To(Equal(0))
				})

				It("<=", func() {
					rs.execInstruction("execute1 inc 1 if c <= 0") // true
					rs.execInstruction("execute2 inc 1 if d <= 5") // still true
					rs.execInstruction("ignore inc 1 if e <= -5")  // false

					Expect(rs["execute1"]).To(Equal(1))
					Expect(rs["execute2"]).To(Equal(1))
					Expect(rs["ignore"]).To(Equal(0))
				})

				It(">", func() {
					rs.execInstruction("execute inc 1 if c > -5") // true
					rs.execInstruction("ignore inc 1 if d > 5")   // false

					Expect(rs["execute"]).To(Equal(1))
					Expect(rs["ignore"]).To(Equal(0))
				})

				It(">=", func() {
					rs.execInstruction("execute1 inc 1 if c >= 0")  // true
					rs.execInstruction("execute2 inc 1 if d >= -5") // still true
					rs.execInstruction("ignore inc 1 if e >= 5")    // false

					Expect(rs["execute1"]).To(Equal(1))
					Expect(rs["execute2"]).To(Equal(1))
					Expect(rs["ignore"]).To(Equal(0))
				})

				It("==", func() {
					rs.execInstruction("execute inc 1 if c == 0")  // true
					rs.execInstruction("ignore1 inc 1 if d == -5") // false
					rs.execInstruction("ignore2 inc 1 if e == 5")  // false

					Expect(rs["execute"]).To(Equal(1))
					Expect(rs["ignore1"]).To(Equal(0))
					Expect(rs["ignore2"]).To(Equal(0))
				})

				It("!=", func() {
					rs.execInstruction("execute1 inc 1 if c != 5")  // true
					rs.execInstruction("execute2 inc 1 if d != -5") // true
					rs.execInstruction("ignore inc 1 if e != 0")    // false

					Expect(rs["execute1"]).To(Equal(1))
					Expect(rs["execute2"]).To(Equal(1))
					Expect(rs["ignore"]).To(Equal(0))
				})
			})

			It("should handle stateful behavior correctly", func() {
				rs.execInstruction("a inc 1 if b > 0") // false
				Expect(rs["a"]).To(Equal(0))

				rs.execInstruction("b inc 1 if c == 0") // true
				rs.execInstruction("a inc 1 if b > 0")  // now true
				Expect(rs["a"]).To(Equal(1))
			})
		})
	})

	Describe("puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day8.txt")
		instructions := strings.Split(string(raw_data), "\n")

		It("answers star 1", func() {
			rs := NewRegisterSet()
			for _, instruction := range instructions {
				if len(instruction) > 0 {
					rs.execInstruction(instruction)
				}
			}
			var max int
			for _, value := range rs {
				if value > max {
					max = value
				}
			}
			fmt.Printf("d8 s1: largest value is %d\n", max)
		})
	})
})
