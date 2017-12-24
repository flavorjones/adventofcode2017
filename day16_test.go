package adventofcode2017_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ProgramDance struct {
	programs []byte
}

func NewProgramDance(size int) *ProgramDance {
	if size > 26 {
		panic("can't create a dance that big")
	}

	programs := make([]byte, size)
	for j := byte(0); j < byte(size); j++ {
		programs[j] = 'a' + j
	}

	return &ProgramDance{programs: programs}
}

var stepSpinRe = regexp.MustCompile(`s(\d+)`)
var stepExchangeRe = regexp.MustCompile(`x(\d+)/(\d+)`)
var stepPartnerRe = regexp.MustCompile(`p(\w+)/(\w+)`)

func (p *ProgramDance) step(step string) {
	switch {
	case stepSpinRe.MatchString(step):
		matches := stepSpinRe.FindStringSubmatch(step)
		spin, _ := strconv.Atoi(matches[1])

		if spin == 0 {
			break
		}

		np := make([]byte, len(p.programs))
		copy(np, p.programs[len(p.programs)-spin:])
		copy(np[spin:], p.programs)
		p.programs = np

	case stepExchangeRe.MatchString(step):
		matches := stepExchangeRe.FindStringSubmatch(step)
		a, _ := strconv.Atoi(matches[1])
		b, _ := strconv.Atoi(matches[2])
		p.programs[a], p.programs[b] = p.programs[b], p.programs[a]

	case stepPartnerRe.MatchString(step):
		matches := stepPartnerRe.FindStringSubmatch(step)
		a := bytes.IndexByte(p.programs, matches[1][0])
		b := bytes.IndexByte(p.programs, matches[2][0])
		p.programs[a], p.programs[b] = p.programs[b], p.programs[a]

	default:
		panic(fmt.Sprintf("error: could not parse step `%s`", step))
	}
}

func (p *ProgramDance) dance(dance string) {
	p.danceN(dance, 1)
}

func (p *ProgramDance) danceN(dance string, repeat int) {
	var nonPartnerSteps []string
	var partnerSteps []string

	for _, step := range strings.Split(dance, ",") {
		if stepPartnerRe.MatchString(step) {
			partnerSteps = append(partnerSteps, step)
		} else {
			nonPartnerSteps = append(nonPartnerSteps, step)
		}
	}

	p.danceN_nonpartner(nonPartnerSteps, repeat)

	p.danceN_partner(partnerSteps, repeat)
}

func (p *ProgramDance) danceN_nonpartner(steps []string, repeat int) {
	//
	//  optimization: do the dance once, and track where programs ended
	//  up. save those position translations in `moveTo` and replay it
	//
	moveTo := make([]int, len(p.programs))
	save := make([]byte, len(p.programs))
	copy(save, p.programs)

	for _, step := range steps {
		p.step(step)
	}

	for jprogram, program := range save {
		moveTo[jprogram] = bytes.IndexByte(p.programs, program)
	}

	for j := 0; j < repeat-1; j++ {
		swap := make([]byte, len(p.programs))
		for jprogram, program := range p.programs {
			swap[moveTo[jprogram]] = program
		}
		p.programs = swap
	}
}

func (p *ProgramDance) danceN_partner(steps []string, repeat int) {
	//
	//  move programs into a hash. compile steps to avoid unnecessary
	//  regexping. use a cache to short-circuit when possible.
	//
	programs := make([]int, 200) // program → position
	for j, program := range p.programs {
		programs[program] = j
	}

	// func to undo the hashification above
	programsToByteSlice := func() []byte {
		np := make([]byte, len(p.programs))
		for _, program := range p.programs {
			np[programs[program]] = program
		}
		return np
	}

	var compiledSteps [][2]byte
	for _, step := range steps {
		matches := stepPartnerRe.FindStringSubmatch(step)
		a := matches[1][0]
		b := matches[2][0]
		compiledSteps = append(compiledSteps, [2]byte{a, b})
	}

	cache := make(map[string]string)
	for j := 0; j < repeat; j++ {
		cachekey := string(programsToByteSlice())

		danceResults, ok := cache[cachekey]
		if ok {
			for j, char := range []byte(danceResults) {
				programs[char] = j
			}
		} else {
			for _, chars := range compiledSteps {
				programs[chars[0]], programs[chars[1]] = programs[chars[1]], programs[chars[0]]
			}
			cache[cachekey] = string(programsToByteSlice())
		}
	}

	p.programs = programsToByteSlice()
}

