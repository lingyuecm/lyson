package lyson

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func verifyDataType(value any) error {
	var err error = nil
	switch value.(type) {
	case int, int64, float64, bool, string, *JsonObject, *JsonArray:
		break
	default:
		err = errors.New("invalid Data Type")
	}
	return err
}

func getInt(value any) (int, error) {
	var result int
	var err error = nil
	switch value.(type) {
	case int:
		result = value.(int)
	case string:
		// The Assumption Is That the Number Is in Decimal Format
		result, err = strconv.Atoi(value.(string))
	case bool:
		if value.(bool) {
			result = 1
		} else {
			result = 0
		}
	default:
		result = 0
		err = errors.New("cannot Convert to Integer")
	}
	return result, err
}

func getLong(value any) (int64, error) {
	var result int64
	var err error = nil
	switch value.(type) {
	case int:
		result = int64(value.(int))
	case int64:
		result = value.(int64)
	case string:
		result, err = strconv.ParseInt(value.(string), 10, 64)
	case bool:
		if value.(bool) {
			result = 1
		} else {
			result = 0
		}
	default:
		result = 0
		err = errors.New("cannot Convert to Integer64")
	}
	return result, err
}

func getDouble(value any) (float64, error) {
	var result float64
	var err error = nil
	switch value.(type) {
	case float64:
		result = value.(float64)
	case int:
		result = float64(value.(int))
	case int64:
		result = float64(value.(int64))
	case string:
		result, err = strconv.ParseFloat(value.(string), 64)
	default:
		result = 0
		err = errors.New("cannot Convert to Float64")
	}
	return result, err
}

func getBoolean(value any) (bool, error) {
	var result bool
	var err error = nil
	switch value.(type) {
	case bool:
		result = value.(bool)
	case int, int64:
		if 0 == value {
			result = false
		} else {
			result = true
		}
	case string:
		if "true" == strings.ToLower(value.(string)) {
			result = true
		} else if "false" == strings.ToLower(value.(string)) {
			result = false
		} else {
			result = false
			err = errors.New("invalid Boolean Text")
		}
	default:
		result = false
		err = errors.New("cannot Convert to Boolean")
	}
	return result, err
}

func getString(value any) (string, error) {
	var result string
	var err error = nil
	switch value.(type) {
	case int, int64:
		result = strconv.FormatInt(value.(int64), 10)
	case float64:
		result = fmt.Sprintf("%f", value.(float64))
	case bool:
		result = strconv.FormatBool(value.(bool))
	case string:
		result = value.(string)
	case *JsonObject:
		result = "JsonObject"
	case *JsonArray:
		result = "JsonArray"
	default:
		result = ""
		err = errors.New("cannot Convert to String")
	}
	return result, err
}

func getJsonObject(value any) (*JsonObject, error) {
	var result *JsonObject
	var err error = nil
	switch value.(type) {
	case *JsonObject:
		result = value.(*JsonObject)
	default:
		result = nil
		err = errors.New("cannot Convert to JSON Object")
	}
	return result, err
}

func getJsonArray(value any) (*JsonArray, error) {
	var result *JsonArray
	var err error = nil
	switch value.(type) {
	case *JsonArray:
		result = value.(*JsonArray)
	default:
		result = nil
		err = errors.New("cannot Convert to JSON Array")
	}
	return result, err
}
