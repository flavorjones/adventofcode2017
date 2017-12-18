package adventofcode2017_test

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/deckarep/golang-set"
	"github.com/fighterlyt/permutation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func permutations(word string) []string {
	bword := []byte(word)
	rval := []string{}

	permutator, _ := permutation.NewPerm(bword, nil)
	for permutation, err := permutator.Next(); err == nil; permutation, err = permutator.Next() {
		rval = append(rval, string(permutation.([]byte)))
	}

	return rval
}

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

func (p PassPhrase) isValid2() bool {
	set := mapset.NewSet()

	for _, word := range strings.Fields(string(p)) {
		if set.Contains(word) {
			return false
		}
		set.Add(word)
		for _, permutation := range permutations(word) {
			set.Add(permutation)
		}
	}

	return true
}

var _ = Describe("Day4", func() {
	Describe("permutations", func() {
		It("should return all valid permutations of a word", func() {
			Expect(permutations("a")).To(ConsistOf("a"))
			Expect(permutations("ab")).To(ConsistOf("ab", "ba"))
			Expect(permutations("abc")).To(ConsistOf("abc", "acb", "bac", "bca", "cab", "cba"))
		})
	})

	Describe("PassPhrase", func() {
		Describe("isValid()", func() {
			It("should be true if the phrase contains only unique words", func() {
				Expect(PassPhrase("aa bb cc dd ee").isValid()).To(BeTrue())
				Expect(PassPhrase("aa bb cc dd aa").isValid()).To(BeFalse(), "aa is repeated")
				Expect(PassPhrase("aa bb cc dd aaa").isValid()).To(BeTrue())
			})
		})

		Describe("isValid2()", func() {
			It("should be true if the phrase contains only unique words and no anagrams", func() {
				Expect(PassPhrase("aa bb cc dd ee").isValid2()).To(BeTrue())
				Expect(PassPhrase("aa bb cc dd aa").isValid2()).To(BeFalse(), "aa is repeated")
				Expect(PassPhrase("aa bb cc dd aaa").isValid2()).To(BeTrue())

				Expect(PassPhrase("abcde fghij").isValid2()).To(BeTrue())
				Expect(PassPhrase("abcde xyz ecdab").isValid2()).To(BeFalse())
				Expect(PassPhrase("a ab abc abd abf abj").isValid2()).To(BeTrue())
				Expect(PassPhrase("iiii oiii ooii oooi oooo").isValid2()).To(BeTrue())
				Expect(PassPhrase("oiii ioii iioi iiio").isValid2()).To(BeFalse())
			})
		})
	})

	Describe("the puzzle", func() {
		raw_data, _ := ioutil.ReadFile("day4.txt")
		phrases := strings.Split(string(raw_data), "\n")

		It("star 1", func() {
			valid_count := 0
			for _, phrase := range phrases {
				if len(phrase) > 0 && PassPhrase(phrase).isValid() {
					valid_count += 1
				}
			}
			fmt.Printf("d4 s1: there are %d valid phrases\n", valid_count)
		})

		It("star 2", func() {
			valid_count := 0
			for _, phrase := range phrases {
				if len(phrase) > 0 && PassPhrase(phrase).isValid2() {
					valid_count += 1
				}
			}
			fmt.Printf("d4 s2: there are %d valid phrases\n", valid_count)
		})
	})
})
