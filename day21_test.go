package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pixelOn = byte('#')
var pixelOff = byte('.')

type ImageStorageRow []byte
type ImageStorage []ImageStorageRow

func NewImageStorage(size int) ImageStorage {
	rval := make(ImageStorage, size)
	for j := 0; j < size; j++ {
		rval[j] = make(ImageStorageRow, size)
	}
	return rval
}

func imageSize(image string) int {
	return int(math.Sqrt(float64(len(image))))
}

func imageMirrors(image ImageStorage) []ImageStorage {
	size := len(image)
	rval := make([]ImageStorage, 2)
	rval[0] = image

	mirror := NewImageStorage(size)
	for jrow := 0; jrow < size; jrow++ {
		for jcol := 0; jcol < size; jcol++ {
			mirror[jrow][size-1-jcol] = image[jrow][jcol]
		}
	}
	rval[1] = mirror

	return rval
}

func imageRotations(image ImageStorage) []ImageStorage {
	size := len(image)
	rval := make([]ImageStorage, 4)
	rval[0] = image

	for j := 1; j < 4; j++ {
		flip := NewImageStorage(size)
		for jrow := 0; jrow < size; jrow++ {
			for jcol := 0; jcol < size; jcol++ {
				flip[size-1-jcol][jrow] = rval[j-1][jrow][jcol]
			}
		}
		rval[j] = flip
	}

	return rval
}

func imagePermutations(image ImageStorage) []ImageStorage {
	var permutations []ImageStorage
	for _, mirror := range imageMirrors(image) {
		for _, rotation := range imageRotations(mirror) {
			permutations = append(permutations, rotation)
		}
	}
	return permutations
}

func stringImage(storage ImageStorage, newlines bool) string {
	size := len(storage)
	var output []byte
	var index func(row, col int) int

	if newlines {
		output = make([]byte, size*(size+1))
		index = func(row, col int) int {
			return row*(size+1) + col
		}
	} else {
		output = make([]byte, size*size)
		index = func(row, col int) int {
			return row*size + col
		}
	}

	for jrow := 0; jrow < size; jrow++ {
		for jcol := 0; jcol < size; jcol++ {
			output[index(jrow, jcol)] = storage[jrow][jcol]
		}
		if newlines {
			output[index(jrow, size)] = '\n'
		}
	}

	return string(output)
}

func storeImage(image string) ImageStorage {
	image = strings.NewReplacer("\n", "", "/", "").Replace(image)
	size := imageSize(image)
	bareImage := []byte(image)
	storage := NewImageStorage(size)
	for jrow := 0; jrow < size; jrow++ {
		for jcol := 0; jcol < size; jcol++ {
			storage[jrow][jcol] = bareImage[jrow*size+jcol]
		}
	}
	return storage
}

func pluckImage(src ImageStorage, srcRow, srcCol int, size int) ImageStorage {
	rval := NewImageStorage(size)
	copyImage(src, srcRow, srcCol, rval, 0, 0, size)
	return rval
}

func copyImage(src ImageStorage, srcRow, srcCol int, dst ImageStorage, dstRow, dstCol int, size int) {
	for jrow := 0; jrow < size; jrow++ {
		for jcol := 0; jcol < size; jcol++ {
			dst[dstRow+jrow][dstCol+jcol] = src[srcRow+jrow][srcCol+jcol]
		}
	}
}

var initialImage = ".#.\n..#\n###"
var fractalArtRuleRe = regexp.MustCompile(`(.*) => (.*)`)

type FractalArt struct {
	image ImageStorage
	rules map[string]ImageStorage // key stored with no newlines or slashes
}

func NewFractalArt(rules string) *FractalArt {
	fa := FractalArt{image: storeImage(initialImage), rules: make(map[string]ImageStorage)}

	for _, rule := range strings.Split(rules, "\n") {
		if len(rule) == 0 {
			continue
		}

		match := fractalArtRuleRe.FindStringSubmatch(rule)
		if len(match) == 0 {
			panic(fmt.Sprintf("error: could not parse rule `%s`", rule))
		}
		pattern := storeImage(match[1])
		result := storeImage(match[2])
		for _, permutation := range imagePermutations(pattern) {
			fa.rules[stringImage(permutation, false)] = result
		}
	}

	return &fa
}

func (fa *FractalArt) Image() string {
	return stringImage(fa.image, true)
}

