package adventofcode2017_test

import (
	"fmt"
	"regexp"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type KnotHashList []byte

type KnotHash struct {
	list     KnotHashList
	position int
	skip     int
}

func NewKnotHash(size int) *KnotHash {
	list := make(KnotHashList, size)
	for j := 0; j < size; j++ {
		list[j] = byte(j)
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
	return int(kh.list[0]) * int(kh.list[1])
}

var SEQUENCE_SUFFIX = []byte{17, 31, 73, 47, 23}
var HASH_ROUNDS = 64

func (kh *KnotHash) fullHash(lengthsDescriptor string) string {
	lengths := append([]byte(lengthsDescriptor), SEQUENCE_SUFFIX...)

	for jround := 0; jround < HASH_ROUNDS; jround++ {
		for j := 0; j < len(lengths); j++ {
			kh.hashStep(int(lengths[j]))
		}
	}

	// make dense hash
	denseHash := densify(kh.list)

	// return hex
	return hexify(denseHash)
}

func (kh *KnotHash) hashStep(length int) {
	kh.list = reverseSubSlice(kh.list, kh.position, length)
	kh.position = (kh.position + length + kh.skip) % len(kh.list)
	kh.skip += 1
}

func reverseSubSlice(slice KnotHashList, start int, length int) KnotHashList {
	if length <= 1 {
		return slice
	}

	var subslice KnotHashList

	if start+length < len(slice) {
		subslice = slice[start : start+length]
	} else {
		subslice = make(KnotHashList, length)
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

func hexify(bytes []byte) string {
	return fmt.Sprintf("%x", bytes)
}

var DENSIFY_BLOCK_SIZE = 16

func densify(numbers KnotHashList) KnotHashList {
	if len(numbers)%DENSIFY_BLOCK_SIZE != 0 {
		panic(fmt.Sprintf("error: len of slice (%d) is not divisible by %d",
			len(numbers), DENSIFY_BLOCK_SIZE))
	}

	nblocks := len(numbers) / DENSIFY_BLOCK_SIZE
	rval := make([]byte, nblocks)

	for jblock := 0; jblock < nblocks; jblock++ {
		blockVal := numbers[jblock*DENSIFY_BLOCK_SIZE]
		for jbyte := 1; jbyte < DENSIFY_BLOCK_SIZE; jbyte++ {
			index := jblock*DENSIFY_BLOCK_SIZE + jbyte
			blockVal = blockVal ^ numbers[index]
		}
		rval[jblock] = blockVal
	}

	return rval
}

var _ = Describe("Day10", func() {
	Describe("KnotHash", func() {
		Describe("NewKnotHash", func() {
			It("has a sane initial state", func() {
				kh := NewKnotHash(5)
				Expect(kh.position).To(Equal(0))
				Expect(kh.skip).To(Equal(0))
				Expect(kh.list).To(Equal(KnotHashList{0, 1, 2, 3, 4}))
			})
		})

		Describe("hash()", func() {
			It("calculates the proper hash", func() {
				Expect(NewKnotHash(5).hash("3, 4, 1, 5")).To(Equal(12))
			})
		})

		Describe("fullHash", func() {
			It("calculates the proper hash", func() {
				Expect(NewKnotHash(256).fullHash("")).
					To(Equal("a2582a3a0e66e6e86e3812dcb672a272"))

				Expect(NewKnotHash(256).fullHash("AoC 2017")).
					To(Equal("33efeb34ea91902bb2f59c9920caa6cd"))

				Expect(NewKnotHash(256).fullHash("1,2,3")).
					To(Equal("3efbe78a8d82f29979031a4aa0b16a9d"))

				Expect(NewKnotHash(256).fullHash("1,2,4")).
					To(Equal("63960835bcdc130f0b66d7ff4f6a5a8e"))
			})
		})

		Describe("hashStep", func() {
			It("performs the correct transformations", func() {
				kh := NewKnotHash(5)

				kh.hashStep(3)
				Expect(kh.position).To(Equal(3))
				Expect(kh.skip).To(Equal(1))
				Expect(kh.list).To(Equal(KnotHashList{2, 1, 0, 3, 4}))

				kh.hashStep(4)
				Expect(kh.position).To(Equal(3))
				Expect(kh.skip).To(Equal(2))
				Expect(kh.list).To(Equal(KnotHashList{4, 3, 0, 1, 2}))

				kh.hashStep(1)
				Expect(kh.position).To(Equal(1))
				Expect(kh.skip).To(Equal(3))
				Expect(kh.list).To(Equal(KnotHashList{4, 3, 0, 1, 2}))

				kh.hashStep(5)
				Expect(kh.position).To(Equal(4))
				Expect(kh.skip).To(Equal(4))
				Expect(kh.list).To(Equal(KnotHashList{3, 4, 2, 1, 0}))

				kh.hashStep(0)
				Expect(kh.position).To(Equal(3))
				Expect(kh.skip).To(Equal(5))
				Expect(kh.list).To(Equal(KnotHashList{3, 4, 2, 1, 0}))
			})
		})
	})

	Describe("hexify", func() {
		It("renders bytes as hex characters", func() {
			Expect(hexify([]byte{64, 7, 255})).To(Equal("4007ff"))
		})
	})

	Describe("densify", func() {
		It("xors each byte of a 16-byte block", func() {
			Expect(densify([]byte{65, 27, 9, 1, 4, 3, 40, 50, 91, 7, 6, 0, 2, 5, 68, 22})).
				To(Equal(KnotHashList{64}))

			Expect(densify([]byte{65, 27, 9, 1, 4, 3, 40, 50, 91, 7, 6, 0, 2, 5, 68, 22, 65, 27, 9, 1, 4, 3, 40, 50, 91, 7, 6, 0, 2, 5, 68, 22})).
				To(Equal(KnotHashList{64, 64}))
		})
	})

	Describe("puzzle", func() {
		lengthsDescriptor := `88,88,211,106,141,1,78,254,2,111,77,255,90,0,54,205`

		It("solves star 1", func() {
			kh := NewKnotHash(256)
			hash := kh.hash(lengthsDescriptor)
			fmt.Printf("d10 s1: hash values is %d\n", hash)
		})

		It("solves star 2", func() {
			kh := NewKnotHash(256)
			hash := kh.fullHash(lengthsDescriptor)
			fmt.Printf("d10 s2: hash values is %s\n", hash)
		})
	})
})
