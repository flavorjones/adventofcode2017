package adventofcode2017_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/kr/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var routingUp = CartesianCoordinates{0, -1}
var routingRight = CartesianCoordinates{1, 0}
var routingDown = CartesianCoordinates{0, 1}
var routingLeft = CartesianCoordinates{-1, 0}
var routingAdjacentCells = []CartesianCoordinates{
	routingUp, routingRight, routingDown, routingLeft,
}

func routingReverse(direction CartesianCoordinates) CartesianCoordinates {
	switch direction {
	case routingUp:
		return routingDown
	case routingRight:
		return routingLeft
	case routingDown:
		return routingUp
	case routingLeft:
		return routingRight
	default:
		panic(pretty.Sprintf("cannot reverse %v", direction))
	}
}

func isValidPosition(pos CartesianCoordinates) bool {
	return pos.y >= 0 && pos.x >= 0
}

func isAlpha(route byte) bool {
	return 'A' <= route && route <= 'Z'
}

type RoutingTable struct {
	position  CartesianCoordinates
	direction CartesianCoordinates
	table     [][]byte
	letters   []byte
}

func NewRoutingTable(table string) *RoutingTable {
	tableLines := strings.Split(table, "\n")

	byteTable := make([][]byte, len(tableLines))
	for j, line := range tableLines {
		byteTable[j] = []byte(line)
	}

	entryPoint := bytes.IndexByte(byteTable[0], '|')
	position := CartesianCoordinates{x: entryPoint, y: 0}

	return &RoutingTable{table: byteTable, position: position, direction: routingDown}
}

func (r *RoutingTable) sendPacket() {
	byteAt := func(pos CartesianCoordinates) byte {
		if !isValidPosition(pos) {
			return 'x'
		}
		return r.table[pos.y][pos.x]
	}

	for isValidPosition(r.position) {
		route := byteAt(r.position)
		// pretty.Printf("at %v I see `%c` heading %v\n", r.position, route, r.direction)

		switch {
		case route == '|' || route == '-':
			r.position = r.position.move(r.direction)

		case isAlpha(route):
			r.letters = append(r.letters, route)
			r.position = r.position.move(r.direction)

		case route == '+':
			moved := false
			for _, peekDir := range routingAdjacentCells {
				if peekDir == routingReverse(r.direction) {
					continue
				}
				peekPos := r.position.move(peekDir)
				peek := byteAt(peekPos)
				if peek == '|' || peek == '-' || isAlpha(peek) {
					r.position = peekPos
					r.direction = peekDir
					moved = true
					break
				}
			}
			if !moved {
				panic(pretty.Sprintf("error: could not discern move at '%c' with direction %v", route, r.direction))
			}

		case route == ' ':
			return

		default:
			panic(fmt.Sprintf("error: don't recognize '%c' at %v", route, r.position))
		}
		// pretty.Printf("new position is %v\n", r.position)
	}
}

var _ = Describe("Day19", func() {
	Describe("RoutingTable", func() {
		table := heredoc.Doc(`
     |          
     |  +--+    
     A  |  C    
 F---|----E|--+ 
     |  |  |  D 
     +B-+  +--+ 
		`)

		It("does the right thing", func() {
			r := NewRoutingTable(table)
			r.sendPacket()
			Expect(string(r.letters)).To(Equal("ABCDEF"))
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day19.txt")
		table := string(rawData)

		It("solves star 1", func() {
			r := NewRoutingTable(table)
			r.sendPacket()
			fmt.Printf("d19 s1: letters encountered are `%s`\n", string(r.letters))
		})
	})
})
