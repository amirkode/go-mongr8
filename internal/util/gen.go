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

	"go/format"

	"io/ioutil"
	"text/template"
)

func GenerateTemplate(tplName, tplPath, outputPath string, tplVar interface {}) error {
	fmt.Println("file path: " + tplPath)
	data, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	t, err := template.New(tplName).Parse(string(data))
	if err != nil {
		return err
	}
	
	output := &bytes.Buffer{}
	err = t.Execute(output, tplVar)
	if err != nil {
		return err
	}

	// fmt.Printf("%s", output.Bytes())
	// fmt.Println("halo")

	formattedOutput, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, formattedOutput, 0644)
}