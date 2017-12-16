package adventofcode2017_test

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CartesianPosition struct {
	x int
	y int
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

// distanceToLocation returns manhattan distance to memory location `location`
func distanceToLocation(location int) int {
	root := findNearestOddSquare(location)

	// special case
	if root == 1 {
		return 0
	}

	// offset from a corner
	offset := (location - 1) % (root - 1)
	distance := (root-1)/2 + int(math.Abs(float64(offset-(root-1)/2)))
	return distance
}

var _ = Describe("Day3", func() {
	Describe("SpiralMemory", func() {
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
