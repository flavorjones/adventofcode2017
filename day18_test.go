package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SoundCard struct {
	registers map[byte]int
	pc        int
	sound     int
	recovered int
}

func NewSoundCard() *SoundCard {
	return &SoundCard{registers: make(map[byte]int)}
}

func (s *SoundCard) getRegister(name byte) int {
	val, ok := s.registers[name]
	if ok {
		return val
	}
	s.registers[name] = 0
	return 0
}

func (s *SoundCard) valueOf(thing string) int {
	if val, err := strconv.Atoi(thing); err == nil {
		return val
	} else {
		return s.getRegister(thing[0])
	}
}

var oneArgSoundCardInstructionRe = regexp.MustCompile(`(snd|rcv) (-?\w+)`)
var twoArgSoundCardInstructionRe = regexp.MustCompile(`(set|add|mul|mod|jgz) (-?\w+) (-?\w+)`)

func (s *SoundCard) execInstruction(instruction string) {
	switch {
	case oneArgSoundCardInstructionRe.MatchString(instruction):
		match := oneArgSoundCardInstructionRe.FindStringSubmatch(instruction)
		srcValue := s.valueOf(match[2])

		switch match[1] {
		case "snd":
			s.sound = srcValue
			s.pc++
		case "rcv":
			if srcValue != 0 {
				s.recovered = s.sound
			}
			s.pc++
		default:
			panic(fmt.Sprintf("error: could not execute instruction `%s`", match[1]))
		}

	case twoArgSoundCardInstructionRe.MatchString(instruction):
		match := twoArgSoundCardInstructionRe.FindStringSubmatch(instruction)
		tgtName := match[2][0]
		srcValue := s.valueOf(match[3])

		switch match[1] {
		case "set":
			s.registers[tgtName] = srcValue
			s.pc++
		case "add":
			s.registers[tgtName] = s.getRegister(tgtName) + srcValue
			s.pc++
		case "mul":
			s.registers[tgtName] = s.getRegister(tgtName) * srcValue
			s.pc++
		case "mod":
			s.registers[tgtName] = s.getRegister(tgtName) % srcValue
			s.pc++
		case "jgz":
			if s.valueOf(match[2]) > 0 {
				s.pc += srcValue
			} else {
				s.pc++
			}
		default:
			panic(fmt.Sprintf("error: could not execute instruction `%s`", match[1]))
		}
	}
}

func (s *SoundCard) execInstructions(rawInstructions string) {
	instructions := strings.Split(rawInstructions, "\n")
	for s.pc < len(instructions) {
		s.execInstruction(instructions[s.pc])
		if s.recovered != 0 {
			break
		}
	}
}

var _ = Describe("Day18", func() {
	Describe("SoundCard", func() {
		var s *SoundCard

		BeforeEach(func() {
			s = NewSoundCard()
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

			Describe("snd/rcv", func() {
				It("should save the sound played when a literal, and recover it", func() {
					s.execInstruction("snd 440")
					s.execInstruction("rcv 1")
					Expect(s.recovered).To(Equal(440))
				})

				It("should save the sound played from a register", func() {
					s.execInstruction("set a 435")
					s.execInstruction("snd a")
					s.execInstruction("rcv -1")
					Expect(s.recovered).To(Equal(435))
				})

				It("should recover based on a register", func() {
					s.execInstruction("set a 435")
					s.execInstruction("set b 2")
					s.execInstruction("snd a")
					s.execInstruction("rcv b")
					Expect(s.recovered).To(Equal(435))
				})

				It("should not recover when rcv arg is zero literal", func() {
					s.execInstruction("snd 440")
					s.execInstruction("rcv 0")
					Expect(s.recovered).To(Equal(0))
				})

				It("should not recover when rcv arg is zero from register", func() {
					s.execInstruction("snd 440")
					s.execInstruction("rcv b")
					Expect(s.recovered).To(Equal(0))
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
		})

		Describe("execInstructions", func() {
			It("runs a bunch of instructions, paying attention to pc", func() {
				instructions := heredoc.Doc(`
					set a 1
					add a 2
					mul a a
					mod a 5
					snd a
					set a 0
					rcv a
					jgz a -1
					set a 1
					jgz a -2
				`)
				s.execInstructions(instructions)
				Expect(s.recovered).To(Equal(4))
			})
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day18.txt")

		It("solves star 1", func() {
			s := NewSoundCard()
			s.execInstructions(string(rawData))
			answer := s.recovered
			fmt.Printf("d18 s1: recovered %d\n", answer)
		})
	})
})
