package adventofcode2017_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAdventofcode2017(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adventofcode2017 Suite")
}
