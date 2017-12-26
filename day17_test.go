package adventofcode2017_test

import (
	"container/list"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SpinLock struct {
	stepSize int
	buffer   *list.List
	cursor   *list.Element
	count    int
}

func NewSpinLock(stepSize int) *SpinLock {
	list := list.New()
	list.PushFront(0)
	cursor := list.Front()
	return &SpinLock{stepSize: stepSize, buffer: list, cursor: cursor}
}

func (s *SpinLock) advanceCursor() {
	s.cursor = s.cursor.Next()
	if s.cursor == nil {
		s.cursor = s.buffer.Front()
	}
}

func (s *SpinLock) insert() {
	s.count++
	for j := 0; j < s.stepSize; j++ {
		s.advanceCursor()
	}
	s.buffer.InsertAfter(s.count, s.cursor)
	s.advanceCursor()
}

func (s *SpinLock) insertN(n int) {
	for j := 1; j <= n; j++ {
		if (j % 100000) == 0 {
			fmt.Printf("insert %d", j)
		}
		s.insert()
	}
}

func (s *SpinLock) cursorOf(desired int) *list.Element {
	for e := s.buffer.Front(); e != nil; e = e.Next() {
		if e.Value == desired {
			return e
		}
	}
	return nil
}

// primarily for testing
func (s *SpinLock) toSlice() []int {
	rval := make([]int, s.buffer.Len())
	for j, e := 0, s.buffer.Front(); e != nil; j, e = j+1, e.Next() {
		rval[j] = e.Value.(int)
	}
	return rval
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
				Expect(s.toSlice()).To(Equal([]int{0}))
			})

			It("should have an initial cursor", func() {
				s := NewSpinLock(12)
				Expect(s.cursor.Value).To(Equal(0))
			})
		})

		Describe("insert()", func() {
			It("steps forward and inserts", func() {
				s := NewSpinLock(3)

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 1}))
				Expect(s.cursor.Value).To(Equal(1))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 2, 1}))
				Expect(s.cursor.Value).To(Equal(2))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 2, 3, 1}))
				Expect(s.cursor.Value).To(Equal(3))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 2, 4, 3, 1}))
				Expect(s.cursor.Value).To(Equal(4))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 5, 2, 4, 3, 1}))
				Expect(s.cursor.Value).To(Equal(5))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 5, 2, 4, 3, 6, 1}))
				Expect(s.cursor.Value).To(Equal(6))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 5, 7, 2, 4, 3, 6, 1}))
				Expect(s.cursor.Value).To(Equal(7))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 5, 7, 2, 4, 3, 8, 6, 1}))
				Expect(s.cursor.Value).To(Equal(8))

				s.insert()
				Expect(s.toSlice()).To(Equal([]int{0, 9, 5, 7, 2, 4, 3, 8, 6, 1}))
				Expect(s.cursor.Value).To(Equal(9))

				for j := 10; j <= 2017; j++ {
					s.insert()
				}
				Expect(s.cursor.Value).To(Equal(2017))
				Expect(s.cursor.Next().Value).To(Equal(638))
			})
		})

		Describe("cursorOf()", func() {
			It("returns the index of the number in the buffer", func() {
				s := NewSpinLock(3)
				s.insertN(9)
				Expect(s.cursorOf(4).Value).To(Equal(4))
				Expect(s.cursorOf(0).Value).To(Equal(0))
				Expect(s.cursorOf(1).Value).To(Equal(1))
			})
		})
	})

	Describe("puzzle", func() {
		It("solves star 1", func() {
			s := NewSpinLock(312)
			s.insertN(2017)
			Expect(s.cursor.Value).To(Equal(2017))
			answer := s.cursor.Next().Value
			fmt.Printf("d17 s1: short-circuit spinlock with %d\n", answer)
		})

		// It("solves star 2", func() {
		// 	s := NewSpinLock(312)
		// 	s.insertN(50000000)
		// 	index := s.indexOf(0)
		// 	answer := s.buffer[index+1]
		// 	fmt.Printf("d17 s2: short-circuit spinlock with %d\n", answer)
		// })
	})
})
