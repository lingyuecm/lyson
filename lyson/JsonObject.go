package lyson

import (
	"bytes"
)

// JsonObject Is the Implementation of JSONObject
type JsonObject struct {
	// entries Are the Entries Wrapped in the JSON Object
	entries map[string]any
}

func NewObject() *JsonObject {
	result := new(JsonObject)
	result.entries = make(map[string]any)
	return result
}

// PutEntry Puts a Value at a Key
func (obj *JsonObject) PutEntry(key string, value any) error {
	err := verifyDataType(value)
	if nil == err {
		obj.entries[key] = value
	}
	return err
}

func (obj *JsonObject) GetObject(key string) (any, bool) {
	value, ok := obj.entries[key]
	return value, ok
}

func (obj *JsonObject) GetInt(key string) (int, bool, error) {
	value, ok := obj.entries[key]
	if !ok {
		return 0, ok, nil
	}
	result, err := getInt(value)
	return result, ok, err
}

func (obj *JsonObject) GetLong(key string) (int64, bool, error) {
	value, ok := obj.entries[key]
	if !ok {
		return 0, ok, nil
	}
	result, err := getLong(value)
	return result, ok, err
}

func (obj *JsonObject) GetDouble(key string) (float64, bool, error) {
	value, ok := obj.entries[key]
	if !ok {
		return 0, ok, nil
	}
	result, err := getDouble(value)
	return result, ok, err
}

func (obj *JsonObject) GetBoolean(key string) (bool, bool, error) {
	value, ok := obj.entries[key]
	if !ok {
		return false, ok, nil
	}
	result, err := getBoolean(value)
	return result, ok, err
}

func (obj *JsonObject) GetString(key string) (string, bool, error) {
	value, ok := obj.entries[key]
	if !ok {
		return "", ok, nil
	}
	result, err := getString(value)
	return result, ok, err
}

func (obj *JsonObject) GetJsonObject(key string) (*JsonObject, bool, error) {
	value, ok := obj.entries[key]
	if !ok {
		return nil, ok, nil
	}
	result, err := getJsonObject(value)
	return result, ok, err
}

func (obj *JsonObject) GetJsonArray(key string) (*JsonArray, bool, error) {
	value, ok := obj.entries[key]
	if !ok {
		return nil, ok, nil
	}
	result, err := getJsonArray(value)
	return result, ok, err
}

func (obj *JsonObject) ToString() string {
	if 0 == len(obj.entries) {
		return "{}"
	}
	var buf bytes.Buffer
	var valueText string
	for key := range obj.entries {
		buf.WriteString(",\"")
		buf.WriteString(key)
		buf.WriteString("\":")

		valueText = TransformToString(obj.entries[key])
		buf.WriteString(valueText)
	}
	result := buf.String()
	return "{" + result[1:] + "}"
}
