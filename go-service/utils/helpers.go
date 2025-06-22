package utils

import (
	"encoding/json"
	"math/rand"
	"time"
)

// RandomInt generates a random integer between min and max (inclusive)
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// ToJSON converts an interface to JSON string
func ToJSON(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

// FromJSON converts JSON string to interface
func FromJSON(jsonStr string, target interface{}) error {
	return json.Unmarshal([]byte(jsonStr), target)
}
