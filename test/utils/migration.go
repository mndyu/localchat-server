package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/mndyu/localchat-server/utils"
)

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
		fieldName := utils.ToSnakeCase(schemaType.Name())
		items := jsonData[fieldName]
		// 行
		for _, i := range items {
			new := reflect.New(schemaType)
			empty := true
			// 列
			for k, v := range i {
				// fieldName := toUpperCamelCase(k)
				fieldName, err := searchFieldCaseInsensitive(s, utils.ToUpperCamelCase(k))
				if err != nil {
					continue
				}
				f := new.Elem().FieldByName(fieldName)
				if f.IsValid() {
					// log.Info(k, ":", v, " : ", reflect.TypeOf(v).Kind())
					err := AssignJSONValue(f, reflect.ValueOf(v))
					if err != nil {
						return err
					}
					empty = false
				} else {
					log.Fatal("why")
				}
			}
			if !empty {
				callback(new.Interface())
			}
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
