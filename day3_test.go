package adventofcode2017_test

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CartesianCoordinates struct {
	x int
	y int
}

// manhattanDistance returns the manhattan distance of the coordinates
func (c CartesianCoordinates) manhattanDistance() int {
	return int(math.Abs(float64(c.x)) + math.Abs(float64(c.y)))
}

// findNearestOddSquare(number int):
// returns the root of the first odd square that's larger than `number`
// so 1, 3, 5, 7 are all valid return values (we return the square root)
func findNearestOddSquare(number int) int {
	for j := 1; ; j += 2 {
		jsq := j * j
		if jsq >= number {
			return j
		}
	}
}

func (c CartesianCoordinates) move(relative CartesianCoordinates) CartesianCoordinates {
	return CartesianCoordinates{c.x + relative.x, c.y + relative.y}
}

// CartesianCoordinates.location returns the location of the coordinates
func (c CartesianCoordinates) location() int {
	if (c == CartesianCoordinates{0, 0}) {
		return 1
	}

	min_root := (2 * int(math.Max(math.Abs(float64(c.x)), math.Abs(float64(c.y))))) - 1
	min_square := int(math.Pow(float64(min_root), 2.0))
	max_square := int(math.Pow(float64(min_root+2), 2.0))

	// brute force it because I'm lazy
	for j := min_square + 1; j <= max_square; j++ {
		if locationToCoordinates(j) == c {
			return j
		}
	}

	return -1
}

// locationToCoordinates returns the cartesian coordinates of `location`
func locationToCoordinates(location int) CartesianCoordinates {
	coords := CartesianCoordinates{}

	// special case
	if location == 1 {
		return coords
	}

	root := findNearestOddSquare(location)
	square := root * root
	offset := (location - 1) % (root - 1) // offset from corner

	if location == square {
		coords.x = (root - 1) / 2
		coords.y = -(root - 1) / 2
	} else if square-(root-1) <= location {
		coords.x = offset - (root-1)/2
		coords.y = -(root - 1) / 2
	} else if square-(2*(root-1)) <= location {
		coords.x = -(root - 1) / 2
		coords.y = (root-1)/2 - offset
	} else if square-(3*(root-1)) <= location {
		coords.x = (root-1)/2 - offset
		coords.y = (root - 1) / 2
	} else if square-(4*(root-1)) <= location {
		coords.x = (root - 1) / 2
		coords.y = offset - (root-1)/2
	}

	return coords
}

// distanceToLocation returns manhattan distance to memory location `location`
func distanceToLocation(location int) int {
	return locationToCoordinates(location).manhattanDistance()
}

