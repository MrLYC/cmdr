package utils_test

import (
	"github.com/mrlyc/cmdr/core/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Replacement", func() {
	DescribeTable("ReplaceString", func(match, template, input, output string, replaced bool) {
		replacement := utils.Replacement{
			Match:    match,
			Template: template,
		}

		realOutput, realReplaced := replacement.ReplaceString(input)
		Expect(realOutput).To(Equal(output))
		Expect(realReplaced).To(Equal(replaced))
	},
		Entry("concat", "", "hello {{ .input }}", "world", "hello world", true),
		Entry("replace", "hi (.*)", "hello {{ index .group 1 }}", "hi world", "hello world", true),
		Entry("urlencode", "", "{{ .input | urlquery }}", " ", "+", true),
		Entry("nothing", "noting", "fail", "ok", "ok", false),
	)
})
