package utils

import "regexp"

func ParseForwardSlash(s string) bool {
	pattern := "/"
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) == 0 {
		return false
	}

	return true
}
