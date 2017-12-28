package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Cartesian3CoordinatesF struct {
	x float64
	y float64
	z float64
}

func (c Cartesian3CoordinatesF) distance() float64 {
	return math.Sqrt(math.Pow(c.x, 2.0) + math.Pow(c.y, 2.0) + math.Pow(c.z, 2.0))
}

type ParticleSet struct {
	acceleration []Cartesian3CoordinatesF
}

func NewParticleSet() *ParticleSet {
	return &ParticleSet{}
}

var accelerationRe = regexp.MustCompile(`.*a=< ?(-?\w+), ?(-?\w+), ?(-?\w+)>`)

func (p *ParticleSet) addParticle(pdesc string) {
	match := accelerationRe.FindStringSubmatch(pdesc)
	if match == nil {
		panic(fmt.Sprintf("error: could not parse `%s`", pdesc))
	}
	x, _ := strconv.ParseFloat(match[1], 64)
	y, _ := strconv.ParseFloat(match[2], 64)
	z, _ := strconv.ParseFloat(match[3], 64)
	c3 := Cartesian3CoordinatesF{x, y, z}
	p.acceleration = append(p.acceleration, c3)
}

var _ = Describe("Day20", func() {
	Describe("Cartesian3CoordinatesF", func() {
		It("calculates the magnitude", func() {
			c := Cartesian3CoordinatesF{3.0, 3.0, 3.0}
			Expect(c.distance()).To(BeNumerically("~", 5.196, 0.001))
		})
	})

	Describe("ParticleSet", func() {
		It("parses acceleration", func() {
			p := NewParticleSet()
			p.addParticle("p=< 3,0,0>, v=< 2,0,0>, a=<-1,0,0>")
			Expect(p.acceleration[0]).To(Equal(Cartesian3CoordinatesF{-1.0, 0.0, 0.0}))

			p.addParticle("p=< 4,0,0>, v=< 0,0,0>, a=<-2,0,0>")
			Expect(p.acceleration[1]).To(Equal(Cartesian3CoordinatesF{-2.0, 0.0, 0.0}))
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day20.txt")

		It("solves star 1", func() {
			p := NewParticleSet()
			for _, pdesc := range strings.Split(string(rawData), "\n") {
				if len(pdesc) == 0 {
					continue
				}
				p.addParticle(pdesc)
			}

			jmin := -1
			min := math.MaxFloat64
			for j, acceleration := range p.acceleration {
				current := acceleration.distance()
				if current < min {
					min = current
					jmin = j
				}
			}
			pretty.Printf("d20 s1: closest particle will be %d %v\n", jmin, p.acceleration[jmin])
		})
	})
})
