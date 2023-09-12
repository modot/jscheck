package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <json_file_name> <schema_file_name>")
		return
	}

	jsonFileName := os.Args[1]
	schemaFileName := os.Args[2]

	err := verifyJSONWithSchema(jsonFileName, schemaFileName)
	if err != nil {
		panic(err)
	}
}

func verifyJSONWithSchema(jsonFilePath string, schemaFilePath string) error {
	// Load JSON file
	jsonFile, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	// Load JSON schema file
	schemaFile, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON schema file: %v", err)
	}

	// Parse JSON
	var jsonData interface{}
	err = json.Unmarshal(jsonFile, &jsonData)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Parse JSON schema
	schemaLoader := gojsonschema.NewStringLoader(string(schemaFile))
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("failed to parse JSON schema: %v", err)
	}

	// Validate JSON against schema
	documentLoader := gojsonschema.NewGoLoader(jsonData)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("failed to validate JSON against schema: %v", err)
	}

	// Check validation result
	if !result.Valid() {
		errMsg := "JSON is not valid according to the schema:\n"
		for _, desc := range result.Errors() {
			errMsg += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}
