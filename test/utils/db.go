package utils

import (
	"fmt"
	"time"

	"github.com/mndyu/localchat-server/config"
	"github.com/mndyu/localchat-server/database/schema"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

var (
	sqlType       string
	connectionURL string
	retrySec      int = 5
)

func replaceIfEmpty(s string, replace string) string {
	if s == "" {
		return replace
	}
	return s
}

// CreateDb :
// テスト用インメモリDBの作成
func CreateDb() *gorm.DB {
	var db *gorm.DB
	var err error

	sqlType = replaceIfEmpty(config.SQLType, "sqlite3")
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
	return db
}

// CloseDb :
// テスト用DBの終了
func CloseDb(db *gorm.DB) {
	if db != nil {
		db.Close()
	}
}

// ResetDB :
// テスト用DBの再作成
func ResetDB(db *gorm.DB) *gorm.DB {
	if db == nil {
		return CreateDb()
	}
	if connectionURL == ":memory:" {
		CloseDb(db)
		return CreateDb()
	}

	// truncate all tables
	for _, s := range schema.All {
		tableName := db.Model("").NewScope(s).TableName()
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
		if err != nil {
			log.Errorf("truncate error: %s", err.Error())
		}
	}
	return db
}

// MigrateAll :
// テスト用DBでマイグレーション
func MigrateAll(db *gorm.DB) {

}
