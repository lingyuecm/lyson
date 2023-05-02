package lyson

import (
	"fmt"
	"strconv"
)

func TransformToString(value any) string {
	switch value.(type) {
	case string:
		return "\"" + value.(string) + "\""
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
	}
	panic("invalid Data Type to Transform to String")
}
