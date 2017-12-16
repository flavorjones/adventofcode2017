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
func manhattanDistance(coordinates CartesianCoordinates) int {
	return int(math.Abs(float64(coordinates.x)) + math.Abs(float64(coordinates.y)))
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
	position := locationToCoordinates(location)
	return manhattanDistance(position)
}

var _ = Describe("Day3", func() {
	Describe("SpiralMemory", func() {
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
