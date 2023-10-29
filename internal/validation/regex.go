package validation

import (
	"regexp"
)

func ValidateWithRegex(input, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}

func FindWithRegex(input, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.FindString(input)
}