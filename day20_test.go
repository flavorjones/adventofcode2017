package adventofcode2017_test

import (
	"math"
	"regexp"
	"strconv"

	"github.com/kr/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Cartesian3Coordinates struct {
	x int
	y int
	z int
}

func (c Cartesian3Coordinates) magnitude() float64 {
	return math.Sqrt(math.Pow(float64(c.x), 2.0) + math.Pow(float64(c.y), 2.0) + math.Pow(float64(c.z), 2.0))
}

type ParticleState struct {
	position         Cartesian3Coordinates
	velocity         Cartesian3Coordinates
	acceleration     Cartesian3Coordinates
	previousPosition Cartesian3Coordinates
}

type ParticleSet struct {
	particles []ParticleState
}

func NewParticleSet() *ParticleSet {
	return &ParticleSet{}
}

var positionRe = regexp.MustCompile(`.*p=< ?(-?\w+), ?(-?\w+), ?(-?\w+)>`)
var velocityRe = regexp.MustCompile(`.*v=< ?(-?\w+), ?(-?\w+), ?(-?\w+)>`)
var accelerationRe = regexp.MustCompile(`.*a=< ?(-?\w+), ?(-?\w+), ?(-?\w+)>`)

func (p *ParticleSet) addParticle(pdesc string) {
	coordinatesFor := func(re *regexp.Regexp) Cartesian3Coordinates {
		match := re.FindStringSubmatch(pdesc)
		if match == nil {
			panic(pretty.Sprintf("error: could not parse `%s` with %v", pdesc, re))
		}
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		z, _ := strconv.Atoi(match[3])
		return Cartesian3Coordinates{x, y, z}
	}

	particle := ParticleState{
		position:     coordinatesFor(positionRe),
		velocity:     coordinatesFor(velocityRe),
		acceleration: coordinatesFor(accelerationRe),
	}
	p.particles = append(p.particles, particle)
}

var _ = Describe("Day20", func() {
	Describe("Cartesian3CoordinatesF", func() {
		It("calculates the magnitude", func() {
			c := Cartesian3Coordinates{3, 3, 3}
			Expect(c.magnitude()).To(BeNumerically("~", 5.196, 0.001))
		})
	})

	Describe("ParticleSet", func() {
		var p *ParticleSet

		BeforeEach(func() {
			p = NewParticleSet()
			p.addParticle("p=< 3,0,0>, v=< 2,0,0>, a=<-1,0,0>")
			p.addParticle("p=< 4,0,0>, v=< 0,0,0>, a=<-2,0,0>")
		})

		It("parses position", func() {
			Expect(p.particles[0].position).To(Equal(Cartesian3Coordinates{3, 0, 0}))
			Expect(p.particles[1].position).To(Equal(Cartesian3Coordinates{4, 0, 0}))
		})

		It("parses velocity", func() {
			Expect(p.particles[0].velocity).To(Equal(Cartesian3Coordinates{2, 0, 0}))
			Expect(p.particles[1].velocity).To(Equal(Cartesian3Coordinates{0, 0, 0}))
		})

		It("parses acceleration", func() {
			Expect(p.particles[0].acceleration).To(Equal(Cartesian3Coordinates{-1, 0, 0}))
			Expect(p.particles[1].acceleration).To(Equal(Cartesian3Coordinates{-2, 0, 0}))
		})
	})

	// Describe("puzzle", func() {
	// 	rawData, _ := ioutil.ReadFile("day20.txt")

	// 	It("solves star 1", func() {
	// 		p := NewParticleSet()
	// 		for _, pdesc := range strings.Split(string(rawData), "\n") {
	// 			if len(pdesc) == 0 {
	// 				continue
	// 			}
	// 			p.addParticle(pdesc)
	// 		}

	// 		jmin := -1
	// 		min := math.MaxFloat64
	// 		for j, acceleration := range p.acceleration {
	// 			current := acceleration.magnitude()
	// 			if current < min {
	// 				min = current
	// 				jmin = j
	// 			}
	// 		}
	// 		pretty.Printf("d20 s1: closest particle will be %d %v\n", jmin, p.acceleration[jmin])
	// 	})
	// })
})
