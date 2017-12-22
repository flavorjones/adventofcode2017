package adventofcode2017_test

import (
	"fmt"
	"regexp"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type KnotHash struct {
	list     []int
	position int
	skip     int
}

func NewKnotHash(size int) *KnotHash {
	list := make([]int, size)
	for j := 0; j < size; j++ {
		list[j] = j
	}
	return &KnotHash{list: list}
}

var lengthsSeparatorRe = regexp.MustCompile(`\s*,\s*`)

func (kh *KnotHash) hash(lengthsDescriptor string) int {
	lengths := lengthsSeparatorRe.Split(lengthsDescriptor, -1)
	for j := 0; j < len(lengths); j++ {
		length, _ := strconv.Atoi(lengths[j])
		kh.hashStep(length)
	}
	return kh.list[0] * kh.list[1]
}

func (kh *KnotHash) hashStep(length int) {
	kh.list = reverseSubSlice(kh.list, kh.position, length)
	kh.position = (kh.position + length + kh.skip) % len(kh.list)
	kh.skip += 1
}

func reverseSubSlice(slice []int, start int, length int) []int {
	if length <= 1 {
		return slice
	}

	var subslice []int

	if start+length < len(slice) {
		subslice = slice[start : start+length]
	} else {
		subslice = make([]int, length)
		copy(subslice[0:len(slice)-start], slice[start:])
		copy(subslice[len(slice)-start:], slice[0:length-(len(slice)-start)])
	}

	// from https://github.com/golang/go/wiki/SliceTricks
	for j := len(subslice)/2 - 1; j >= 0; j-- {
		opp := len(subslice) - 1 - j
		subslice[j], subslice[opp] = subslice[opp], subslice[j]
	}

	if start+length < len(slice) {
		copy(slice[start:start+length-1], subslice)
	} else {
		copy(slice[start:], subslice[0:len(slice)-start])
		copy(slice[0:length-(len(slice)-start)], subslice[len(slice)-start:])
	}

	return slice
}

var _ = Describe("Day10", func() {
	Describe("KnotHash", func() {
		Describe("NewKnotHash", func() {
			It("has a sane initial state", func() {
				kh := NewKnotHash(5)
				Expect(kh.position).To(Equal(0))
				Expect(kh.skip).To(Equal(0))
				Expect(kh.list).To(Equal([]int{0, 1, 2, 3, 4}))
			})
		})

		Describe("hash()", func() {
			It("calculates the proper hash", func() {
				Expect(NewKnotHash(5).hash("3, 4, 1, 5")).To(Equal(12))
			})
		})

		Describe("hashStep", func() {
			It("performs the correct transformations", func() {
				kh := NewKnotHash(5)

				kh.hashStep(3)
				Expect(kh.position).To(Equal(3))
				Expect(kh.skip).To(Equal(1))
				Expect(kh.list).To(Equal([]int{2, 1, 0, 3, 4}))

				kh.hashStep(4)
				Expect(kh.position).To(Equal(3))
				Expect(kh.skip).To(Equal(2))
				Expect(kh.list).To(Equal([]int{4, 3, 0, 1, 2}))

				kh.hashStep(1)
				Expect(kh.position).To(Equal(1))
				Expect(kh.skip).To(Equal(3))
				Expect(kh.list).To(Equal([]int{4, 3, 0, 1, 2}))

				kh.hashStep(5)
				Expect(kh.position).To(Equal(4))
				Expect(kh.skip).To(Equal(4))
				Expect(kh.list).To(Equal([]int{3, 4, 2, 1, 0}))

				kh.hashStep(0)
				Expect(kh.position).To(Equal(3))
				Expect(kh.skip).To(Equal(5))
				Expect(kh.list).To(Equal([]int{3, 4, 2, 1, 0}))
			})
		})
	})

	Describe("puzzle", func() {
		lengthsDescriptor := `88,88,211,106,141,1,78,254,2,111,77,255,90,0,54,205`

		It("solves star 1", func() {
			kh := NewKnotHash(256)
			hash := kh.hash(lengthsDescriptor)
			fmt.Printf("d10 s1: hash values is %d\n", hash)
		})
	})
})
