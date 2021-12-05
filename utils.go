package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type Change struct {
	Number   int64
	Field    string
	OldValue string
	NewValue string
}

func ToUnderscore(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetFields(f interface{}) map[string]interface{} {
	inrec, _ := json.Marshal(f)
	var values map[string]interface{}
	json.Unmarshal(inrec, &values)

	return values
}

func ConvertFieldsToUnderscore(items map[string]interface{}) map[string]interface{} {
	values := make(map[string]interface{})

	for key, value := range items {
		values[ToUnderscore(key)] = value
	}

	return values
}

func GetChanges(id int64, items map[string]interface{}, items2 map[string]interface{}) []Change {
	var changes []Change

	for key, value := range items {
		new_value := fmt.Sprintf("%v", value)
		old_value := fmt.Sprintf("%v", items2[key])
		if new_value != old_value {
			changes = append(changes, Change{Number: id, OldValue: old_value, NewValue: new_value, Field: key})
		}
	}

	return changes
}
