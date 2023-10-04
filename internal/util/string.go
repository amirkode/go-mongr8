/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
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
		return errors.New("The provided string is not a valid Snake-case style string")
	}

	return nil
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