var _ = Describe("Day3", func() {
	Describe("SpiralMemory", func() {
		Describe("CartesianCoordinates", func() {
			Describe("move", func() {
				It("moves the relative amount", func() {
					here := CartesianCoordinates{11, 22}
					relative := CartesianCoordinates{-1, 5}
					Expect(here.move(relative)).To(Equal(CartesianCoordinates{10, 27}))
				})
			})

			Describe("manhattanDistance", func() {
				It("returns the manhattan distance of the coords", func() {
					Expect(CartesianCoordinates{11, -5}.manhattanDistance()).To(Equal(16))
				})
			})

			Describe("location", func() {
				It("returns the location at the coordinates", func() {
					Expect(coordinatesToLocation(CartesianCoordinates{0, 0})).To(Equal(1))
					Expect(coordinatesToLocation(CartesianCoordinates{1, 0})).To(Equal(2))
					Expect(coordinatesToLocation(CartesianCoordinates{1, 1})).To(Equal(3))
					Expect(coordinatesToLocation(CartesianCoordinates{0, 1})).To(Equal(4))
					Expect(coordinatesToLocation(CartesianCoordinates{-1, 1})).To(Equal(5))
					Expect(coordinatesToLocation(CartesianCoordinates{-1, 0})).To(Equal(6))
					Expect(coordinatesToLocation(CartesianCoordinates{-1, -1})).To(Equal(7))
					Expect(coordinatesToLocation(CartesianCoordinates{0, -1})).To(Equal(8))
					Expect(coordinatesToLocation(CartesianCoordinates{1, -1})).To(Equal(9))
					Expect(coordinatesToLocation(CartesianCoordinates{2, -1})).To(Equal(10))
					Expect(coordinatesToLocation(CartesianCoordinates{2, 0})).To(Equal(11))
					Expect(coordinatesToLocation(CartesianCoordinates{2, 1})).To(Equal(12))
					Expect(coordinatesToLocation(CartesianCoordinates{2, 2})).To(Equal(13))
					Expect(coordinatesToLocation(CartesianCoordinates{1, 2})).To(Equal(14))
					Expect(coordinatesToLocation(CartesianCoordinates{0, 2})).To(Equal(15))
					Expect(coordinatesToLocation(CartesianCoordinates{-1, 2})).To(Equal(16))
					Expect(coordinatesToLocation(CartesianCoordinates{-2, 2})).To(Equal(17))
					Expect(coordinatesToLocation(CartesianCoordinates{-2, 1})).To(Equal(18))
					Expect(coordinatesToLocation(CartesianCoordinates{-2, 0})).To(Equal(19))
					Expect(coordinatesToLocation(CartesianCoordinates{-2, -1})).To(Equal(20))
					Expect(coordinatesToLocation(CartesianCoordinates{-2, -2})).To(Equal(21))
					Expect(coordinatesToLocation(CartesianCoordinates{-1, -2})).To(Equal(22))
					Expect(coordinatesToLocation(CartesianCoordinates{0, -2})).To(Equal(23))
					Expect(coordinatesToLocation(CartesianCoordinates{1, -2})).To(Equal(24))
					Expect(coordinatesToLocation(CartesianCoordinates{2, -2})).To(Equal(25))
					Expect(coordinatesToLocation(CartesianCoordinates{3, -2})).To(Equal(26))
				})
			})
		})

		Describe("locationToCoordinates", func() {
			It("returns the coordinates of a location", func() {
				Expect(locationToCoordinates(1)).To(Equal(CartesianCoordinates{0, 0}))
				Expect(locationToCoordinates(2)).To(Equal(CartesianCoordinates{1, 0}))
				Expect(locationToCoordinates(3)).To(Equal(CartesianCoordinates{1, 1}))
				Expect(locationToCoordinates(4)).To(Equal(CartesianCoordinates{0, 1}))
				Expect(locationToCoordinates(5)).To(Equal(CartesianCoordinates{-1, 1}))
				Expect(locationToCoordinates(6)).To(Equal(CartesianCoordinates{-1, 0}))
				Expect(locationToCoordinates(7)).To(Equal(CartesianCoordinates{-1, -1}))
				Expect(locationToCoordinates(8)).To(Equal(CartesianCoordinates{0, -1}))
				Expect(locationToCoordinates(9)).To(Equal(CartesianCoordinates{1, -1}))
				Expect(locationToCoordinates(10)).To(Equal(CartesianCoordinates{2, -1}))
				Expect(locationToCoordinates(11)).To(Equal(CartesianCoordinates{2, 0}))
				Expect(locationToCoordinates(12)).To(Equal(CartesianCoordinates{2, 1}))
				Expect(locationToCoordinates(13)).To(Equal(CartesianCoordinates{2, 2}))
				Expect(locationToCoordinates(14)).To(Equal(CartesianCoordinates{1, 2}))
				Expect(locationToCoordinates(15)).To(Equal(CartesianCoordinates{0, 2}))
				Expect(locationToCoordinates(16)).To(Equal(CartesianCoordinates{-1, 2}))
				Expect(locationToCoordinates(17)).To(Equal(CartesianCoordinates{-2, 2}))
				Expect(locationToCoordinates(18)).To(Equal(CartesianCoordinates{-2, 1}))
				Expect(locationToCoordinates(19)).To(Equal(CartesianCoordinates{-2, 0}))
				Expect(locationToCoordinates(20)).To(Equal(CartesianCoordinates{-2, -1}))
				Expect(locationToCoordinates(21)).To(Equal(CartesianCoordinates{-2, -2}))
				Expect(locationToCoordinates(22)).To(Equal(CartesianCoordinates{-1, -2}))
				Expect(locationToCoordinates(23)).To(Equal(CartesianCoordinates{0, -2}))
				Expect(locationToCoordinates(24)).To(Equal(CartesianCoordinates{1, -2}))
				Expect(locationToCoordinates(25)).To(Equal(CartesianCoordinates{2, -2}))
				Expect(locationToCoordinates(26)).To(Equal(CartesianCoordinates{3, -2}))
			})
		})

		Describe("distanceToLocation", func() {
			It("returns the manhattan distance to the memory location", func() {
				Expect(distanceToLocation(1)).To(Equal(0))
				Expect(distanceToLocation(2)).To(Equal(1))
				Expect(distanceToLocation(3)).To(Equal(2))
				Expect(distanceToLocation(4)).To(Equal(1))
				Expect(distanceToLocation(5)).To(Equal(2))
				Expect(distanceToLocation(6)).To(Equal(1))
				Expect(distanceToLocation(7)).To(Equal(2))
				Expect(distanceToLocation(8)).To(Equal(1))
				Expect(distanceToLocation(9)).To(Equal(2))

				Expect(distanceToLocation(10)).To(Equal(3))
				Expect(distanceToLocation(11)).To(Equal(2))
				Expect(distanceToLocation(12)).To(Equal(3))
				Expect(distanceToLocation(13)).To(Equal(4))
				Expect(distanceToLocation(14)).To(Equal(3))
				Expect(distanceToLocation(15)).To(Equal(2))
				Expect(distanceToLocation(16)).To(Equal(3))
				Expect(distanceToLocation(17)).To(Equal(4))
				Expect(distanceToLocation(18)).To(Equal(3))
				Expect(distanceToLocation(19)).To(Equal(2))
				Expect(distanceToLocation(20)).To(Equal(3))
				Expect(distanceToLocation(21)).To(Equal(4))
				Expect(distanceToLocation(22)).To(Equal(3))
				Expect(distanceToLocation(23)).To(Equal(2))
				Expect(distanceToLocation(24)).To(Equal(3))
				Expect(distanceToLocation(25)).To(Equal(4))

				Expect(distanceToLocation(26)).To(Equal(5))

				Expect(distanceToLocation(1024)).To(Equal(31))
			})
		})

		Describe("findNearestOddSquare", func() {
			It("returns the root of the nearest odd square (without going under)", func() {
				Expect(findNearestOddSquare(1)).To(Equal(1))
				Expect(findNearestOddSquare(2)).To(Equal(3))
				Expect(findNearestOddSquare(9)).To(Equal(3))
				Expect(findNearestOddSquare(10)).To(Equal(5))
				Expect(findNearestOddSquare(24)).To(Equal(5))
				Expect(findNearestOddSquare(25)).To(Equal(5))
				Expect(findNearestOddSquare(26)).To(Equal(7))
			})
		})

		Describe("puzzle", func() {
			It("star 1", func() {
				fmt.Printf("d3 s1: %d\n", distanceToLocation(265149))
			})
		})
	})
})
