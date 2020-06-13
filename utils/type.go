package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// SearchFieldCaseInsensitive :
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

func IndirectValue(a reflect.Value) reflect.Value {
	if a.Type().Kind() == reflect.Ptr {
		return a.Elem()
	}
	return a
}

func IndirectType(a reflect.Type) reflect.Type {
	if a.Kind() == reflect.Ptr {
		return a.Elem()
	}
	return a
}

func Indirect(a interface{}) interface{} {
	if reflect.TypeOf(a).Kind() == reflect.Ptr {
		return reflect.Indirect(reflect.ValueOf(a)).Interface()
	}
	return a
}
func Convert(a interface{}, t reflect.Type) interface{} {
	v := reflect.ValueOf(a)
	if v.Type().ConvertibleTo(t) == false {
		return nil
	}
	return v.Convert(t).Interface()
}

// AssignJSONValue :
// reflect.Value 同士で値代入
func AssignJSONValue(a, b reflect.Value) error {
	if a.CanSet() == false {
		return fmt.Errorf("cant set %v", a)
	}

	ae := IndirectValue(a)
	be := IndirectValue(b)
	if !be.IsValid() && a.Type().Kind() == reflect.Ptr {
		// *ptr = nil
		a.Set(reflect.Zero(a.Type()))
		return nil
	}
	if !ae.IsValid() {
		ae = reflect.New(a.Type().Elem()).Elem()
	}
	if !be.IsValid() {
		return fmt.Errorf("invalid %v(%v), %v(%v)", a.Type(), a, b.Type(), b)
	}

	switch ae.Interface().(type) {
	case time.Time:
		if beString, ok := be.Interface().(string); ok {
			t, err := time.Parse(time.RFC3339, beString)
			if err != nil {
				return err
			}
			be = InterfaceToValue(t)
		}
	case decimal.Decimal:
		d, err := ToDecimal(be.Interface())
		if err != nil {
			return err
		}
		be = InterfaceToValue(d)
	}

	switch n := be.Interface().(type) {
	case float64:
		if float64(int(n)) == n {
			be = InterfaceToValue(int64(n))
		}
	}

	if be.Type().ConvertibleTo(ae.Type()) == false {
		return fmt.Errorf("cant convert %s to %s", be.Type().String(), ae.Type().String())
	}
	bv := be.Convert(ae.Type())
	ae.Set(bv)

	// assign ptr
	if a.Type().Kind() == reflect.Ptr && ae.CanAddr() {
		a.Set(ae.Addr())
	}

	return nil
}

func CompareJsons(expected interface{}, target interface{}) (bool, error) {
	if expected == nil && target == nil {
		return true, nil
	}
	expectedType := reflect.TypeOf(expected)
	targetType := reflect.TypeOf(target)
	if expectedType == nil && targetType == nil {
		return true, nil
	}
	if expectedType == nil || targetType == nil {
		return false, fmt.Errorf("nil type %v(%v), %v(%v)", expectedType, expected, targetType, target)
	}
	if expectedType != targetType && !expectedType.ConvertibleTo(targetType) {
		return false, fmt.Errorf("types dont match: %s, %s", reflect.TypeOf(expected), reflect.TypeOf(target))
	}

	switch expected.(type) {
	case Jsonmap, map[string]interface{}:
		var ejm Jsonmap = InterfaceToJsonmap(expected)
		var tjm Jsonmap = InterfaceToJsonmap(target)
		for k, ev := range ejm {
			tv := tjm[k]
			result, err := CompareJsons(ev, tv)
			if !result || err != nil {
				return false, err
			}
		}
	case []interface{}:
		ear := expected.([]interface{})
		tar := target.([]interface{})
		if len(ear) != len(tar) {
			return false, nil
		}
		for i, item := range ear {
			r, err := CompareJsons(item, tar[i])
			if !r || err != nil {
				return false, err
			}
		}
	default:
		if expected != target {
			return false, nil
		}
	}
	return true, nil
}
