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

	output := fmt.Sprintf("type %s struct {\n", structName[0])
	output += convert(result)
	output += "}"
	fmt.Println(output)
}

func convert(jsonObj map[string]interface{}) string {
	output := ""
	for k, v := range jsonObj {
		var value interface{}
		value = v
		valueType := reflect.TypeOf(value).String()
		updatedKey, jsonPart := convertJsonkeyToGoKey(k)
		nested := ""
		nestedValueTypes := ""
		if reflect.TypeOf(v).Kind() == reflect.Map { //if map
			if mapValue, ok := v.(map[string]interface{}); ok {
				nestedValueTypes = "struct {"
				nested = convert(mapValue)

			}
		} else if reflect.TypeOf(v).Kind() == reflect.Slice { //if array
			valueType = "[]interface{}"
		}
		if nestedValueTypes == "" {
			output += fmt.Sprintf("\t%s %s %s\n", updatedKey, valueType, jsonPart)	
		} else {
			output += fmt.Sprintf("\t%s %s\n", updatedKey, nestedValueTypes)
			output += fmt.Sprintf("\t%s", nested)
			output += fmt.Sprintf("\t} %s\n", jsonPart)
		}
		//fmt.Println(nested)
	}
	return output
}

func convertJsonkeyToGoKey(k string) (string, string) {
	snake_case_delim := "_"
	keyParts := strings.Split(k, snake_case_delim)
	for i,part := range keyParts {
		keyParts[i] = strings.Title(part)
	}
	updatedKey := strings.Join(keyParts, "")
	jsonPart := fmt.Sprintf("`json:\"%s\"`", k)
	return updatedKey, jsonPart
}