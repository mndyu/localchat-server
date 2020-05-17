package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

// case 変換のパターン
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchUnderscore = regexp.MustCompile("_([A-Za-z])")

// ToSnakeCase :
// camel case -> snake case に変換
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// toCamelCase :
// snake case -> lower camel case に変換
func ToCamelCase(str string) string {
	if str[0] == '_' {
		str = str[1:]
	}
	camel := matchUnderscore.ReplaceAllStringFunc(str, func(m string) string {
		return strings.ToUpper(m)
	})
	camel = matchUnderscore.ReplaceAllString(camel, "${1}")
	return camel
}

// toUpperCamelCase :
// snake case -> upper camel case に変換
func ToUpperCamelCase(str string) string {
	str = strings.ToUpper(string(str[0])) + str[1:]
	camel := matchUnderscore.ReplaceAllStringFunc(str, func(m string) string {
		return strings.ToUpper(m)
	})
	camel = matchUnderscore.ReplaceAllString(camel, "${1}")
	return camel
}

// searchFieldCaseInsensitive :
// struct から field を case insensitive で検索
func SearchFieldCaseInsensitive(target interface{}, name string) (string, error) {
	value := reflect.TypeOf(target)
	ln := strings.ToLower(name)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := field.Name
		if strings.ToLower(fieldName) == ln {
			return fieldName, nil
		}
	}
	return "", fmt.Errorf("field not found %s in %v", name, target)
}

func InterfaceToValue(a interface{}) reflect.Value {
	return reflect.ValueOf(interface{}(a))
}

func ToDecimal(a interface{}) (d decimal.Decimal, err error) {
	switch v := a.(type) {
	case decimal.Decimal:
		d = v
	case int32:
		d = decimal.NewFromInt32(v)
	case int64:
		d = decimal.NewFromInt(v)
	case float32:
		d = decimal.NewFromFloat32(v)
	case float64:
		d = decimal.NewFromFloat(v)
	case string:
		d, err = decimal.NewFromString(v)
	default:
		err = fmt.Errorf("unsupported type: %s", reflect.TypeOf(a).Name())
	}
	return
}

// AssignJSONValue :
// reflect.Value 同士で値代入
func AssignJSONValue(a, b reflect.Value) error {
	if a.CanSet() == false {
		return fmt.Errorf("cant set %v", a)
	}

	switch a.Interface().(type) {
	case decimal.Decimal:
		d, err := ToDecimal(b.Interface())
		if err != nil {
			return err
		}
		b = InterfaceToValue(d)
	}

	switch n := b.Interface().(type) {
	case float64:
		if float64(int(n)) == n {
			b = InterfaceToValue(int64(n))
		}
	}

	if b.Type().ConvertibleTo(a.Type()) == false {
		return fmt.Errorf("cant convert %s to %s", b.Type().String(), a.Type().String())
	}
	bv := b.Convert(a.Type())
	a.Set(bv)

	return nil
}
