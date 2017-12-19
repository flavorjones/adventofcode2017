package adventofcode2017_test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deckarep/golang-set"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MemoryBankSet struct {
	banks []int
}

func NewMemoryBankSet(banks_decl string) *MemoryBankSet {
	banks_list := strings.Fields(banks_decl)
	banks := make([]int, len(banks_list))
	for j, bank := range banks_list {
		banks[j], _ = strconv.Atoi(bank)
	}
	return &MemoryBankSet{banks: banks}
}

func (mbs *MemoryBankSet) tick() {
	jlargest := findIndexOfLargest(mbs.banks)
	blocks := mbs.banks[jlargest]
	mbs.banks[jlargest] = 0

	for j := jlargest + 1; blocks > 0; j++ {
		if j >= len(mbs.banks) {
			j = 0
		}

		mbs.banks[j] += 1
		blocks -= 1
	}
}

func (mbs *MemoryBankSet) debug() int {
	cache := mapset.NewSet()
	steps := 0

	for {
		mbs.tick()
		steps += 1
		key := makeKeyFrom(mbs.banks)
		if cache.Contains(key) {
			return steps
		}
		cache.Add(key)
	}
}

func findIndexOfLargest(int_slice []int) int {
	jlargest := -1
	max := -1

	for j := 0; j < len(int_slice); j++ {
		if int_slice[j] > max {
			jlargest = j
			max = int_slice[j]
		}
	}

	return jlargest
}

func makeKeyFrom(int_slice []int) string {
	pieces := make([]string, len(int_slice))
	for j, number := range int_slice {
		pieces[j] = strconv.Itoa(number)
	}
	return strings.Join(pieces, ",")
}

var _ = Describe("Day6", func() {
	Describe("findIndexOfLargest", func() {
		It("finds the index of the largest, ties go to the lower index", func() {
			Expect(findIndexOfLargest([]int{1, 3, 5, 7, 3})).To(Equal(3))
			Expect(findIndexOfLargest([]int{1, 7, 5, 3, 3})).To(Equal(1))
			Expect(findIndexOfLargest([]int{1, 7, 5, 7, 3})).To(Equal(1))
		})
	})

	Describe("makeKeyFrom", func() {
		It("returns a string representing the int slice", func() {
			Expect(makeKeyFrom([]int{1, 3, 5, 7, 3})).To(Equal("1,3,5,7,3"))
		})
	})

	Describe("NewMemoryBankSet", func() {
		It("creates a memory bank with the right initial state", func() {
			mbs := NewMemoryBankSet("0 9 1 8 2 8 3 7 4 6")
			Expect(mbs.banks).To(Equal([]int{0, 9, 1, 8, 2, 8, 3, 7, 4, 6}))
		})
	})

	Describe("MemoryBankSet", func() {
		Describe("tick", func() {
			It("rebalances the largest bank", func() {
				mbs := NewMemoryBankSet("0 2 7 0")

				mbs.tick()
				Expect(mbs.banks).To(Equal([]int{2, 4, 1, 2}))

				mbs.tick()
				Expect(mbs.banks).To(Equal([]int{3, 1, 2, 3}))

				mbs.tick()
				Expect(mbs.banks).To(Equal([]int{0, 2, 3, 4}))

				mbs.tick()
				Expect(mbs.banks).To(Equal([]int{1, 3, 4, 1}))

				mbs.tick()
				Expect(mbs.banks).To(Equal([]int{2, 4, 1, 2}))
			})
		})

		Describe("debug", func() {
			It("runs until it sees the same state again", func() {
				mbs := NewMemoryBankSet("0 2 7 0")
				mbs.debug()
				Expect(mbs.banks).To(Equal([]int{2, 4, 1, 2}))
			})

			It("returns the number of steps taken", func() {
				mbs := NewMemoryBankSet("0 2 7 0")
				Expect(mbs.debug()).To(Equal(5))
			})
		})
	})

	Describe("puzzle", func() {
		It("solves star 1", func() {
			mbs := NewMemoryBankSet("5	1	10	0	1	7	13	14	3	12	8	10	7	12	0	6")
			steps := mbs.debug()
			fmt.Printf("d6 s1: took %d steps to find infinite loop\n", steps)
		})
	})
})
