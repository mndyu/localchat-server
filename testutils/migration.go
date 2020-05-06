package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// case 変換のパターン
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchUnderscore = regexp.MustCompile("_([A-Za-z])")

// toSnakeCase :
// camel case -> snake case に変換
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// toCamelCase :
// snake case -> lower camel case に変換
func toCamelCase(str string) string {
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
func toUpperCamelCase(str string) string {
	str = strings.ToUpper(string(str[0])) + str[1:]
	camel := matchUnderscore.ReplaceAllStringFunc(str, func(m string) string {
		return strings.ToUpper(m)
	})
	camel = matchUnderscore.ReplaceAllString(camel, "${1}")
	return camel
}

// searchFieldCaseInsensitive :
// struct から field を case insensitive で検索
func searchFieldCaseInsensitive(target interface{}, name string) (string, error) {
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

// assignValue :
// reflect.Value 同士で値代入
func assignValue(a, b reflect.Value) bool {
	if a.CanSet() == false {
		return false
	}
	// if b.Type().ConvertibleTo(a.Type()) == false {
	// 	return false
	// }
	bv := b.Convert(a.Type())
	a.Set(bv)

	return true
}

// ReadJSON :
// JSON から DB の要素を読み取る
// JSON の形式:
//   {
//     "スネークケースのテーブル名": [
//       {
//         "スネークケースのカラム名": <値>
// 				 ...
//       }
//       ...
//     ],
//     ...
//   }
// 行ごとに callback 呼び出し
func ReadJSON(path string, schemas []interface{}, callback func(a interface{})) error {
	// ファイル読み込み
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("error reading %s", path)
	}
	// JSON 読み込み
	var jsonData map[string][]map[string]interface{}
	json.Unmarshal(raw, &jsonData)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("error parsing %s as json", path)
	}
	// テーブル
	for _, s := range schemas {
		schemaType := reflect.TypeOf(s)
		fieldName := toSnakeCase(schemaType.Name())
		items := jsonData[fieldName]
		// 行
		for _, i := range items {
			new := reflect.New(schemaType)
			// 列
			for k, v := range i {
				// fieldName := toUpperCamelCase(k)
				fieldName, err := searchFieldCaseInsensitive(s, toUpperCamelCase(k))
				if err != nil {
					continue
				}
				f := new.Elem().FieldByName(fieldName)
				if f.IsValid() {
					log.Info(k, ":", v)
					assignValue(f, reflect.ValueOf(v))
				} else {
					log.Fatal("why")
				}
			}
			callback(new.Interface())
		}
	}
	return nil
}

// exampleReadJSON :
// ReadJSON の使用例
func exampleReadJSON() {
	type vendorCards struct {
		CardID   int
		CardName string
	}
	type masterLangs struct {
		LangID   int
		LangName string
		LangCode string
	}
	schemas := []interface{}{
		vendorCards{},
		masterLangs{},
	}
	ReadJSON("a.json", schemas, func(a interface{}) {
		fmt.Println(reflect.TypeOf(a).Name(), ":", a)
	})
}
