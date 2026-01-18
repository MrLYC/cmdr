package strategy

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/mrlyc/cmdr/core"
)

type RewriteStrategy struct {
	config  *StrategyConfig
	tmpl    *template.Template
	enabled bool
}

func (s *RewriteStrategy) Name() string {
	return "rewrite"
}

func (s *RewriteStrategy) Prepare(uri string) (string, error) {
	if s.tmpl == nil {
		return uri, nil
	}

	// Parse to original URI
	var buf bytes.Buffer

	// Template data
	data := struct {
		URI      string
		Scheme   string
		Host     string
		Path     string
		Query    string
		Fragment string
	}{
		URI: uri,
	}

	// Simple URI parsing (for common cases)
	if idx := strings.Index(uri, "://"); idx > 0 {
		data.Scheme = uri[:idx]
		rest := uri[idx+3:]

		if idx2 := strings.Index(rest, "/"); idx2 >= 0 {
			data.Host = rest[:idx2]
			data.Path = rest[idx2:]

			// Split path and query/fragment
			if idx3 := strings.Index(data.Path, "?"); idx3 >= 0 {
				data.Query = data.Path[idx3+1:]
				data.Path = data.Path[:idx3]

				if idx4 := strings.Index(data.Query, "#"); idx4 >= 0 {
					data.Fragment = data.Query[idx4+1:]
					data.Query = data.Query[:idx4]
				}
			} else if idx4 := strings.Index(data.Path, "#"); idx4 >= 0 {
				data.Fragment = data.Path[idx4+1:]
				data.Path = data.Path[:idx4]
			}
		} else {
			data.Host = rest
		}
	}

	// Execute template
	if err := s.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	result := buf.String()

	// If template is empty or same as original, return original
	if result == "" || result == uri {
		return uri, nil
	}

	return result, nil
}

func (s *RewriteStrategy) ShouldRetry(err error) bool {
	// Rewrite strategy doesn't handle retries
	return false
}

func (s *RewriteStrategy) ShouldFallback(err error) bool {
	// Rewrite strategy doesn't trigger fallback
	return false
}

func (s *RewriteStrategy) Configure(cfg core.Configuration) error {
	rewriteRule := cfg.GetString("download.rewrite.rule")
	if rewriteRule == "" {
		s.config = &StrategyConfig{Enabled: false}
		s.enabled = false
		return nil
	}

	s.config = &StrategyConfig{
		Enabled:     true,
		RewriteRule: rewriteRule,
	}

	// Parse condition
	if cfg.IsSet("download.rewrite.condition") {
		s.config.Condition = &StrategyCondition{}

		// Parse schemes
		schemes := cfg.GetStringSlice("download.rewrite.condition.schemes")
		if len(schemes) > 0 {
			s.config.Condition.Schemes = schemes
		}

		// Parse hosts
		hosts := cfg.GetStringSlice("download.rewrite.condition.hosts")
		if len(hosts) > 0 {
			s.config.Condition.Hosts = hosts
		}

		// Parse patterns
		patterns := cfg.GetStringSlice("download.rewrite.condition.patterns")
		if len(patterns) > 0 {
			s.config.Condition.Patterns = patterns
		}

		// If condition is set, strategy is only enabled when matching
		if s.config.Condition.Schemes != nil || s.config.Condition.Hosts != nil || s.config.Condition.Patterns != nil {
			s.config.Enabled = false
		}
	}

	// Parse template
	tmpl, err := template.New("rewrite").Parse(rewriteRule)
	if err != nil {
		return err
	}

	s.tmpl = tmpl
	s.enabled = s.config.Enabled

	return nil
}

func (s *RewriteStrategy) IsEnabled(uri string) bool {
	if s.config == nil {
		return s.enabled
	}
	if s.config.Condition == nil {
		return s.enabled
	}
	return s.config.Matches(uri)
}

func (s *RewriteStrategy) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func (s *RewriteStrategy) IsEnabledConfigured() bool {
	return s.config != nil && s.config.RewriteRule != ""
}

func (s *RewriteStrategy) GetRewrittenURI(uri string) (string, error) {
	return s.Prepare(uri)
}

func NewRewriteStrategy() *RewriteStrategy {
	return &RewriteStrategy{}
}
