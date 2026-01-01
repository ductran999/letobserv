package xstrings

import "strings"

func IsEmpty(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) == 0
}
