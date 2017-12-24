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

func (p *ProgramDance) steps(steps string) {
	for _, step := range strings.Split(steps, ",") {
		if len(step) == 0 {
			continue
		}
		p.step(step)
	}
}

var _ = Describe("Day16", func() {
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
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day16.txt")

		It("solves star 1", func() {
			p := NewProgramDance(16)
			p.steps(string(rawData))
			fmt.Printf("d16 s1: final order is %s\n", p.programs)
		})
	})
})
