package utils

import "regexp"

func ParseForwardSlash(s string) bool {
	pattern := "/"
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	return len(matches) != 0
}
