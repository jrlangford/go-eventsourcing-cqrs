// Package jsontools implements functions to manipulate json data.
package jsontools

import "encoding/json"

// JsonToTypeInterface Unmarshals structs into interfaces to a specific type.
func JsonToTypeInterface[T any](jData []byte) (interface{}, error) {
	var object T
	err := json.Unmarshal(jData, &object)
	return &object, err
}
