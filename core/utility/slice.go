package utility

import "strings"

func ContainsString(needle string, haystack []string) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}

func ContainsStringStartingWith(prefix string, a []string) bool {
	for _, element := range a {
		if strings.HasPrefix(element, prefix) {
			return true
		}
	}
	return false
}