package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var virusUp = CartesianCoordinates{0, 1}
var virusRight = CartesianCoordinates{1, 0}
var virusDown = CartesianCoordinates{0, -1}
var virusLeft = CartesianCoordinates{-1, 0}

var virusTurnLeft = map[CartesianCoordinates]CartesianCoordinates{
	virusUp:    virusLeft,
	virusLeft:  virusDown,
	virusDown:  virusRight,
	virusRight: virusUp,
}

var virusTurnRight = map[CartesianCoordinates]CartesianCoordinates{
	virusUp:    virusRight,
	virusRight: virusDown,
	virusDown:  virusLeft,
	virusLeft:  virusUp,
}

type SporificaVirus struct {
	infected   map[CartesianCoordinates]bool
	position   CartesianCoordinates
	direction  CartesianCoordinates
	infections int
}

func NewSporificaVirus(nodeMap string) *SporificaVirus {
	sv := SporificaVirus{direction: virusUp, infected: make(map[CartesianCoordinates]bool)}

	nodeMapLines := strings.Split(nodeMap, "\n")
	size := len(nodeMapLines[0])
	offset := (size - 1) / 2
	for jrow, line := range nodeMapLines {
		for jcol, char := range []byte(line) {
			if char == '#' {
				coords := CartesianCoordinates{jcol - offset, offset - jrow}
				sv.infected[coords] = true
			}
		}
	}

	return &sv
}

func (sv *SporificaVirus) NodeInfected(coords CartesianCoordinates) bool {
	result, ok := sv.infected[coords]
	if !ok {
		return false
	}
	return result
}

func (sv *SporificaVirus) Burst() {
	if sv.NodeInfected(sv.position) {
		sv.infected[sv.position] = false
		sv.direction = virusTurnRight[sv.direction]
	} else {
		sv.infections++
		sv.infected[sv.position] = true
		sv.direction = virusTurnLeft[sv.direction]
	}
	sv.position = sv.position.move(sv.direction)
}

var _ = Describe("Day22", func() {
	Describe("SporificaVirus", func() {
		var testMap = heredoc.Doc(`
			..#
			#..
			...
		`)

		Describe("NewSporificaVirus", func() {
			It("points up", func() {
				sv := NewSporificaVirus(testMap)
				Expect(sv.direction).To(Equal(virusUp))
			})

			It("positions itself in the middle of the map", func() {
				sv := NewSporificaVirus(testMap)
				Expect(sv.position).To(Equal(CartesianCoordinates{0, 0}))
			})

			It("stores infected node positions", func() {
				sv := NewSporificaVirus(testMap)
				Expect(sv.NodeInfected(CartesianCoordinates{0, 0})).To(BeFalse())
				Expect(sv.NodeInfected(CartesianCoordinates{1, 1})).To(BeTrue())
				Expect(sv.NodeInfected(CartesianCoordinates{-1, 0})).To(BeTrue())
			})
		})

		Describe("Burst", func() {
			Context("when on a clean node", func() {
				var sv *SporificaVirus

				BeforeEach(func() {
					sv = NewSporificaVirus(testMap)
					Expect(sv.NodeInfected(sv.position)).To(BeFalse())
				})

				It("infects the node, turns left, moves forward", func() {
					position := sv.position
					sv.Burst()
					Expect(sv.NodeInfected(position)).To(BeTrue())
					Expect(sv.direction).To(Equal(virusLeft))
					Expect(sv.position).To(Equal(CartesianCoordinates{-1, 0}))
					Expect(sv.infections).To(Equal(1))
				})
			})

			Context("when on an infected node", func() {
				var sv *SporificaVirus

				BeforeEach(func() {
					sv = NewSporificaVirus(testMap)
					sv.Burst()
					Expect(sv.NodeInfected(sv.position)).To(BeTrue())
					Expect(sv.direction).To(Equal(virusLeft))
				})

				It("cleans the node, turns right, moves forward", func() {
					position := sv.position
					sv.Burst()
					Expect(sv.NodeInfected(position)).To(BeFalse())
					Expect(sv.direction).To(Equal(virusUp))
					Expect(sv.position).To(Equal(CartesianCoordinates{-1, 1}))
					Expect(sv.infections).To(Equal(1))
				})
			})

			It("does the right thing for an ad-hoc test", func() {
				sv := NewSporificaVirus(testMap)
				for j := 1; j <= 70; j++ {
					sv.Burst()
				}
				Expect(sv.infections).To(Equal(41))
			})

			It("does the right thing for an ad-hoc test", func() {
				sv := NewSporificaVirus(testMap)
				for j := 1; j <= 10000; j++ {
					sv.Burst()
				}
				Expect(sv.infections).To(Equal(5587))
			})
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day22.txt")
		nodeMap := string(rawData)

		It("solves star 1", func() {
			sv := NewSporificaVirus(nodeMap)
			Expect(sv.NodeInfected(CartesianCoordinates{-12, 12})).To(BeTrue())

			for j := 1; j <= 10000; j++ {
				sv.Burst()
			}
			fmt.Printf("d22 s1: there were %d infections\n", sv.infections)
		})
	})
})
