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
	score, _ := parseGroup(sp.stream, 1)
	return score
}

// parseGroup returns (score, length)
func parseGroup(stream []byte, depth int) (int, int) {
	// pretty.Println("parseGroup: ", depth, string(stream))
	if stream[0] != '{' {
		panic("group doesn't start with `{`")
	}
	score := depth
	for jbyte := 1; jbyte < len(stream); {
		switch stream[jbyte] {
		case '}':
			return score, jbyte
		case '{':
			childScore, childLen := parseGroup(stream[jbyte:], depth+1)
			score += childScore
			jbyte += childLen + 1
		case '<':
			garbageLen := parseGarbage(stream[jbyte:])
			jbyte += garbageLen + 1
		default:
			jbyte++
		}
	}
	return -1, -1
}

func parseGarbage(stream []byte) int {
	// pretty.Println("parseGarbage: ", string(stream))
	if stream[0] != '<' {
		panic("garbage doesn't start with `{`")
	}
	for jbyte := 1; jbyte < len(stream); {
		switch stream[jbyte] {
		case '>':
			return jbyte
		case '!':
			jbyte += 2
		default:
			jbyte++
		}
	}
	return -1
}

var _ = Describe("Day9", func() {
	Describe("StreamProcessor", func() {
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

	Describe("puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day9.txt")
		stream := string(raw_data)

		It("answers star 1 correctly", func() {
			score := NewStreamProcessor(stream).score()
			fmt.Printf("d9 s1: stream score is %d\n", score)
		})
	})
})
