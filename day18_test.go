package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DuetCpu struct {
	id        int
	registers map[byte]int
	pc        int
	incoming  chan int
	outgoing  chan int
	sentCount int
	mulCount  int
}

func NewDuetCpu(id int) *DuetCpu {
	d := DuetCpu{registers: make(map[byte]int), incoming: make(chan int, 100)}
	d.id = id
	d.registers['p'] = d.id
	return &d
}

func (s *DuetCpu) setOutgoing(outgoing chan int) {
	s.outgoing = outgoing
}

func (s *DuetCpu) getRegister(name byte) int {
	val, ok := s.registers[name]
	if ok {
		return val
	}
	s.registers[name] = 0
	return 0
}

func (s *DuetCpu) valueOf(thing string) int {
	if val, err := strconv.Atoi(thing); err == nil {
		return val
	} else {
		return s.getRegister(thing[0])
	}
}

var oneArgDuetCpuInstructionRe = regexp.MustCompile(`(snd|rcv) (-?\w+)`)
var twoArgDuetCpuInstructionRe = regexp.MustCompile(`(set|add|sub|mul|mod|jgz|jnz) (-?\w+) (-?\w+)`)

func (s *DuetCpu) execInstruction(instruction string) {
	switch {
	case oneArgDuetCpuInstructionRe.MatchString(instruction):
		match := oneArgDuetCpuInstructionRe.FindStringSubmatch(instruction)

		switch match[1] {
		case "snd":
			srcValue := s.valueOf(match[2])
			s.outgoing <- srcValue
			s.sentCount++
			s.pc++
		case "rcv":
			tgtName := match[2][0]
			select {
			case s.registers[tgtName] = <-s.incoming:
				s.pc++
			case <-time.After(time.Second):
				s.pc = math.MaxInt32 // should terminate program
			}
		default:
			panic(fmt.Sprintf("error: could not execute instruction `%s`", match[1]))
		}

	case twoArgDuetCpuInstructionRe.MatchString(instruction):
		match := twoArgDuetCpuInstructionRe.FindStringSubmatch(instruction)
		tgtName := match[2][0]
		srcValue := s.valueOf(match[3])

		switch match[1] {
		case "set":
			s.registers[tgtName] = srcValue
			s.pc++
		case "add":
			s.registers[tgtName] = s.getRegister(tgtName) + srcValue
			s.pc++
		case "sub":
			s.registers[tgtName] = s.getRegister(tgtName) - srcValue
			s.pc++
		case "mul":
			s.registers[tgtName] = s.getRegister(tgtName) * srcValue
			s.pc++
			s.mulCount++
		case "mod":
			s.registers[tgtName] = s.getRegister(tgtName) % srcValue
			s.pc++
		case "jgz":
			if s.valueOf(match[2]) > 0 {
				s.pc += srcValue
			} else {
				s.pc++
			}
		case "jnz":
			if s.valueOf(match[2]) != 0 {
				s.pc += srcValue
			} else {
				s.pc++
			}
		default:
			panic(fmt.Sprintf("error: could not execute instruction `%s`", match[1]))
		}
	}
}

func (s *DuetCpu) execInstructions(rawInstructions string) {
	instructions := strings.Split(rawInstructions, "\n")
	instructionLen := len(instructions)
	if len(instructions[instructionLen-1]) == 0 {
		instructionLen -= 1
	}
	for s.pc < instructionLen {
		s.execInstruction(instructions[s.pc])
	}
}

