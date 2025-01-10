package controller

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVoteService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Voters frontend Suite")
}
