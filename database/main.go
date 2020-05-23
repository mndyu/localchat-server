package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mndyu/localchat-server/config"
	"github.com/mndyu/localchat-server/database/schema"
	"github.com/mndyu/localchat-server/utils"
	log "github.com/sirupsen/logrus"
)

func SetupDatabase(db *gorm.DB) {
	migrateSchemas(db)
}

func replaceIfEmpty(s string, replace string) string {
	if s == "" {
		return replace
	}
	return s
}

// migrateSchemas :
// スキーマ作成
func migrateSchemas(db *gorm.DB) {
	// Migrate the schemas
	for _, schema := range schema.All {
		// schema := replaceSchema(s)
		if err := db.AutoMigrate(schema).Error; err != nil {
			t := reflect.TypeOf(schema).Name()
			log.Fatalf("auto migrate %s failed: %s", t, err.Error())
		}
	}
}

func setupGormDB(db *gorm.DB) *gorm.DB {
	db.SingularTable(true)
	return db
}

func Connect(sqlType string, connectionURL string, retrySec int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	sqlType = replaceIfEmpty(sqlType, "sqlite3")
	connectionURL = replaceIfEmpty(config.GetConnectionURL(), ":memory:")

	for {
		log.Infof("DB connection: %s %s", sqlType, connectionURL)
		db, err = gorm.Open(sqlType, connectionURL) // DBMS
		if err == nil {
			// success
			break
		}
		log.Errorf("DB connection failed: %s", err.Error())
		log.Infof("DB connection: retrying in %d seconds ...", retrySec)
		time.Sleep(time.Duration(retrySec) * time.Second)
	}

	db = setupGormDB(db)
	return db, err
}

// ReadSeedFile :
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
func ReadSeedFile(path string, schemas []interface{}, callback func(a interface{})) error {
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
				fieldName, err := utils.SearchFieldCaseInsensitive(s, utils.ToUpperCamelCase(k))
				if err != nil {
					continue
				}
				f := new.Elem().FieldByName(fieldName)
				if f.IsValid() {
					// log.Info(k, ":", v, " : ", reflect.TypeOf(v).Kind())
					err := utils.AssignJSONValue(f, reflect.ValueOf(v))
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
