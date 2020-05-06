package testutils

import (
	"fmt"
	"reflect"
)

// StatusComparisonFunc :
// status 比較関数
type StatusComparisonFunc func(expected int, actual int) (bool, error)

// BodyComparisonFunc :
// body 比較関数
type BodyComparisonFunc func(expected interface{}, actual interface{}) (bool, error)

// CompareStatus :
// status 比較
func CompareStatus(expected int, actual int) (bool, error) {
	return expected == actual, nil
}

// CompareStatusClass :
// status の class 比較
func CompareStatusClass(expected int, actual int) (bool, error) {
	return expected/100 == actual/100, nil
}

// CompareObj :
// obj 同士を比較
func CompareObj(a, b interface{}) (bool, error) {
	if reflect.DeepEqual(a, b) {
		return true, nil
	}
	return false, fmt.Errorf("Not equal: expected: %s, actual  : %s", a, b)
}

// CompareObjFields :
// obj 同士を比較 (キー指定 & struct 以外)
func CompareObjFields(a, b interface{}, keys []string) (bool, error) {
	// 型取得
	var ta = reflect.TypeOf(a).Kind()
	var tb = reflect.TypeOf(b).Kind()

	if len(keys) == 0 || ta == reflect.Array {
		// 全キー比較, または配列の場合は単純に比較
		return CompareObj(a, b)
	} else if ta != tb || ta != reflect.Map {
		// a, b の型が異なる, または型が map 以外の時
		return false, fmt.Errorf("Invalid operation: %#v != %#v, %#v, %#v", ta, tb, a, b)
	}

	var va = reflect.ValueOf(a)
	var vb = reflect.ValueOf(b)

	for _, keyStr := range keys {
		// キーから値取得
		key := reflect.ValueOf(keyStr)
		var ma = va.MapIndex(key)
		var mb = vb.MapIndex(key)

		if !ma.IsValid() || !mb.IsValid() {
			// どちらかが nil の場合
			if ma.IsValid() != mb.IsValid() {
				// どちらか片方だけが nil の場合
				return false, fmt.Errorf("Not equal (key: %s): \n"+
					"expected: %s\n"+
					"actual  : %s", keyStr, a, b)
			}
			continue
		}

		// キーの値同士を比較
		var ia = ma.Interface()
		var ib = mb.Interface()

		if !reflect.DeepEqual(ia, ib) {
			return false, fmt.Errorf("Not equal (key: %s): \n"+
				"expected: %s\n"+
				"actual  : %s", keyStr, a, b)
		}
	}
	return true, nil
}

// CompareObjAsJSONString :
// obj 同士を JSON に変換して比較
func CompareObjAsJSONString(a, b interface{}) (bool, error) {
	as, err := ToJSON(a)
	if err != nil {
		return false, err
	}
	bs, err := ToJSON(b)
	if err != nil {
		return false, err
	}
	return as == bs, nil
}
