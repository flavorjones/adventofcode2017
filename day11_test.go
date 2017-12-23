package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// 30/60/90 triangle edges are 1, sqrt(3); and hypotenuse is 2
var LONG_LEG = math.Sqrt(3)
var SHORT_LEG = 1.0
var HYPOTENUSE = 2.0

var translations = map[string]CartesianCoordinatesF{ // direction → translation
	"ne": CartesianCoordinatesF{LONG_LEG, SHORT_LEG},
	"se": CartesianCoordinatesF{LONG_LEG, -SHORT_LEG},
	"s":  CartesianCoordinatesF{0, -HYPOTENUSE},
	"sw": CartesianCoordinatesF{-LONG_LEG, -SHORT_LEG},
	"nw": CartesianCoordinatesF{-LONG_LEG, SHORT_LEG},
	"n":  CartesianCoordinatesF{0, HYPOTENUSE},
}

type CartesianCoordinatesF struct {
	x float64
	y float64
}

func (c CartesianCoordinatesF) move(relative CartesianCoordinatesF) CartesianCoordinatesF {
	return CartesianCoordinatesF{c.x + relative.x, c.y + relative.y}
}

func (c CartesianCoordinatesF) distanceToOrigin() float64 {
	return math.Sqrt(math.Pow(c.x, 2.0) + math.Pow(c.y, 2.0))
}

type Hextile struct {
	position CartesianCoordinatesF
}

func NewHextile() *Hextile {
	return &Hextile{}
}

func (h *Hextile) move(direction string) {
	translation, ok := translations[direction]
	if !ok {
		panic(fmt.Sprintf("error: could not find direction `%s`\n", direction))
	}

	h.position = h.position.move(translation)
}

func (h *Hextile) moveMany(directionsStr string) {
	for _, direction := range strings.Split(directionsStr, ",") {
		direction = strings.TrimSpace(direction)
		h.move(direction)
	}
}

func (h *Hextile) stepsAway() int {
	steps := 0
	pos := h.position

	for pos.distanceToOrigin() > 0.1 { // until we get to the origin
		// look at all six possible moves
		best_try := pos
		for _, translation := range translations {
			// and pick the one that moves us closest to the origin
			try := pos.move(translation)
			if try.distanceToOrigin() < best_try.distanceToOrigin() {
				best_try = try
			}
		}
		pos = best_try
		steps++
	}
	return steps
}

var _ = Describe("Day11", func() {
	Describe("Hextile", func() {
		Describe("move()", func() {
			It("moves ne", func() {
				h := NewHextile()
				h.move("ne")
				Expect(h.position.x).To(BeNumerically("~", LONG_LEG))
				Expect(h.position.y).To(BeNumerically("~", SHORT_LEG))
				Expect(h.position.distanceToOrigin()).To(BeNumerically("~", HYPOTENUSE))
			})
		})

		Describe("moveMany()", func() {
			It("takes multiple steps", func() {
				h := NewHextile()
				h.moveMany("ne,ne,ne")
				Expect(h.position.x).To(BeNumerically("~", 3*LONG_LEG))
				Expect(h.position.y).To(BeNumerically("~", 3*SHORT_LEG))
			})

			It("takes multiple steps", func() {
				h := NewHextile()
				h.moveMany("ne,ne,sw,sw")
				Expect(h.position.x).To(BeNumerically("~", 0))
				Expect(h.position.y).To(BeNumerically("~", 0))
			})
		})

		Describe("stepsAway()", func() {
			It("returns the minimal number of steps back to origin", func() {
				h := NewHextile()
				h.moveMany("ne,ne,ne")
				Expect(h.stepsAway()).To(Equal(3))
			})

			It("returns the minimal number of steps back to origin", func() {
				h := NewHextile()
				h.moveMany("ne,ne,sw,sw")
				Expect(h.stepsAway()).To(Equal(0))
			})

			It("returns the minimal number of steps back to origin", func() {
				h := NewHextile()
				h.moveMany("ne,ne,s,s")
				Expect(h.stepsAway()).To(Equal(2))
			})

			It("returns the minimal number of steps back to origin", func() {
				h := NewHextile()
				h.moveMany("se,sw,se,sw,sw")
				Expect(h.stepsAway()).To(Equal(3))
			})
		})
	})

	Describe("puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day11.txt")
		walk := string(raw_data)

		It("solves star 1", func() {
			h := NewHextile()
			h.moveMany(walk)
			fmt.Printf("d11 s1: child is %d steps away\n", h.stepsAway())
		})
	})
})
