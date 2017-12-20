package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ProgramNode struct {
	name     string
	weight   int
	children ProgramNodes
	parent   *ProgramNode
}

func (pn *ProgramNode) recursiveWeight() int {
	weight := pn.weight
	for _, child := range pn.children {
		weight += child.recursiveWeight()
	}
	return weight
}

func (pn *ProgramNode) weightCheck() (*ProgramNode, int) {
	// terminate on leaf nodes
	if len(pn.children) == 0 {
		return nil, -1
	}

	// depth-first search ...
	for _, child := range pn.children {
		problemNode, rightWeight := child.weightCheck()
		if problemNode != nil {
			return problemNode, rightWeight
		}
	}

	// look at children, bucket recursive weights
	childWeightMap := make(map[int]int) // weight → count
	for _, child := range pn.children {
		weight := child.recursiveWeight()
		_, ok := childWeightMap[weight]
		if ok {
			childWeightMap[weight] += 1
		} else {
			childWeightMap[weight] = 1
		}
	}
	if len(childWeightMap) == 1 {
		return nil, -2
	}

	// if there's an outlier, go through children and find it
	var problemWeight, okWeight int
	for weight, count := range childWeightMap {
		if count == 1 {
			problemWeight = weight
		} else {
			okWeight = weight
		}
	}
	for _, child := range pn.children {
		if child.recursiveWeight() == problemWeight {
			return child, child.weight + okWeight - problemWeight
		}
	}

	return nil, -3
}

type ProgramNodes []*ProgramNode

func (pns ProgramNodes) names() []string {
	rval := make([]string, len(pns))
	for j, programNode := range pns {
		rval[j] = programNode.name
	}
	return rval
}

var programSelfDescriptionRe = regexp.MustCompile(`^(\w+) \((\d+)\)(?: -> (.*))?`)

func NewProgramTree(description string) *ProgramNode {
	lines := strings.Split(description, "\n")

	programMap := make(map[string]*ProgramNode)
	childMap := make(map[string][]string)

	// create nodes for each program, and map child names to parent name
	for _, line := range lines {
		if len(line) > 0 {
			matches := programSelfDescriptionRe.FindStringSubmatch(line)
			name := matches[1]
			weight, _ := strconv.Atoi(matches[2])
			children := matches[3]

			programNode := ProgramNode{name: name, weight: weight}
			programMap[name] = &programNode

			if len(children) > 0 {
				for _, child := range strings.Split(string(children), ", ") {
					childMap[name] = append(childMap[name], child)
				}
			}
		}
	}

	// set up parent/child relationships
	for parentName, parentNode := range programMap {
		childrenNames := childMap[parentName]
		for _, childName := range childrenNames {
			childNode, ok := programMap[childName]
			if !ok {
				panic(fmt.Sprintf("could not find child named %s\n", childName))
			}
			childNode.parent = parentNode
			parentNode.children = append(parentNode.children, childNode)
		}
	}

	// find the root and return it
	for _, program := range programMap {
		if program.parent == nil {
			return program
		}
	}

	return nil
}

var _ = Describe("Day7", func() {
	Describe("ProgramTree", func() {
		testData := heredoc.Doc(`
			pbga (66)
			xhth (57)
			ebii (61)
			havc (66)
			ktlj (57)
			fwft (72) -> ktlj, cntj, xhth
			qoyq (66)
			padx (45) -> pbga, havc, qoyq
			tknk (41) -> ugml, padx, fwft
			jptl (61)
			ugml (68) -> gyxo, ebii, jptl
			gyxo (61)
			cntj (57)
		`)

		var root *ProgramNode

		BeforeEach(func() {
			root = NewProgramTree(testData)
		})

		Describe("NewProgramTree", func() {
			It("constructs a tree of names", func() {
				Expect(root.name).To(Equal("tknk"))
				Expect(root.children.names()).To(ConsistOf([]string{"ugml", "padx", "fwft"}))
			})

			It("stores each program's weight", func() {
				Expect(root.weight).To(Equal(41))
			})
		})

		Describe("recursiveWeight", func() {
			It("adds the weight of itself to the recursive weight of all children", func() {
				weightMap := make(map[string]int) // child name to weight
				for _, child := range root.children {
					weightMap[child.name] = child.recursiveWeight()
				}
				Expect(weightMap["ugml"]).To(Equal(251))
				Expect(weightMap["padx"]).To(Equal(243))
				Expect(weightMap["fwft"]).To(Equal(243))
			})
		})

		Describe("weightCheck", func() {
			It("returns the node that is the wrong weight", func() {
				wrongNode, _ := root.weightCheck()
				Expect(wrongNode.name).To(Equal("ugml"))
			})

			It("returns the weight that the node should be", func() {
				_, rightWeight := root.weightCheck()
				Expect(rightWeight).To(Equal(60))
			})
		})
	})

	Describe("puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day7.txt")
		tree_description := string(raw_data)

		It("answers star 1 correctly", func() {
			pt := NewProgramTree(tree_description)
			fmt.Printf("d7 s1: tree root is %s\n", pt.name)
		})

		It("answers star 2 correctly", func() {
			pt := NewProgramTree(tree_description)
			wrongNode, rightWeight := pt.weightCheck()
			fmt.Printf("d7 s2: wrong node %s, should have weight %d\n", wrongNode.name, rightWeight)
		})
	})
})
