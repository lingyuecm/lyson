package lyson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func TransformToString(value any) string {
	if nil == value {
		return "null"
	}

	switch value.(type) {
	case string:
		return "\"" + escapeString(value.(string)) + "\""
	case int:
		return strconv.Itoa(value.(int))
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case float64:
		return fmt.Sprintf("%f", value.(float64))
	case bool:
		if value.(bool) {
			return "true"
		} else {
			return "false"
		}
	case *JsonObject:
		return value.(*JsonObject).ToString()
	case *JsonArray:
		return value.(*JsonArray).ToString()
	}
	panic("invalid Data Type to Transform to String")
}

func ParseObject(jsonText string) *JsonObject {
	if regexp.MustCompile("^\\s*$").MatchString(jsonText) {
		return nil
	}
	return parseObject(strings.TrimSpace(jsonText), 0, make(map[string]int))
}

func ParseArray(jsonText string) *JsonArray {
	if regexp.MustCompile("^\\s*$").MatchString(jsonText) {
		return nil
	}
	return parseArray(strings.TrimSpace(jsonText), 0, make(map[string]int))
}
