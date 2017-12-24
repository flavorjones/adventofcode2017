package adventofcode2017_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type NumberGenerator struct {
	seed   int
	factor int
}

func NewNumberGenerator(seed int, factor int) *NumberGenerator {
	return &NumberGenerator{seed: seed, factor: factor}
}

func (ng *NumberGenerator) next() int {
	ng.seed = (ng.seed * ng.factor) % 2147483647
	return ng.seed
}

const low16BitMask = 65535

func sameLow16Bits(a, b int) bool {
	return a&low16BitMask == b&low16BitMask
}

func judgeCount(n1, n2 *NumberGenerator) int {
	count := 0
	for j := 0; j < 40000000; j++ {
		if sameLow16Bits(n1.next(), n2.next()) {
			count++
		}
	}
	return count
}

var _ = Describe("Day15", func() {
	Describe("sameLow16Bits", func() {
		It("returns true if the lower 16 bits are identical", func() {
			Expect(sameLow16Bits(1092455, 430625591)).To(BeFalse())
			Expect(sameLow16Bits(1181022009, 1233683848)).To(BeFalse())
			Expect(sameLow16Bits(245556042, 1431495498)).To(BeTrue())
			Expect(sameLow16Bits(1744312007, 137874439)).To(BeFalse())
			Expect(sameLow16Bits(1352636452, 285222916)).To(BeFalse())
		})
	})

	Describe("NumberGenerator", func() {
		Describe("NewNumberGenerator", func() {
			It("takes a seed value", func() {
				ng := NewNumberGenerator(65, 16807)
				Expect(ng.seed).To(Equal(65))
			})
		})

		Describe("next", func() {
			It("generates the next value, and saves that as the next seed value", func() {
				ng := NewNumberGenerator(65, 16807)
				val := ng.next()
				Expect(val).To(Equal(1092455))
				Expect(ng.seed).To(Equal(1092455))
			})

			It("generates a predictable sequence of numbers", func() {
				ng := NewNumberGenerator(65, 16807)
				Expect(ng.next()).To(Equal(1092455))
				Expect(ng.next()).To(Equal(1181022009))
				Expect(ng.next()).To(Equal(245556042))
				Expect(ng.next()).To(Equal(1744312007))
				Expect(ng.next()).To(Equal(1352636452))

				ng = NewNumberGenerator(8921, 48271)
				Expect(ng.next()).To(Equal(430625591))
				Expect(ng.next()).To(Equal(1233683848))
				Expect(ng.next()).To(Equal(1431495498))
				Expect(ng.next()).To(Equal(137874439))
				Expect(ng.next()).To(Equal(285222916))
			})
		})
	})

	Describe("judgeCount", func() {
		It("finds the number of samelow16bits in the first 40 million numbers", func() {
			n1 := NewNumberGenerator(65, 16807)
			n2 := NewNumberGenerator(8921, 48271)
			Expect(judgeCount(n1, n2)).To(Equal(588))
		})
	})

	Describe("puzzle", func() {
		It("solves star 1", func() {
			n1 := NewNumberGenerator(277, 16807)
			n2 := NewNumberGenerator(349, 48271)
			count := judgeCount(n1, n2)
			fmt.Printf("d15 s1: judge counted %d numbers", count)
		})
	})
})