var _ = Describe("Day16", func() {
	rawData, _ := ioutil.ReadFile("day16.txt")
	danceMoves := string(rawData)

	Describe("ProgramDance", func() {
		Describe("NewProgramDance", func() {
			It("takes a number of programs and sets up the dance", func() {
				p := NewProgramDance(5)
				Expect(p.programs).To(Equal([]byte("abcde")))

				p = NewProgramDance(16)
				Expect(p.programs).To(Equal([]byte("abcdefghijklmnop")))
			})
		})

		Describe("step()", func() {
			var p *ProgramDance

			BeforeEach(func() {
				p = NewProgramDance(5)
			})

			It("the spin", func() {
				p.step("s3")
				Expect(p.programs).To(Equal([]byte("cdeab")))
			})

			It("the exchange", func() {
				p.step("x3/4")
				Expect(p.programs).To(Equal([]byte("abced")))
			})

			It("the partner", func() {
				p.step("pe/b")
				Expect(p.programs).To(Equal([]byte("aecdb")))
			})

			It("puts it all together", func() {
				p.step("s1")
				Expect(p.programs).To(Equal([]byte("eabcd")))
				p.step("x3/4")
				Expect(p.programs).To(Equal([]byte("eabdc")))
				p.step("pe/b")
				Expect(p.programs).To(Equal([]byte("baedc")))
			})
		})

		Describe("dance()", func() {
			It("performs the steps", func() {
				p := NewProgramDance(5)
				p.dance("s1,x3/4,pe/b")
				Expect(p.programs).To(Equal([]byte("baedc")))
			})

			It("performs the steps", func() {
				p := NewProgramDance(16)
				p.dance(danceMoves)
				Expect(p.programs).To(Equal([]byte("kpbodeajhlicngmf")))
			})
		})

		Describe("danceN()", func() {
			It("performs the steps multiple times", func() {
				p := NewProgramDance(5)
				p.danceN("s1,x3/4,pe/b", 2)
				Expect(p.programs).To(Equal([]byte("ceadb")))
			})
		})

		Describe("equivalence", func() {
			It("dance x N and danceN should be equivalent", func() {
				p := NewProgramDance(16)
				p.dance(danceMoves)
				p.dance(danceMoves)
				Expect(p.programs).To(Equal([]byte("dkfcagielbnjohpm")))

				p = NewProgramDance(16)
				p.danceN(danceMoves, 2)
				Expect(p.programs).To(Equal([]byte("dkfcagielbnjohpm")))

				p = NewProgramDance(16)
				p.dance(danceMoves)
				p.dance(danceMoves)
				p.dance(danceMoves)
				Expect(p.programs).To(Equal([]byte("bhdakljmfocgpeni")))

				p = NewProgramDance(16)
				p.danceN(danceMoves, 3)
				Expect(p.programs).To(Equal([]byte("bhdakljmfocgpeni")))

				p = NewProgramDance(16)
				for j := 0; j < 10; j++ {
					p.dance(danceMoves)
				}
				Expect(p.programs).To(Equal([]byte("jhgpadkcbfmnolei")))

				p = NewProgramDance(16)
				p.danceN(danceMoves, 10)
				Expect(p.programs).To(Equal([]byte("jhgpadkcbfmnolei")))
			})
		})
	})

	Describe("puzzle", func() {
		It("solves star 1", func() {
			p := NewProgramDance(16)
			p.dance(danceMoves)
			fmt.Printf("d16 s1: final order is %s\n", p.programs)
		})

		It("solves star 2", func() {
			p := NewProgramDance(16)
			p.danceN(danceMoves, 1000000000)
			fmt.Printf("d16 s2: final order is %s\n", p.programs)
		})
	})
})