var _ = Describe("Day18", func() {
	Describe("DuetCpu", func() {
		var s *DuetCpu
		var c chan int

		BeforeEach(func() {
			s = NewDuetCpu(0)
			c = make(chan int, 100)
			s.setOutgoing(c)
		})

		Describe("execInstruction", func() {
			It("should increment the program counter", func() {
				Expect(s.pc).To(Equal(0))
				s.execInstruction("set a 11")
				Expect(s.pc).To(Equal(1))
				s.execInstruction("add a 12")
				Expect(s.pc).To(Equal(2))
				s.execInstruction("mul a 21")
				Expect(s.pc).To(Equal(3))
				s.execInstruction("mod a 31")
				Expect(s.pc).To(Equal(4))
				s.execInstruction("snd a")
				Expect(s.pc).To(Equal(5))
				_ = <-c
				s.incoming <- 3
				s.execInstruction("rcv a")
				Expect(s.pc).To(Equal(6))
			})

			Describe("set", func() {
				It("should save a value to a register", func() {
					s.execInstruction("set a 1")
					Expect(s.getRegister('a')).To(Equal(1))
				})

				It("should copy a register to a register", func() {
					s.execInstruction("set b -2")
					s.execInstruction("set a b")
					Expect(s.getRegister('a')).To(Equal(-2))
				})
			})

			Describe("add", func() {
				It("should add a value to a register", func() {
					s.execInstruction("set a 3")
					s.execInstruction("add a 5")
					Expect(s.getRegister('a')).To(Equal(8))
				})

				It("should add a register value to a register", func() {
					s.execInstruction("set a 4")
					s.execInstruction("set b 9")
					s.execInstruction("add a b")
					Expect(s.getRegister('a')).To(Equal(13))
				})
			})

			Describe("sub", func() {
				It("should subtract a value from a register", func() {
					s.execInstruction("set a 3")
					s.execInstruction("sub a -5")
					Expect(s.getRegister('a')).To(Equal(8))
				})

				It("should subtract a register value from a register", func() {
					s.execInstruction("set a 4")
					s.execInstruction("set b 9")
					s.execInstruction("sub a b")
					Expect(s.getRegister('a')).To(Equal(-5))
				})
			})

			Describe("mul", func() {
				It("should multiply a value with a register", func() {
					s.execInstruction("set a 5")
					s.execInstruction("mul a 5")
					Expect(s.getRegister('a')).To(Equal(25))
				})

				It("should multiply a register value with a register", func() {
					s.execInstruction("set a 6")
					s.execInstruction("set b 7")
					s.execInstruction("mul a b")
					Expect(s.getRegister('a')).To(Equal(42))
				})
			})

			Describe("mod", func() {
				It("should modulo a value with a register", func() {
					s.execInstruction("set a 33")
					s.execInstruction("mod a 6")
					Expect(s.getRegister('a')).To(Equal(3))
				})

				It("should modulo a register value with a register", func() {
					s.execInstruction("set a 32")
					s.execInstruction("set b 5")
					s.execInstruction("mod a b")
					Expect(s.getRegister('a')).To(Equal(2))
				})
			})

			Describe("snd", func() {
				It("sends a literal value on a channel", func() {
					s.execInstruction("snd 11")
					result := <-c
					Expect(result).To(Equal(11))
				})

				It("sends a register value on a channel", func() {
					s.execInstruction("set a 12")
					s.execInstruction("snd a")
					result := <-c
					Expect(result).To(Equal(12))
				})

				It("increments a sent counter", func() {
					s.execInstruction("set a 1")
					s.execInstruction("snd 11")
					_ = <-c
					s.execInstruction("snd 12")
					_ = <-c
					Expect(s.sentCount).To(Equal(2))
				})
			})

			Describe("rcv", func() {
				It("write the received value to a register", func() {
					s.incoming <- 33
					s.execInstruction("rcv a")
					Expect(s.getRegister('a')).To(Equal(33))
				})

				It("times out", func() {
					s.execInstruction("rcv b")
				})
			})

			Describe("jgz", func() {
				It("sets program counter forward based on literal", func() {
					s.execInstruction("jgz 1 10")
					Expect(s.pc).To(Equal(10))
				})

				It("sets program counter backwards based on register", func() {
					s.execInstruction("set a 1")
					s.execInstruction("jgz a -10")
					Expect(s.pc).To(Equal(-9))
				})

				It("increments pc if literal arg is zero", func() {
					s.execInstruction("jgz 0 10")
					Expect(s.pc).To(Equal(1))
				})

				It("increments pc if register arg is less than zero", func() {
					s.execInstruction("set a -1")
					s.execInstruction("jgz a 10")
					Expect(s.pc).To(Equal(2))
				})
			})

			Describe("jnz", func() {
				It("sets program counter forward based on literal", func() {
					s.execInstruction("jnz 1 10")
					Expect(s.pc).To(Equal(10))
				})

				It("sets program counter backwards based on register", func() {
					s.execInstruction("set a 1")
					s.execInstruction("jnz a -10")
					Expect(s.pc).To(Equal(-9))
				})

				It("increments pc if literal arg is zero", func() {
					s.execInstruction("jnz 0 10")
					Expect(s.pc).To(Equal(1))
				})

				It("sets program counter if register arg is less than zero", func() {
					s.execInstruction("set a -1")
					s.execInstruction("jnz a 10")
					Expect(s.pc).To(Equal(11))
				})
			})
		})

		Describe("execInstructions", func() {
			// It("runs a bunch of instructions, paying attention to pc", func() {
			// 	instructions := heredoc.Doc(`
			// 		set a 1
			// 		add a 2
			// 		mul a a
			// 		mod a 5
			// 		snd a
			// 		set a 0
			// 		rcv a
			// 		jgz a -1
			// 		set a 1
			// 		jgz a -2
			// 	`)
			// 	s.execInstructions(instructions)
			//  Expect(s.recovered).To(Equal(4))
			// })
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day18.txt")
		instructions := string(rawData)

		// It("solves star 1", func() {
		// 	s := NewDuetCpu()
		// 	s.execInstructions(string(rawData))
		// 	answer := s.recovered
		// 	fmt.Printf("d18 s1: recovered %d\n", answer)
		// })

		It("solves star 2", func() {
			s0 := NewDuetCpu(0)
			s1 := NewDuetCpu(1)
			s0.setOutgoing(s1.incoming)
			s1.setOutgoing(s0.incoming)

			go s0.execInstructions(instructions)
			s1.execInstructions(instructions)
			fmt.Printf("d18 s2: cpu 1 sent a value %d times\n", s1.sentCount)
		})
	})
})

var _ = Describe("Day23", func() {
	It("counts how many times `mul` is invoked", func() {
		instructions := heredoc.Doc(`
			set a 1
			add a 2
			mul a a
			mod a 5
			set a 0
			set a 1
		`)
		s := NewDuetCpu(0)
		s.execInstructions(instructions)
		Expect(s.mulCount).To(Equal(1))
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day23.txt")
		instructions := string(rawData)

		It("solves star 1", func() {
			s := NewDuetCpu(0)
			s.execInstructions(instructions)
			fmt.Printf("d23 s1: mul was called %d times\n", s.mulCount)
		})
	})
})
