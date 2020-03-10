package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/xeipuuv/gojsonschema"
)

const resources = "../../../../terraform/resources/azure_policy"

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("policy_validation", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func testJSONPresent() error {
	files := getJSONPolicies()
	if len(files) == 0 {
		return errors.New("there are no JSON files to test")
	}
	return nil
}

func getJSONPolicies() []os.FileInfo {
	f, err := ioutil.ReadDir(resources)
	if err != nil {
		panic("Failed to read or open JSON policy directory")
	}
	return f
}

func testValidJSON() error {
	files := getJSONPolicies()
	for _, f := range files {

		fb, err := ioutil.ReadFile(resources + string(os.PathSeparator) + f.Name())
		if err != nil {
			log.Fatalf("Failed to read JSON file wih name %v", f.Name())
			return err
		}
		var j interface{}

		err = json.NewDecoder(bytes.NewReader(fb)).Decode(&j)
		if err != nil {
			log.Fatalf("Failed to decode JSON from file with name %v", f.Name())
			return err
		}
	}
	return nil
}

func testValidSchemaJSON() error {

	files := getJSONPolicies()
	var success error = nil
	for _, f := range files {

		schemaLoader := gojsonschema.NewReferenceLoader("https://schema.management.azure.com/schemas/2019-06-01/policyDefinition.json")

		fb, err := ioutil.ReadFile(resources + string(os.PathSeparator) + f.Name())
		documentLoader := gojsonschema.NewStringLoader(string(fb))

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			log.Printf("[ERROR] Cannot validate %v due to %v", f.Name(), err)
			return err
		}

		if !result.Valid() {
			success = errors.New("one or more documents failed validation")
			for _, err := range result.Errors() {
				fmt.Printf("Failed to validate %v - %s\n", f.Name(), err)
			}
		}
	}
	return success
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a directory of Azure Policy files in JSON format`, testJSONPresent)
	s.Step(`^the documents must be valid JSON`, testValidJSON)
	s.Step(`^the JSON must be valid against the Microsoft schema`, testValidSchemaJSON)
}
