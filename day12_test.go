package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/deckarep/golang-set"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

type Process struct {
	Pid   string
	Conns []string
}

type PipeMapper struct {
	processes map[string]Process
}

func NewPipeMapper() *PipeMapper {
	return &PipeMapper{processes: make(map[string]Process)}
}

var pidRecordRe = regexp.MustCompile(`(\d+) <-> (.*)`)
var pidSeparatorRe = regexp.MustCompile(`\s*,\s*`)

func (pm *PipeMapper) parseRecord(record string) {
	matches := pidRecordRe.FindStringSubmatch(record)
	pid := matches[1]
	connections := pidSeparatorRe.Split(matches[2], -1)
	process := Process{Pid: pid, Conns: connections}
	pm.processes[pid] = process
}

func (pm *PipeMapper) parseRecords(records string) {
	for _, record := range strings.Split(records, "\n") {
		if len(record) == 0 {
			continue
		}
		pm.parseRecord(record)
	}
}

func (pm *PipeMapper) countPidGroup(pid string) int {
	seen := mapset.NewSet()
	root := pm.processes[pid]

	return pm.recursiveCountPidGroup(&seen, root)
}

func (pm *PipeMapper) recursiveCountPidGroup(seen *mapset.Set, process Process) int {
	if (*seen).Contains(process.Pid) {
		return 0
	}

	count := 1
	(*seen).Add(process.Pid)

	for _, otherPid := range process.Conns {
		count += pm.recursiveCountPidGroup(seen, pm.processes[otherPid])
	}

	return count
}

func (pm *PipeMapper) countGroups() int {
	remaining := mapset.NewSet()
	for pid := range pm.processes {
		remaining.Add(pid)
	}

	count := 0

	for remaining.Cardinality() > 0 {
		// pick one, any one
		remainingIter := remaining.Iterator()
		pid := (<-remainingIter.C).(string)
		remainingIter.Stop()

		root := pm.processes[pid]
		seen := mapset.NewSet()

		pm.recursiveCountPidGroup(&seen, root)

		for pid := range seen.Iter() {
			remaining.Remove(pid)
		}

		count++
	}

	return count
}

var _ = Describe("Day12", func() {
	Describe("PipeMapper", func() {
		testData := heredoc.Doc(`
			0 <-> 2
			1 <-> 1
			2 <-> 0, 3, 4
			3 <-> 2, 4
			4 <-> 2, 3, 6
			5 <-> 6
			6 <-> 4, 5
		`)

		Describe("parseRecord", func() {
			It("parses a record and builds a tree of processes that are connected", func() {
				pm := NewPipeMapper()
				pm.parseRecord(`0 <-> 2`)
				Expect(pm.processes["0"]).To(MatchAllFields(Fields{
					"Pid":   Equal("0"),
					"Conns": ConsistOf([]string{"2"}),
				}))

				pm.parseRecord(`2 <-> 0, 3, 4`)
				Expect(pm.processes["2"]).To(MatchAllFields(Fields{
					"Pid":   Equal("2"),
					"Conns": ConsistOf([]string{"0", "3", "4"}),
				}))
			})
		})

		Describe("parseRecords", func() {
			It("parses multiple records", func() {
				pm := NewPipeMapper()
				pm.parseRecords(testData)

				Expect(pm.processes["0"]).To(MatchAllFields(Fields{
					"Pid":   Equal("0"),
					"Conns": ConsistOf([]string{"2"}),
				}))

				pm.parseRecord(`2 <-> 0, 3, 4`)
				Expect(pm.processes["2"]).To(MatchAllFields(Fields{
					"Pid":   Equal("2"),
					"Conns": ConsistOf([]string{"0", "3", "4"}),
				}))
			})
		})

		Describe("countPidGroup()", func() {
			It("counts the number of processes in the same group as the pid", func() {
				pm := NewPipeMapper()
				pm.parseRecords(testData)

				Expect(pm.countPidGroup("0")).To(Equal(6))
			})
		})

		Describe("countGroups()", func() {
			It("counts the number of independent groups", func() {
				pm := NewPipeMapper()
				pm.parseRecords(testData)

				Expect(pm.countGroups()).To(Equal(2))
			})
		})
	})

	Describe("puzzle", func() {
		rawData, _ := ioutil.ReadFile("day12.txt")
		cookedData := string(rawData)

		It("solves star 1", func() {
			pm := NewPipeMapper()
			pm.parseRecords(cookedData)
			count := pm.countPidGroup("0")
			fmt.Printf("d12 s1: there are %d processes\n", count)
		})

		It("solves star 2", func() {
			pm := NewPipeMapper()
			pm.parseRecords(cookedData)
			count := pm.countGroups()
			fmt.Printf("d12 s2: there are %d process groups\n", count)
		})
	})
})
