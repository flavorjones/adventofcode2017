package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/MakeNowJust/heredoc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ProgramNode struct {
	name     string
	children ProgramNodes
	parent   *ProgramNode
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
			children := matches[3]

			programNode := ProgramNode{name: name}
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

		It("constructs a tree of names", func() {
			root := NewProgramTree(testData)
			Expect(root.name).To(Equal("tknk"))
			Expect(root.children.names()).To(ConsistOf([]string{"ugml", "padx", "fwft"}))
		})
	})

	Describe("puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day7.txt")
		tree_description := string(raw_data)

		It("answers star 1 correctly", func() {
			pt := NewProgramTree(tree_description)
			fmt.Printf("d7 s1: tree root is %s\n", pt.name)
		})
	})
})
