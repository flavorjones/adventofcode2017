package adventofcode2017_test

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SpiralMemoryLocation int

type CartesianCoordinates struct {
	x int
	y int
}

// returns the manhattan distance of the coordinates
func (c CartesianCoordinates) manhattanDistance() int {
	return int(math.Abs(float64(c.x)) + math.Abs(float64(c.y)))
}

func (c CartesianCoordinates) move(relative CartesianCoordinates) CartesianCoordinates {
	return CartesianCoordinates{c.x + relative.x, c.y + relative.y}
}

// returns the location of the coordinates
func (c CartesianCoordinates) location() SpiralMemoryLocation {
	if (c == CartesianCoordinates{0, 0}) {
		return 1
	}

	min_root := (2 * int(math.Max(math.Abs(float64(c.x)), math.Abs(float64(c.y))))) - 1
	min_square := SpiralMemoryLocation(math.Pow(float64(min_root), 2.0))
	max_square := SpiralMemoryLocation(math.Pow(float64(min_root+2), 2.0))

	// brute force it because I'm lazy
	for j := SpiralMemoryLocation(min_square + 1); j <= max_square; j++ {
		if j.coordinates() == c {
			return j
		}
	}

	return -1
}

// returns the cartesian coordinates of `location`
func (location SpiralMemoryLocation) coordinates() CartesianCoordinates {
	coords := CartesianCoordinates{}

	// special case
	if location == 1 {
		return coords
	}

	root := findNearestOddSquare(int(location))
	square := root * root
	offset := (int(location) - 1) % (root - 1) // offset from corner

	if int(location) == square {
		coords.x = (root - 1) / 2
		coords.y = -(root - 1) / 2
	} else if square-(root-1) <= int(location) {
		coords.x = offset - (root-1)/2
		coords.y = -(root - 1) / 2
	} else if square-(2*(root-1)) <= int(location) {
		coords.x = -(root - 1) / 2
		coords.y = (root-1)/2 - offset
	} else if square-(3*(root-1)) <= int(location) {
		coords.x = (root-1)/2 - offset
		coords.y = (root - 1) / 2
	} else if square-(4*(root-1)) <= int(location) {
		coords.x = (root - 1) / 2
		coords.y = offset - (root-1)/2
	}

	return coords
}

