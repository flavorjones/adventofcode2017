package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type StreamProcessor struct {
	stream []byte
}

func NewStreamProcessor(stream string) *StreamProcessor {
	return &StreamProcessor{stream: []byte(stream)}
}

func (sp *StreamProcessor) score() int {
	score, _, _ := parseGroup(sp.stream, 1)
	return score
}

func (sp *StreamProcessor) garbage() int {
	_, _, nGarbage := parseGroup(sp.stream, 1)
	return nGarbage
}

// parseGroup returns (score, length, nGarbageChars)
func parseGroup(stream []byte, depth int) (int, int, int) {
	// pretty.Println("parseGroup: ", depth, string(stream))
	if stream[0] != '{' {
		panic("group doesn't start with `{`")
	}
	nGarbage := 0
	score := depth
	for jbyte := 1; jbyte < len(stream); {
		switch stream[jbyte] {
		case '}':
			return score, jbyte, nGarbage
		case '{':
			childScore, childLen, nGarbageChars := parseGroup(stream[jbyte:], depth+1)
			score += childScore
			nGarbage += nGarbageChars
			jbyte += childLen + 1
		case '<':
			garbageLen, nGarbageChars := parseGarbage(stream[jbyte:])
			nGarbage += nGarbageChars
			jbyte += garbageLen + 1
		default:
			jbyte++
		}
	}
	return -1, -1, -1
}

// parseGarbage returns length, nGarbageChars
func parseGarbage(stream []byte) (int, int) {
	// pretty.Println("parseGarbage: ", string(stream))
	if stream[0] != '<' {
		panic("garbage doesn't start with `{`")
	}
	ngarbage := 0
	for jbyte := 1; jbyte < len(stream); {
		switch stream[jbyte] {
		case '>':
			return jbyte, ngarbage
		case '!':
			jbyte += 2
		default:
			ngarbage++
			jbyte++
		}
	}
	return -1, -1
}

var _ = Describe("Day9", func() {
	Describe("StreamProcessor", func() {
		Describe("score()", func() {
			It("calculates a score for bare groups", func() {
				Expect(NewStreamProcessor(`{}`).score()).To(Equal(1))
				Expect(NewStreamProcessor(`{{{}}}`).score()).To(Equal(6))
				Expect(NewStreamProcessor(`{{},{}}`).score()).To(Equal(5))
				Expect(NewStreamProcessor(`{{{},{},{{}}}}`).score()).To(Equal(16))
			})

			It("parses garbage correctly", func() {
				Expect(NewStreamProcessor(`{<a>,<a>,<a>,<a>}`).score()).To(Equal(1))
				Expect(NewStreamProcessor(`{{<ab>},{<ab>},{<ab>},{<ab>}}`).score()).To(Equal(9))
				Expect(NewStreamProcessor(`{{<!!>},{<!!>},{<!!>},{<!!>}}`).score()).To(Equal(9))
				Expect(NewStreamProcessor(`{{<a!>},{<a!>},{<a!>},{<ab>}}`).score()).To(Equal(3))
			})
		})

		Describe("garbage()", func() {
			It("counts garbage characters", func() {
				Expect(NewStreamProcessor(`{<>}`).garbage()).To(Equal(0))
				Expect(NewStreamProcessor(`{<random characters>}`).garbage()).To(Equal(17))
				Expect(NewStreamProcessor(`{<<<<>}`).garbage()).To(Equal(3))
			})

			It("doesn't count cancelled characters or `!`", func() {
				Expect(NewStreamProcessor(`{<{!>}>}`).garbage()).To(Equal(2))
				Expect(NewStreamProcessor(`{<!!>}`).garbage()).To(Equal(0))
				Expect(NewStreamProcessor(`{<!!!>>}`).garbage()).To(Equal(0))
				Expect(NewStreamProcessor(`{<{o"i!a,<{i<a>}`).garbage()).To(Equal(10))
			})
		})
	})

	Describe("puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day9.txt")
		stream := string(raw_data)

		It("answers star 1 correctly", func() {
			score := NewStreamProcessor(stream).score()
			fmt.Printf("d9 s1: stream score is %d\n", score)
		})

		It("answers star 2 correctly", func() {
			garbage := NewStreamProcessor(stream).garbage()
			fmt.Printf("d9 s2: stream garbage had %d chars\n", garbage)
		})
	})
})
