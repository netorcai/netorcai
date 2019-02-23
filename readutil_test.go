package netorcai

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadIntInString(t *testing.T) {
	_, err := ReadIntInString(nil, "meh", 64, 0, 10)
	assert.Error(t, err, "No error on missing field")

	str := `{"meh":42}`
	var data map[string]interface{}
	json.Unmarshal([]byte(str), &data)

	_, err = ReadIntInString(data, "meh", 64, 0, 10)
	assert.Error(t, err, "No error on non-string value")
}

func TestReadFloatInString(t *testing.T) {
	_, err := ReadFloatInString(nil, "meh", 64, 0, 10)
	assert.Error(t, err, "No error on missing field")

	str := `{"meh":42}`
	var data map[string]interface{}
	json.Unmarshal([]byte(str), &data)

	_, err = ReadFloatInString(data, "meh", 64, 0, 10)
	assert.Error(t, err, "No error on non-string value")
}
