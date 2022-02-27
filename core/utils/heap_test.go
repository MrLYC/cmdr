package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core/utils"
)

var _ = Describe("Heap", func() {
	It("should be a sorted heap", func() {
		h := utils.NewSortedHeap(3)
		h.Add("a", 1.0)
		h.Add("c", 3.0)
		h.Add("b", 2.0)

		c, score := h.PopMax()
		Expect(c).To(Equal("c"))
		Expect(score).To(Equal(3.0))

		b, score := h.PopMax()
		Expect(b).To(Equal("b"))
		Expect(score).To(Equal(2.0))

		a, score := h.PopMax()
		Expect(a).To(Equal("a"))
		Expect(score).To(Equal(1.0))
	})
})
