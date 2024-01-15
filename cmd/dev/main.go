package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/kr/pretty"
)

func main() {
	err := testBackslash()
	if err != nil {
		log.Printf("Test Backslash in key: %v", err)
	}
	err = testNewline()
	if err != nil {
		log.Printf("Test newline in key: %v", err)
	}
	err = testQuote()
	if err != nil {
		log.Printf("Test quote in key: %v", err)
	}
	file, err := testLoadYAMLFile()
	if err != nil {
		log.Printf("Test loading YAML file %s: %v", file, err)
	}
	file, err = testLoadJSONFile()
	if err != nil {
		log.Printf("Test loading JSON file %s: %v", file, err)
	}
}

func testNewline() error {
	spec := map[string]string{
		"a":    "a",
		"c\nc": "cc",
		"d":    "d",
	}
	y, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}

	//y, err := os.ReadFile(filename)
	decoder := yaml.NewDecoder(bytes.NewReader(y))
	var out map[string]interface{}
	for i := 0; ; i++ {
		err = decoder.Decode(&out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error on iteration %d: %v", i, err)
		}
		pretty.Print(i, out)
	}
	log.Printf("newline PASSED")
	return nil
}

func testBackslash() error {
	spec := map[string]any{
		"outer": map[string]string{
			"b\\b": "b",
			"d":    "d",
		},
	}
	y, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}

	//y, err := os.ReadFile(filename)
	decoder := yaml.NewDecoder(bytes.NewReader(y))
	var out map[string]interface{}
	for i := 0; ; i++ {
		err = decoder.Decode(&out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error on iteration %d: %v", i, err)
		}
	}
	log.Printf("backslash PASSED")
	return nil
}

func testQuote() error {
	spec := `---
	outer:
    "a\"b\"c": a
    "d\"e\"f": d
`
	decoder := yaml.NewDecoder(strings.NewReader(spec))
	var out map[string]interface{}
	for i := 0; ; i++ {
		err := decoder.Decode(&out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error on iteration %d: %v", i, err)
		}
	}
	log.Printf("quote PASSED")
	return nil
}

/**
bi3/dashboards/4u89zGTVk/dashboard-Vm10mGTVz.yaml <-- plain text/HTML??
bi3/dashboards/7u49jPd7k/dashboard--49ZBkKnk.yaml <-- newline
bi3/dashboards/d40980f9-6366-458b-aac0-2c0a24e9fb6b/dashboard-b4e8bc14-9c16-4181-897c-e30631806d45.yaml <-- newline at beginning of expression

**/

/*
	"bi2/dashboards/WoVNxSA4k/dashboard-ae65738c-671b-4a51-b2e9-49d583b28793.yaml", //backslash
	"bi2/dashboards/WoVNxSA4k/dashboard-cdc9172e-b37e-49fb-8e12-0e9ac806ea44.yaml", //backslash
	"bi2/dashboards/WoVNxSA4k/dashboard-fcce96d1-a898-4726-97ff-59e70c80ad95.yaml", //backslash
	"bi2/dashboards/QqKpDP0Vz/dashboard-e7be048f-625e-486b-a4d3-37a393fe834f.yaml", //backslash
	"bi2/dashboards/FkdQ9Lsnz/dashboard-e2773f54-aebc-407a-8b48-2c7192db2921.yaml", //quote
*/

func testLoadYAMLFile() (string, error) {
	file := "e85464d4-efcd-4351-ad68-3ccc935162fa.yaml"
	log.Printf("Loading %s", file)
	y, err := os.ReadFile(file)
	if err != nil {
		return file, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(y))
	var out map[string]interface{}
	for i := 0; ; i++ {
		err = decoder.Decode(&out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("Error on iteration %d: %v", i, err)
		}
	}
	log.Printf("load %s file PASSED", file)
	return file, nil
}

func testLoadJSONFile() (string, error) {
	spec := map[string]any{}

	file := "e2773f54-aebc-407a-8b48-2c7192db2921.json"
	b, err := os.ReadFile(file)
	if err != nil {
		return file, err
	}
	err = json.Unmarshal(b, &spec)

	y, err := yaml.Marshal(spec)
	if err != nil {
		return file, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(y))
	var out map[string]interface{}
	for i := 0; ; i++ {
		err = decoder.Decode(&out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("Error on iteration %d: %v", i, err)
		}
	}
	log.Printf("load %s file PASSED", file)
	return file, nil
}
