package adventofcode2017_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SpinLockBuffer []int

type SpinLock struct {
	stepSize int
	buffer   SpinLockBuffer
	position int
	count    int
}

func NewSpinLock(stepSize int) *SpinLock {
	return &SpinLock{position: 0, stepSize: stepSize, buffer: []int{0}}
}

func (s *SpinLock) insert() {
	s.count++
	s.position = (s.position+s.stepSize)%len(s.buffer) + 1

	// insert without allocating a whole new slice
	s.buffer = append(s.buffer, 0)
	copy(s.buffer[s.position+1:], s.buffer[s.position:])
	s.buffer[s.position] = s.count
}

func (s *SpinLock) insert2017() {
	for j := 1; j <= 2017; j++ {
		s.insert()
	}
}

var _ = Describe("Day17", func() {
	Describe("SpinLock", func() {
		Describe("NewSpinLock", func() {
			It("takes a step size argument", func() {
				s := NewSpinLock(12)
				Expect(s.stepSize).To(Equal(12))
			})

			It("should have an initial circular buffer state", func() {
				s := NewSpinLock(12)
				Expect(s.buffer).To(Equal(SpinLockBuffer{0}))
			})

			It("should have an initial pointer position", func() {
				s := NewSpinLock(12)
				Expect(s.position).To(Equal(0))
			})
		})

		Describe("insert()", func() {
			It("steps forward and inserts", func() {
				s := NewSpinLock(3)

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 1}))
				Expect(s.position).To(Equal(1))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 2, 1}))
				Expect(s.position).To(Equal(1))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 2, 3, 1}))
				Expect(s.position).To(Equal(2))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 2, 4, 3, 1}))
				Expect(s.position).To(Equal(2))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 5, 2, 4, 3, 1}))
				Expect(s.position).To(Equal(1))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 5, 2, 4, 3, 6, 1}))
				Expect(s.position).To(Equal(5))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 5, 7, 2, 4, 3, 6, 1}))
				Expect(s.position).To(Equal(2))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 5, 7, 2, 4, 3, 8, 6, 1}))
				Expect(s.position).To(Equal(6))

				s.insert()
				Expect(s.buffer).To(Equal(SpinLockBuffer{0, 9, 5, 7, 2, 4, 3, 8, 6, 1}))
				Expect(s.position).To(Equal(1))

				for j := 10; j <= 2017; j++ {
					s.insert()
				}
				Expect(s.buffer[s.position]).To(Equal(2017))
				Expect(s.buffer[s.position+1]).To(Equal(638))
			})
		})
	})

	Describe("puzzle", func() {
		It("solves star 1", func() {
			s := NewSpinLock(312)
			s.insert2017()
			Expect(s.buffer[s.position]).To(Equal(2017))
			answer := s.buffer[s.position+1]
			fmt.Printf("d17 s1: short-circuit spinlock with %d\n", answer)
		})
	})
})
