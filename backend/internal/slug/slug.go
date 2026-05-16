package slug

import (
	"regexp"
	"strings"
)

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func FromName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonSlug.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "item"
	}
	return s
}
