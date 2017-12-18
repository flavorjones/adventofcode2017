package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/deckarep/golang-set"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type PassPhrase string

func (p PassPhrase) isValid() bool {
	set := mapset.NewSet()

	for _, word := range strings.Fields(string(p)) {
		if set.Contains(word) {
			return false
		}
		set.Add(word)
	}

	return true
}

var _ = Describe("Day4", func() {
	Describe("PassPhrase", func() {
		It("should only contain unique words", func() {
			Expect(PassPhrase("aa bb cc dd ee").isValid()).To(BeTrue())
			Expect(PassPhrase("aa bb cc dd aa").isValid()).To(BeFalse(), "aa is repeated")
			Expect(PassPhrase("aa bb cc dd aaa").isValid()).To(BeTrue())
		})
	})

	Describe("the puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day4.txt")
		phrases := strings.Split(string(raw_data), "\n")

		It("star 1", func() {
			valid_count := 0
			for _, phrase := range phrases {
				if len(phrase) > 0 {
					if PassPhrase(phrase).isValid() {
						valid_count += 1
					}
				}
			}
			fmt.Printf("d4 s1: there are %d valid phrases\n", valid_count)
		})
	})
})
