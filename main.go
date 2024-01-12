package main

import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"strings"
	"reflect"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go filename.json")
		return
	}

	fileName := os.Args[1]
	structName := strings.Split(fileName, ".")

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(byteValue, &result); err != nil {
		fmt.Println("Error Decoding JSON.", err)
		return
	}

	snake_case_delim := "_"
	//space_delim := " "
	output := fmt.Sprintf("type %s struct {\n", structName[0])
	for k,v := range result {
		var value interface{}
		keyParts := strings.Split(k, snake_case_delim)
		for i, part := range keyParts {
			keyParts[i] = strings.Title(part)
		}
		finalKey := strings.Join(keyParts, "")
		value = v
		value = reflect.TypeOf(value)
		output += fmt.Sprintf("\t%s %s\n", finalKey, value)
	}
	output += "}"
	fmt.Println(output)
}