func (fa *FractalArt) ZoomAndEnhance() {
	size := len(fa.image)
	var chunkSize, nextChunkSize int

	if size%2 == 0 {
		chunkSize = 2
		nextChunkSize = 3
	} else if size%3 == 0 {
		chunkSize = 3
		nextChunkSize = 4
	} else {
		panic(fmt.Sprintf("error: can't apply rules to image of size %d", size))
	}
	nchunks := size / chunkSize

	nextImage := NewImageStorage(size * nextChunkSize / chunkSize)
	for chunkRow := 0; chunkRow < nchunks; chunkRow++ {
		for chunkCol := 0; chunkCol < nchunks; chunkCol++ {
			stringImage := stringImage(pluckImage(fa.image, chunkRow*chunkSize, chunkCol*chunkSize, chunkSize), false)
			result, ok := fa.rules[stringImage]
			if !ok {
				panic(fmt.Sprintf("error: could not find rule for `%s`", stringImage))
			}
			copyImage(result, 0, 0, nextImage, chunkRow*nextChunkSize, chunkCol*nextChunkSize, nextChunkSize)
		}
	}
	fa.image = nextImage
}

func (fa *FractalArt) PixelCount() int {
	count := 0
	for jrow := 0; jrow < len(fa.image); jrow++ {
		for jcol := 0; jcol < len(fa.image); jcol++ {
			if fa.image[jrow][jcol] == '#' {
				count++
			}
		}
	}
	return count
}

var _ = Describe("Day21", func() {
	Describe("image manipulation", func() {
		Describe("pack/unpack", func() {
			It("are complementary", func() {
				image := "#..#....."
				Expect(stringImage(storeImage(image), false)).To(Equal(image))
			})

			It("unpack can have newlines", func() {
				image := "#..#....."
				imagen := "#..\n#..\n...\n"
				Expect(stringImage(storeImage(image), true)).To(Equal(imagen))
			})

			It("pack generates a 2d byte slice", func() {
				image := "#..#....."
				Expect(storeImage(image)).To(Equal(ImageStorage{
					ImageStorageRow{'#', '.', '.'},
					ImageStorageRow{'#', '.', '.'},
					ImageStorageRow{'.', '.', '.'},
				}))
			})
		})

		Describe("imageMirrors", func() {
			It("returns the image and its mirror", func() {
				image := storeImage("#..#.....")
				permutations := []ImageStorage{
					storeImage("#..#....."),
					storeImage("..#..#..."),
				}
				Expect(imageMirrors(image)).To(ConsistOf(permutations))
			})
		})

		Describe("imageRotations", func() {
			It("returns the image and its rotations", func() {
				image := storeImage("#..#.....")
				permutations := []ImageStorage{
					storeImage("#..#....."),
					storeImage(".##......"),
					storeImage(".....#..#"),
					storeImage("......##."),
				}
				Expect(imageRotations(image)).To(ConsistOf(permutations))
			})
		})

		Describe("imagePermutations", func() {
			It("returns all permutations", func() {
				image := storeImage("#..#.....")
				permutations := []ImageStorage{
					storeImage("#..#....."),
					storeImage(".##......"),
					storeImage(".....#..#"),
					storeImage("......##."),

					storeImage("..#..#..."),
					storeImage(".......##"),
					storeImage("...#..#.."),
					storeImage("##......."),
				}
				Expect(imagePermutations(image)).To(ConsistOf(permutations))
			})
		})
	})

	Describe("FractalArt", func() {
		var fa *FractalArt

		Describe("NewFractalArt()", func() {
			It("sets the art to the starting pattern", func() {
				fa = NewFractalArt("")
				Expect(fa.Image()).To(Equal(".#.\n..#\n###\n"))
			})
		})

		Describe("ZoomAndEnhance()", func() {
			BeforeEach(func() {
				testRules := heredoc.Doc(`
					../.# => ##./#../...
					.#./..#/### => #..#/..../..../#..#
				`)

				fa = NewFractalArt(testRules)
			})

			Context("size 3", func() {
				It("matches 3x3 rules", func() {
					fa.ZoomAndEnhance()
					Expect(fa.Image()).To(Equal("#..#\n....\n....\n#..#\n"))
					Expect(fa.PixelCount()).To(Equal(4))
				})
			})

			Context("size 4", func() {
				It("matches 2x2 rules for all four quadrants", func() {
					fa.ZoomAndEnhance()
					fa.ZoomAndEnhance()
					Expect(fa.Image()).To(Equal("##.##.\n#..#..\n......\n##.##.\n#..#..\n......\n"))
					Expect(fa.PixelCount()).To(Equal(12))
				})
			})
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day21.txt")
		rules := string(rawData)

		It("solves star 1", func() {
			fa := NewFractalArt(rules)
			fa.ZoomAndEnhance()
			fa.ZoomAndEnhance()
			fa.ZoomAndEnhance()
			fa.ZoomAndEnhance()
			fa.ZoomAndEnhance()
			count := fa.PixelCount()
			fmt.Printf("d21 s1: after 5 iterations there are %d pixels on\n", count)
		})
	})
})
