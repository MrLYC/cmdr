package utils

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/mrlyc/cmdr/core"
)

type Replacement struct {
	Match    string
	Template string
}

func (r *Replacement) ReplaceString(s string) (string, bool) {
	logger := core.GetLogger()

	regex := regexp.MustCompile(r.Match)
	group := regex.FindStringSubmatch(s)
	if group == nil {
		return s, false
	}

	tmpl := template.Must(template.New("").Parse(r.Template))
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, map[string]interface{}{
		"input": s,
		"group": group,
	})
	if err != nil {
		return s, false
	}

	replaced := buf.String()
	logger.Debug("replaced", map[string]interface{}{
		"location": s,
		"match":    r.Match,
		"replaced": replaced,
	})

	return replaced, true
}

type Replacements []*Replacement

func (r Replacements) ReplaceString(s string) (string, bool) {
	for _, r := range r {
		s, ok := r.ReplaceString(s)
		if ok {
			return s, ok
		}
	}

	return s, false
}
