/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package util

import (
	"bytes"
	"log"
	"os"

	"go/format"

	"io/ioutil"
	"text/template"
)

func GenerateTemplate(tplName, tplPath, outputPath string, tplVar interface {}, formatSource bool) error {
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

	// skip source formatting
	if !formatSource {
		return os.WriteFile(outputPath, output.Bytes(), 0644)	
	}

	formattedOutput, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, formattedOutput, 0644)
	if err == nil {
		log.Println("File was generated:", outputPath)
	}

	return err
}