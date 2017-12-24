package adventofcode2017_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	diskHeight = 128 // nrows
	diskWidth  = 128 // ncols
)

type Disk struct {
	rows []string
}

func NewDisk(key string) *Disk {
	rows := make([]string, diskHeight)
	for j := 0; j < diskHeight; j++ {
		rowkey := fmt.Sprintf("%s-%d", key, j)
		rows[j] = NewKnotHash(256).fullHash(rowkey)
	}
	return &Disk{rows: rows}
}

var asciiHexVal = map[byte]byte{
	'0': 0, '1': 1, '2': 2, '3': 3,
	'4': 4, '5': 5, '6': 6, '7': 7,
	'8': 8, '9': 9, 'a': 10, 'b': 11,
	'c': 12, 'd': 13, 'e': 14, 'f': 15,
}

func (d *Disk) used(row, col int) bool {
	nbyte := col / 4
	nbit := col % 4
	mask := byte(1 << uint(3-nbit)) // low bit becomes high bit

	val, ok := asciiHexVal[d.rows[row][nbyte]]
	if !ok {
		panic(fmt.Sprintf("error: could not find val for `%s`", d.rows[row][nbyte]))
	}

	return val&mask > 0
}

func (d *Disk) usedCount() int {
	count := 0
	for jrow := 0; jrow < diskHeight; jrow++ {
		for jcol := 0; jcol < diskWidth; jcol++ {
			if d.used(jrow, jcol) {
				count++
			}
		}
	}
	return count
}

func (d *Disk) regionCount() int {
	// make a mutable copy of the used blocks
	bitmap := make([][]bool, diskHeight)
	for j := 0; j < diskHeight; j++ {
		bitmap[j] = make([]bool, diskWidth)
	}

	for jrow := 0; jrow < diskHeight; jrow++ {
		for jcol := 0; jcol < diskWidth; jcol++ {
			if d.used(jrow, jcol) {
				bitmap[jrow][jcol] = true
			}
		}
	}

	count := 0
	for jrow := 0; jrow < diskHeight; jrow++ {
		for jcol := 0; jcol < diskWidth; jcol++ {
			if bitmap[jrow][jcol] {
				count++
				pos := CartesianCoordinates{x: jrow, y: jcol}
				markAdjacent(&bitmap, pos)
			}
		}
	}

	return count
}

var defragAdjacentCells = []CartesianCoordinates{
	CartesianCoordinates{-1, 0},
	CartesianCoordinates{0, -1},
	CartesianCoordinates{0, 1},
	CartesianCoordinates{1, 0},
}

func markAdjacent(bitmapp *([][]bool), p CartesianCoordinates) {
	bitmap := (*bitmapp)
	if !bitmap[p.x][p.y] {
		panic(fmt.Sprintf("error: [%d, %d] is not free", p.x, p.y))
	}

	bitmap[p.x][p.y] = false

	for _, translation := range defragAdjacentCells {
		np := p.move(translation)
		if np.x >= 0 && np.x < len(bitmap) && np.y >= 0 && np.y < len(bitmap[np.x]) {
			if bitmap[np.x][np.y] {
				markAdjacent(bitmapp, np)
			}
		}
	}
}

var _ = Describe("Day14", func() {
	Describe("Disk", func() {
		Describe("NewDisk", func() {
			It("takes a string and uses KnotHash to calculate free and used blocks", func() {
				d := NewDisk("flqrgnkx")

				Expect(d.used(0, 0)).To(BeTrue())
				Expect(d.used(0, 1)).To(BeTrue())
				Expect(d.used(0, 2)).To(BeFalse())
				Expect(d.used(0, 3)).To(BeTrue())
				Expect(d.used(0, 4)).To(BeFalse())
				Expect(d.used(0, 5)).To(BeTrue())
				Expect(d.used(0, 6)).To(BeFalse())
				Expect(d.used(0, 7)).To(BeFalse())

				Expect(d.used(2, 0)).To(BeFalse())
				Expect(d.used(2, 1)).To(BeFalse())
				Expect(d.used(2, 2)).To(BeFalse())
				Expect(d.used(2, 3)).To(BeFalse())
				Expect(d.used(2, 4)).To(BeTrue())
				Expect(d.used(2, 5)).To(BeFalse())
				Expect(d.used(2, 6)).To(BeTrue())
				Expect(d.used(2, 7)).To(BeFalse())

				Expect(d.usedCount()).To(Equal(8108))

				Expect(d.regionCount()).To(Equal(1242))
			})
		})
	})

	Describe("puzzle", func() {
		key := "uugsqrei"

		It("solves star 1", func() {
			d := NewDisk(key)
			c := d.usedCount()
			fmt.Printf("d14 s1: there are %d used blocks\n", c)
		})

		It("solves star 2", func() {
			d := NewDisk(key)
			c := d.regionCount()
			fmt.Printf("d14 s2: there are %d regions\n", c)
		})
	})
})
