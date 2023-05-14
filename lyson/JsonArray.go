package lyson

import "bytes"

// JsonArray Is the Implementation of JSONArray
type JsonArray struct {
	// elements Are the Elements Wrapped in the JSON Array
	elements []any
}

func NewArray() *JsonArray {
	result := new(JsonArray)
	result.elements = make([]any, 0, 10)
	return result
}

// AddElement Adds an Element to the End of the Array
func (arr *JsonArray) AddElement(element any) error {
	err := verifyDataType(element)
	if nil == err {
		arr.elements = append(arr.elements, element)
	}
	return err
}

func (arr *JsonArray) GetObject(index int) any {
	return arr.elements[index]
}

func (arr *JsonArray) GetInt(index int) (int, error) {
	value := arr.elements[index]
	return getInt(value)
}

func (arr *JsonArray) GetLong(index int) (int64, error) {
	value := arr.elements[index]
	return getLong(value)
}

func (arr *JsonArray) GetDouble(index int) (float64, error) {
	value := arr.elements[index]
	return getDouble(value)
}

func (arr *JsonArray) GetBoolean(index int) (bool, error) {
	value := arr.elements[index]
	return getBoolean(value)
}

func (arr *JsonArray) GetString(index int) (string, error) {
	value := arr.elements[index]
	return getString(value)
}

func (arr *JsonArray) GetJsonObject(index int) (*JsonObject, error) {
	value := arr.elements[index]
	return getJsonObject(value)
}

func (arr *JsonArray) GetJsonArray(index int) (*JsonArray, error) {
	value := arr.elements[index]
	return getJsonArray(value)
}

func (arr *JsonArray) ToString() string {
	if 0 == len(arr.elements) {
		return "[]"
	}
	var buf bytes.Buffer
	var valueText string
	for _, e := range arr.elements {
		buf.WriteString(",")
		valueText = TransformToString(e)
		buf.WriteString(valueText)
	}
	result := buf.String()
	return "[" + result[1:] + "]"
}
