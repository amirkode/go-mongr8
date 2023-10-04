/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package util

import (
	"bytes"
	"fmt"
	"os"

	"io/ioutil"
	"text/template"
)

func GenerateTemplate(tplPath, outputPath string, tplVar interface {}) error {
	fmt.Println("file path: " + tplPath)
	data, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		return err
	}
	
	output := &bytes.Buffer{}
	err = t.Execute(output, tplVar)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, output.Bytes(), 0644)
}