package adventofcode2017_test

import (
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/kr/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Cartesian3Coordinates struct {
	x int
	y int
	z int
}

func (c Cartesian3Coordinates) distanceTo(o Cartesian3Coordinates) float64 {
	return Cartesian3Coordinates{c.x - o.x, c.y - o.y, c.z - o.z}.magnitude()
}

func (c Cartesian3Coordinates) magnitude() float64 {
	return math.Sqrt(math.Pow(float64(c.x), 2.0) + math.Pow(float64(c.y), 2.0) + math.Pow(float64(c.z), 2.0))
}

func (c Cartesian3Coordinates) move(relative Cartesian3Coordinates) Cartesian3Coordinates {
	return Cartesian3Coordinates{c.x + relative.x, c.y + relative.y, c.z + relative.z}
}

type ParticleState struct {
	position         Cartesian3Coordinates
	velocity         Cartesian3Coordinates
	acceleration     Cartesian3Coordinates
	previousPosition Cartesian3Coordinates
	collided         bool
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

func (p *ParticleSet) addParticles(pdesc string) {
	for _, line := range strings.Split(pdesc, "\n") {
		if len(line) == 0 {
			continue
		}
		p.addParticle(line)
	}
}

func (p *ParticleSet) addParticle(pdesc string) {
	if len(pdesc) == 0 {
		return
	}

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
		collided:     false,
	}
	p.particles = append(p.particles, particle)
}

func (p *ParticleSet) tick(collisionDetection bool) {
	for j, _ := range p.particles {
		if p.particles[j].collided {
			continue
		}
		p.particles[j].previousPosition = p.particles[j].position
		p.particles[j].velocity = p.particles[j].velocity.move(p.particles[j].acceleration)
		p.particles[j].position = p.particles[j].position.move(p.particles[j].velocity)
	}

	if collisionDetection {
		for j := 0; j < len(p.particles); j++ {
			if p.particles[j].collided {
				continue
			}
			for k := j + 1; k < len(p.particles); k++ {
				if p.particles[k].collided {
					continue
				}
				if p.particles[j].position == p.particles[k].position {
					p.particles[j].collided = true
					p.particles[k].collided = true
				}
			}
		}
	}
}

func (p *ParticleSet) tickToSteadyState(collisionDetection bool) {
	for {
		for j := 0; j < 100; j++ {
			p.tick(collisionDetection)
		}

		allReceding := true
	search:
		for j := 0; j < len(p.particles); j++ {
			if p.particles[j].collided {
				continue
			}
			for k := j + 1; k < len(p.particles); k++ {
				if p.particles[k].collided {
					continue
				}
				relativeVelocity := p.particles[j].position.distanceTo(p.particles[k].position) -
					p.particles[j].previousPosition.distanceTo(p.particles[k].previousPosition)
				if relativeVelocity < 0 {
					allReceding = false
					break search
				}
			}
		}
		if allReceding {
			break
		}
	}
}

func (p *ParticleSet) closestToOrigin() (int, ParticleState) {
	jmin := -1
	min := math.MaxFloat64
	for j, particle := range p.particles {
		current := particle.position.magnitude()
		if current < min {
			min = current
			jmin = j
		}
	}
	return jmin, p.particles[jmin]
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

		Describe("addParticle()", func() {
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

		Describe("tick()", func() {
			It("updates position", func() {
				p.tick(false)
				Expect(p.particles[0].position).To(Equal(Cartesian3Coordinates{4, 0, 0}))
				Expect(p.particles[1].position).To(Equal(Cartesian3Coordinates{2, 0, 0}))
			})

			It("updates previousPosition", func() {
				p.tick(false)
				Expect(p.particles[0].previousPosition).To(Equal(Cartesian3Coordinates{3, 0, 0}))
				Expect(p.particles[1].previousPosition).To(Equal(Cartesian3Coordinates{4, 0, 0}))
			})

			It("updates velocity", func() {
				p.tick(false)
				Expect(p.particles[0].velocity).To(Equal(Cartesian3Coordinates{1, 0, 0}))
				Expect(p.particles[1].velocity).To(Equal(Cartesian3Coordinates{-2, 0, 0}))
			})

			Describe("collision detection", func() {
				BeforeEach(func() {
					p = NewParticleSet()
					p.addParticles(heredoc.Doc(`
						p=<-6,0,0>, v=< 3,0,0>, a=< 0,0,0>    
						p=<-4,0,0>, v=< 2,0,0>, a=< 0,0,0>
						p=<-2,0,0>, v=< 1,0,0>, a=< 0,0,0>
						p=< 3,0,0>, v=<-1,0,0>, a=< 0,0,0>
  				`))
				})

				It("can ignore collisions", func() {
					p.tick(false)
					p.tick(false)
					Expect(p.particles[0].collided).To(BeFalse())
					Expect(p.particles[1].collided).To(BeFalse())
					Expect(p.particles[2].collided).To(BeFalse())
					Expect(p.particles[3].collided).To(BeFalse())
				})

				It("can detect collisions", func() {
					p.tick(true)
					p.tick(true)
					Expect(p.particles[0].collided).To(BeTrue())
					Expect(p.particles[1].collided).To(BeTrue())
					Expect(p.particles[2].collided).To(BeTrue())
					Expect(p.particles[3].collided).To(BeFalse())
				})
			})
		})

		Describe("tickToSteadyState", func() {
			It("iteratively calls tick() until all particles are receding from each other", func() {
				p.tickToSteadyState(false)
			})
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day20.txt")

		It("solves star 1", func() {
			p := NewParticleSet()
			p.addParticles(string(rawData))
			p.tickToSteadyState(false)
			jmin, particle := p.closestToOrigin()
			pretty.Printf("d20 s1: closest particle will be %d %v\n", jmin, particle)
		})

		It("solves star 2", func() {
			p := NewParticleSet()
			p.addParticles(string(rawData))
			p.tickToSteadyState(true)

			count := 0
			for _, particle := range p.particles {
				if !particle.collided {
					count++
				}
			}
			pretty.Printf("d20 s2: there are %d uncollided particles\n", count)
		})
	})
})
