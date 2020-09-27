package gopherhole

import (
	"mime"
	"path/filepath"
	"strings"
)

type mimeType string

func (m *mimeType) guess(file string) *mimeType {
	ext := filepath.Ext(file)
	name := mime.TypeByExtension(ext)
	name, _, _ = mime.ParseMediaType(name)
	*m = mimeType(name)
	return m
}

func (m *mimeType) string() string {
	if m == nil {
		return ""
	}

	return string(*m)
}

func (m *mimeType) isNone() bool {
	return m.string() == ""
}

func (m *mimeType) in(matchers []string) bool {
	if m.isNone() || len(matchers) == 0 {
		return false
	}

	for _, matcher := range matchers {
		if m.string() == matcher ||
			(strings.HasSuffix(matcher, "/") &&
				strings.Index(m.string(), matcher) == 0) {
			return true
		}
	}

	return false
}
