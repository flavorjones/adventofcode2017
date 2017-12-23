package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	scannerDirectionDown = false
	scannerDirectionUp   = true
)

type scannerDirection bool

type ScannerState struct {
	srange    int
	position  int
	direction scannerDirection
}

func (ss *ScannerState) tock() {
	if ss.direction == scannerDirectionDown {
		ss.position += 1
		if ss.position >= ss.srange {
			ss.position -= 2
			ss.direction = !ss.direction
		}
	} else {
		ss.position -= 1
		if ss.position < 0 {
			ss.position += 2
			ss.direction = !ss.direction
		}
	}
}

type Trip struct {
	packetPos     int
	scannerStates []*ScannerState
	severity      int
}

type ScannersDescriptor map[int]int

func (s *ScannersDescriptor) maxDepth() int {
	max := -1
	for depth := range *s {
		if depth > max {
			max = depth
		}
	}
	return max
}

func NewTrip(f *Firewall) *Trip {
	scannerStates := make([]*ScannerState, f.scannersDescriptor.maxDepth()+1)
	for sd, sr := range f.scannersDescriptor {
		scannerStates[sd] = &ScannerState{srange: sr, position: 0}
	}
	return &Trip{packetPos: -1, scannerStates: scannerStates}
}

func (t *Trip) tick() {
	t.packetPos += 1
	scanner := t.scannerStates[t.packetPos]
	if scanner == nil {
		return
	}
	if scanner.position == 0 {
		t.severity += scanner.srange * t.packetPos
	}
}

func (t *Trip) tock() {
	for _, jscanner := range t.scannerStates {
		if jscanner == nil {
			continue
		}
		jscanner.tock()
	}
}

type Firewall struct {
	scannersDescriptor ScannersDescriptor
}

func NewFirewall(scannersDesc string) *Firewall {
	scannersDescriptor := make(ScannersDescriptor)

	for _, s := range strings.Split(scannersDesc, "\n") {
		if len(s) == 0 {
			continue
		}

		parsed := strings.Split(s, ":")
		sDepth, _ := strconv.Atoi(strings.TrimSpace(parsed[0]))
		sRange, _ := strconv.Atoi(strings.TrimSpace(parsed[1]))
		scannersDescriptor[sDepth] = sRange
	}

	return &Firewall{scannersDescriptor: scannersDescriptor}
}

func (f *Firewall) tripSeverity() int {
	trip := NewTrip(f)
	for trip.packetPos < len(trip.scannerStates)-1 {
		trip.tick()
		trip.tock()
	}
	return trip.severity
}

var _ = Describe("Day13", func() {
	Describe("Firewall", func() {
		testInput := heredoc.Doc(`
			0: 3
			1: 2
			4: 4
			6: 4
		`)

		Describe("NewFirewall", func() {
			It("takes a scanner description and builds a data structure", func() {
				f := NewFirewall(testInput)
				Expect(len(f.scannersDescriptor)).To(Equal(4))
			})
		})

		Describe("tripSeverity()", func() {
			It("takes returns the calculated severity of a trip that starts at t=0", func() {
				f := NewFirewall(testInput)
				Expect(f.tripSeverity()).To(Equal(24))
			})
		})

		Describe("ScannerDescriptor", func() {
			Describe("maxDepth()", func() {
				It("returns the max depth of the set of scanners", func() {
					f := NewFirewall(testInput)
					Expect(f.scannersDescriptor.maxDepth()).To(Equal(6))
				})
			})
		})

		Describe("Trip", func() {
			Describe("tick()", func() {
				It("advances the packet and checks if scanner caught us", func() {
					f := NewFirewall(testInput)
					t := NewTrip(f)

					t.tick()
					Expect(t.packetPos).To(Equal(0))
					Expect(t.scannerStates[0].position).To(Equal(0))
					Expect(t.severity).To(Equal(0)) // we were caught, but severity was 0

					t.tick()
					Expect(t.packetPos).To(Equal(1))
					Expect(t.scannerStates[1].position).To(Equal(0))
					Expect(t.severity).To(Equal(2)) // caught by scanner 1 depth 2 → 2

					t.tick()
					Expect(t.packetPos).To(Equal(2))

					t.tick()
					Expect(t.packetPos).To(Equal(3))

					t.tick()
					Expect(t.packetPos).To(Equal(4))
					Expect(t.scannerStates[4].position).To(Equal(0))
					Expect(t.severity).To(Equal(18)) // caught by scanner 4 depth 4 → 16
				})
			})

			Describe("tock()", func() {
				It("advances each of the scanners", func() {
					f := NewFirewall(testInput)
					t := NewTrip(f)
					Expect(t.scannerStates[0].position).To(Equal(0))
					Expect(t.scannerStates[1].position).To(Equal(0))
					Expect(t.scannerStates[4].position).To(Equal(0))
					Expect(t.scannerStates[6].position).To(Equal(0))

					t.tock()
					Expect(t.scannerStates[0].position).To(Equal(1))
					Expect(t.scannerStates[1].position).To(Equal(1))
					Expect(t.scannerStates[4].position).To(Equal(1))
					Expect(t.scannerStates[6].position).To(Equal(1))

					t.tock()
					Expect(t.scannerStates[0].position).To(Equal(2))
					Expect(t.scannerStates[1].position).To(Equal(0))
					Expect(t.scannerStates[4].position).To(Equal(2))
					Expect(t.scannerStates[6].position).To(Equal(2))

					t.tock()
					Expect(t.scannerStates[0].position).To(Equal(1))
					Expect(t.scannerStates[1].position).To(Equal(1))
					Expect(t.scannerStates[4].position).To(Equal(3))
					Expect(t.scannerStates[6].position).To(Equal(3))

					t.tock()
					Expect(t.scannerStates[0].position).To(Equal(0))
					Expect(t.scannerStates[1].position).To(Equal(0))
					Expect(t.scannerStates[4].position).To(Equal(2))
					Expect(t.scannerStates[6].position).To(Equal(2))
				})
			})
		})
	})

	Describe("puzzle", func() {
		It("solves star 1", func() {
			rawData, _ := ioutil.ReadFile("day13.txt")
			f := NewFirewall(string(rawData))
			sev := f.tripSeverity()
			fmt.Printf("d13 s1: trip severity is %d\n", sev)
		})
	})
})
