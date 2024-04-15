/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package util

import (
	"errors"
	"regexp"
	"strings"
)

func capitalizeFirstLetter(s string, lowerRest bool) string {
	if len(s) == 0 {
		return s
	}

	rest := s[1:]
	if lowerRest {
		rest = strings.ToLower(rest)
	}

	return strings.ToUpper(s[0:1]) + rest
}

func checkSnakeCase(s string) error {
	pattern := `^([a-zA-Z].*[a-zA-Z0-9])+(_([a-zA-Z].*[a-zA-Z0-9])+)*$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(s) {
		return errors.New("the provided string is not a valid Snake-case style string")
	}

	return nil
}

// This split string into a slice of smaller strings
// by an exception function comparator
// Example:
// - exceptFunc: returns true if the char is alpha numeric
// - source: "hello world, let's make a better-world"
// - return: ["hello", "world", "let", "s", "make", "a", "better", "world"]
func splitStringByExcept(exceptFunc func(rune) bool, source string) []string {
	res := make([]string, 0)
	curr := ""
	for _, c := range source {
		if exceptFunc(c) {
			curr += string(c)
		} else if curr != "" {
			res = append(res, curr)
			curr = ""
		}
	}

	// append the last one
	if curr != "" {
		res = append(res, curr)
	}

	return res
}

func ToSnakeCase(s string) string {
	alphaNumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]$`)
	var splitted = splitStringByExcept(func(r rune) bool {
		return alphaNumericRegex.MatchString(string(r))
	}, s)

	res := ""
	// could have used `strings.Join`, but it need to be lowercase
	for idx, currSplit := range splitted {
		res += strings.ToLower(currSplit)
		if idx < len(splitted)-1 {
			res += "_"
		}
	}

	return res
}

func ToCapitalizedCamelCase(s string) string {
	var splitted []string
	if checkSnakeCase(s) == nil {
		splitted = strings.Split(s, "_")
	} else {
		splitted = strings.Split(s, " ")
	}

	res := ""
	for _, currSplit := range splitted {
		res += capitalizeFirstLetter(currSplit, true)
	}

	return res
}
