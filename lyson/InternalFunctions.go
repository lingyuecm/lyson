package lyson

import (
	"errors"
	"fmt"
	"lyson/constant"
	"regexp"
	"strconv"
	"strings"
)

const endIndexKey = "endIndex"

func verifyDataType(value any) error {
	if nil == value {
		return nil
	}
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
	case int64:
		result = int(value.(int64))
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

func parseObject(jsonText string, startIndex int, endIndex map[string]int) *JsonObject {
	length := len(jsonText)

	bStack := NewStack[byte]()

	joStack := NewStack[*JsonObject]()
	rootObject := NewObject()
	joStack.Push(rootObject)

	var currentByte byte
	var peekByte byte
	key := ""

	tokenBytes := make([]byte, 0, 10)
	unicodeBytes := make([]byte, 0, 4)

	for m := startIndex; m < length; m++ {
		currentByte = jsonText[m]
		switch currentByte {
		case constant.StartObject:
			if bStack.Size() == 0 {
				// Start Parsing
				bStack.Push(currentByte)

				result := NewObject()
				_ = joStack.Peek().PutEntry("result", result)
				joStack.Push(result)
				break
			}
			peekByte = bStack.Peek()
			if constant.Quote == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
				break
			}
			if constant.Colon == peekByte {
				bStack.Push(currentByte)

				child := NewObject()
				_ = joStack.Peek().PutEntry(key, child)
				joStack.Push(child)
			}
			key = ""
		case constant.StartArray:
			peekByte = bStack.Peek()
			if constant.Quote == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
				break
			}
			if "" != key {
				_ = joStack.Peek().PutEntry(key, parseArray(jsonText, m, endIndex))
				m = endIndex[endIndexKey]
				key = ""
			}
		case constant.Quote:
			peekByte = bStack.Peek()
			if constant.BackSlash == peekByte {
				tokenBytes = append(tokenBytes, findEscape(currentByte))
				_ = bStack.Pop()
				break
			}
			if constant.StartObject == peekByte {
				bStack.Push(currentByte)
				tokenBytes = make([]byte, 0, 10)
			} else if constant.Colon == peekByte {
				bStack.Push(currentByte)
				tokenBytes = append(make([]byte, 0, 10), currentByte)
			} else if constant.Quote == peekByte {
				_ = bStack.Pop()
				peekByte = bStack.Peek()
				if constant.StartObject == peekByte {
					key = string(tokenBytes)
				} else if constant.Colon == peekByte {
					tokenBytes = append(tokenBytes, currentByte)
				}
			}
		case constant.Colon:
			peekByte = bStack.Peek()
			if constant.Quote == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
				break
			}
			bStack.Push(currentByte)
			tokenBytes = make([]byte, 0, 10)
		case constant.Comma:
			peekByte = bStack.Peek()
			if constant.Quote == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
				break
			}
			_ = bStack.Pop()
			if "" != key {
				_ = joStack.Peek().PutEntry(key, translateValue(strings.TrimSpace(string(tokenBytes))))
			}
		case constant.EndObject:
			peekByte = bStack.Peek()
			if constant.Quote == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
				break
			}
			peekByte = bStack.Pop()
			if constant.Colon == peekByte {
				_ = bStack.Pop()
			}
			if "" != key {
				_ = joStack.Peek().PutEntry(key, translateValue(strings.TrimSpace(string(tokenBytes))))
			}

			if 0 == bStack.Size() {
				endIndex[endIndexKey] = m
				return joStack.Pop()
			}
			rootObject = joStack.Pop()
			key = ""
		case constant.BackSlash:
			peekByte = bStack.Peek()
			if constant.Quote == peekByte {
				bStack.Push(currentByte)
			} else if constant.BackSlash == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
				_ = bStack.Pop()
			}
		case constant.LetterU:
			peekByte = bStack.Peek()
			if constant.BackSlash == peekByte {
				bStack.Push(currentByte)
				unicodeBytes = make([]byte, 0, 4)
			} else if constant.Quote == peekByte || constant.Colon == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
			}
		default:
			peekByte = bStack.Peek()
			if constant.BackSlash == peekByte {
				tokenBytes = append(tokenBytes, findEscape(currentByte))
				_ = bStack.Pop()
			} else if constant.Quote == peekByte || constant.Colon == peekByte {
				tokenBytes = append(tokenBytes, currentByte)
			} else if constant.LetterU == peekByte {
				unicodeBytes = append(unicodeBytes, currentByte)
				if 4 == len(unicodeBytes) {
					code, _ := strconv.ParseInt(string(unicodeBytes), 16, 32)
					ba := []byte(string([]rune{rune(code)}))
					tokenBytes = append(tokenBytes, ba...)

					_ = bStack.Pop()
					_ = bStack.Pop()
				}
			}
		}
	}
	return rootObject
}

func parseArray(jsonText string, startIndex int, endIndex map[string]int) *JsonArray {
	return NewArray()
}

func findEscape(b byte) byte {
	switch b {
	case '\\':
		return '\\'
	case '\'':
		return '\''
	case '"':
		return '"'
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'r':
		return '\r'
	default:
		panic("invalid Escape")
	}
}

func escapeString(value string) string {
	rs := []rune(value)
	length := len(rs)
	result := make([]rune, 0, length)
	for m := 0; m < length; m++ {
		switch rs[m] {
		case '\\':
			result = append(result, '\\', '\\')
		case '"':
			result = append(result, '\\', '"')
		case '\n':
			result = append(result, '\\', 'n')
		case '\t':
			result = append(result, '\\', 't')
		case '\b':
			result = append(result, '\\', 'b')
		case '\f':
			result = append(result, '\\', 'f')
		case '\r':
			result = append(result, '\\', 'r')
		default:
			if rs[m] < 128 {
				result = append(result, rs[m])
				break
			}
			result = append(result, '\\', 'u')
			result = append(result, []rune(strconv.FormatInt(int64(rs[m]), 16))...)
		}
	}
	return string(result)
}

func translateValue(value string) any {
	lowerValue := strings.ToLower(value)
	if regexp.MustCompile("^\"(.|\\s)*\"$").MatchString(value) {
		// Strings
		return value[1 : len(value)-1]
	} else if regexp.MustCompile("^\\d+$").MatchString(value) {
		result, _ := strconv.ParseInt(value, 10, 64)
		return result
	} else if regexp.MustCompile("^\\d+\\.\\d+$").MatchString(value) {
		result, _ := strconv.ParseFloat(value, 64)
		return result
	} else if "true" == lowerValue || "false" == lowerValue {
		result, _ := strconv.ParseBool(lowerValue)
		return result
	} else if "null" == value {
		return nil
	}
	return nil
}
