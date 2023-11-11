/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
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