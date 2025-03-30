package service

import "regexp"

func isInvalidName(name string) bool {
	return !regexp.MustCompile(`^[ \w+]{1,128}$`).MatchString(name) ||
		regexp.MustCompile("  ").MatchString(name) || name[0] == ' ' || name[len(name)-1] == ' '
}