// returns manhattan distance to memory location `location`
func (location SpiralMemoryLocation) distance() int {
	return location.coordinates().manhattanDistance()
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
					Expect(CartesianCoordinates{0, 0}.location()).To(Equal(SpiralMemoryLocation(1)))
					Expect(CartesianCoordinates{1, 0}.location()).To(Equal(SpiralMemoryLocation(2)))
					Expect(CartesianCoordinates{1, 1}.location()).To(Equal(SpiralMemoryLocation(3)))
					Expect(CartesianCoordinates{0, 1}.location()).To(Equal(SpiralMemoryLocation(4)))
					Expect(CartesianCoordinates{-1, 1}.location()).To(Equal(SpiralMemoryLocation(5)))
					Expect(CartesianCoordinates{-1, 0}.location()).To(Equal(SpiralMemoryLocation(6)))
					Expect(CartesianCoordinates{-1, -1}.location()).To(Equal(SpiralMemoryLocation(7)))
					Expect(CartesianCoordinates{0, -1}.location()).To(Equal(SpiralMemoryLocation(8)))
					Expect(CartesianCoordinates{1, -1}.location()).To(Equal(SpiralMemoryLocation(9)))
					Expect(CartesianCoordinates{2, -1}.location()).To(Equal(SpiralMemoryLocation(10)))
					Expect(CartesianCoordinates{2, 0}.location()).To(Equal(SpiralMemoryLocation(11)))
					Expect(CartesianCoordinates{2, 1}.location()).To(Equal(SpiralMemoryLocation(12)))
					Expect(CartesianCoordinates{2, 2}.location()).To(Equal(SpiralMemoryLocation(13)))
					Expect(CartesianCoordinates{1, 2}.location()).To(Equal(SpiralMemoryLocation(14)))
					Expect(CartesianCoordinates{0, 2}.location()).To(Equal(SpiralMemoryLocation(15)))
					Expect(CartesianCoordinates{-1, 2}.location()).To(Equal(SpiralMemoryLocation(16)))
					Expect(CartesianCoordinates{-2, 2}.location()).To(Equal(SpiralMemoryLocation(17)))
					Expect(CartesianCoordinates{-2, 1}.location()).To(Equal(SpiralMemoryLocation(18)))
					Expect(CartesianCoordinates{-2, 0}.location()).To(Equal(SpiralMemoryLocation(19)))
					Expect(CartesianCoordinates{-2, -1}.location()).To(Equal(SpiralMemoryLocation(20)))
					Expect(CartesianCoordinates{-2, -2}.location()).To(Equal(SpiralMemoryLocation(21)))
					Expect(CartesianCoordinates{-1, -2}.location()).To(Equal(SpiralMemoryLocation(22)))
					Expect(CartesianCoordinates{0, -2}.location()).To(Equal(SpiralMemoryLocation(23)))
					Expect(CartesianCoordinates{1, -2}.location()).To(Equal(SpiralMemoryLocation(24)))
					Expect(CartesianCoordinates{2, -2}.location()).To(Equal(SpiralMemoryLocation(25)))
					Expect(CartesianCoordinates{3, -2}.location()).To(Equal(SpiralMemoryLocation(26)))
				})
			})
		})

		Describe("SpiralMemoryLocation", func() {
			Describe("coordinates", func() {
				It("returns the coordinates of a location", func() {
					Expect(SpiralMemoryLocation(1).coordinates()).To(Equal(CartesianCoordinates{0, 0}))
					Expect(SpiralMemoryLocation(2).coordinates()).To(Equal(CartesianCoordinates{1, 0}))
					Expect(SpiralMemoryLocation(3).coordinates()).To(Equal(CartesianCoordinates{1, 1}))
					Expect(SpiralMemoryLocation(4).coordinates()).To(Equal(CartesianCoordinates{0, 1}))
					Expect(SpiralMemoryLocation(5).coordinates()).To(Equal(CartesianCoordinates{-1, 1}))
					Expect(SpiralMemoryLocation(6).coordinates()).To(Equal(CartesianCoordinates{-1, 0}))
					Expect(SpiralMemoryLocation(7).coordinates()).To(Equal(CartesianCoordinates{-1, -1}))
					Expect(SpiralMemoryLocation(8).coordinates()).To(Equal(CartesianCoordinates{0, -1}))
					Expect(SpiralMemoryLocation(9).coordinates()).To(Equal(CartesianCoordinates{1, -1}))
					Expect(SpiralMemoryLocation(10).coordinates()).To(Equal(CartesianCoordinates{2, -1}))
					Expect(SpiralMemoryLocation(11).coordinates()).To(Equal(CartesianCoordinates{2, 0}))
					Expect(SpiralMemoryLocation(12).coordinates()).To(Equal(CartesianCoordinates{2, 1}))
					Expect(SpiralMemoryLocation(13).coordinates()).To(Equal(CartesianCoordinates{2, 2}))
					Expect(SpiralMemoryLocation(14).coordinates()).To(Equal(CartesianCoordinates{1, 2}))
					Expect(SpiralMemoryLocation(15).coordinates()).To(Equal(CartesianCoordinates{0, 2}))
					Expect(SpiralMemoryLocation(16).coordinates()).To(Equal(CartesianCoordinates{-1, 2}))
					Expect(SpiralMemoryLocation(17).coordinates()).To(Equal(CartesianCoordinates{-2, 2}))
					Expect(SpiralMemoryLocation(18).coordinates()).To(Equal(CartesianCoordinates{-2, 1}))
					Expect(SpiralMemoryLocation(19).coordinates()).To(Equal(CartesianCoordinates{-2, 0}))
					Expect(SpiralMemoryLocation(20).coordinates()).To(Equal(CartesianCoordinates{-2, -1}))
					Expect(SpiralMemoryLocation(21).coordinates()).To(Equal(CartesianCoordinates{-2, -2}))
					Expect(SpiralMemoryLocation(22).coordinates()).To(Equal(CartesianCoordinates{-1, -2}))
					Expect(SpiralMemoryLocation(23).coordinates()).To(Equal(CartesianCoordinates{0, -2}))
					Expect(SpiralMemoryLocation(24).coordinates()).To(Equal(CartesianCoordinates{1, -2}))
					Expect(SpiralMemoryLocation(25).coordinates()).To(Equal(CartesianCoordinates{2, -2}))
					Expect(SpiralMemoryLocation(26).coordinates()).To(Equal(CartesianCoordinates{3, -2}))
				})
			})
		})

		Describe("distanceToLocation", func() {
			It("returns the manhattan distance to the memory location", func() {
				Expect(SpiralMemoryLocation(1).distance()).To(Equal(0))
				Expect(SpiralMemoryLocation(2).distance()).To(Equal(1))
				Expect(SpiralMemoryLocation(3).distance()).To(Equal(2))
				Expect(SpiralMemoryLocation(4).distance()).To(Equal(1))
				Expect(SpiralMemoryLocation(5).distance()).To(Equal(2))
				Expect(SpiralMemoryLocation(6).distance()).To(Equal(1))
				Expect(SpiralMemoryLocation(7).distance()).To(Equal(2))
				Expect(SpiralMemoryLocation(8).distance()).To(Equal(1))
				Expect(SpiralMemoryLocation(9).distance()).To(Equal(2))

				Expect(SpiralMemoryLocation(10).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(11).distance()).To(Equal(2))
				Expect(SpiralMemoryLocation(12).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(13).distance()).To(Equal(4))
				Expect(SpiralMemoryLocation(14).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(15).distance()).To(Equal(2))
				Expect(SpiralMemoryLocation(16).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(17).distance()).To(Equal(4))
				Expect(SpiralMemoryLocation(18).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(19).distance()).To(Equal(2))
				Expect(SpiralMemoryLocation(20).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(21).distance()).To(Equal(4))
				Expect(SpiralMemoryLocation(22).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(23).distance()).To(Equal(2))
				Expect(SpiralMemoryLocation(24).distance()).To(Equal(3))
				Expect(SpiralMemoryLocation(25).distance()).To(Equal(4))

				Expect(SpiralMemoryLocation(26).distance()).To(Equal(5))

				Expect(SpiralMemoryLocation(1024).distance()).To(Equal(31))
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
				fmt.Printf("d3 s1: %d\n", SpiralMemoryLocation(265149).distance())
			})
		})
	})
})
