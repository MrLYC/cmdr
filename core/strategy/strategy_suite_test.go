package strategy

import (
	"testing"

	. "github.com/onsi/ginkgo"
)

func TestStrategy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Strategy Suite")
}
