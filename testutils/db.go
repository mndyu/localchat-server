package testutils

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

// CreateDb :
// テスト用インメモリDBの作成
func CreateDb() *gorm.DB {
	// db 作成 (インメモリ)
	var err error
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err.Error()))
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

// RecreateDb :
// テスト用DBの再作成
func RecreateDb(db *gorm.DB) *gorm.DB {
	CloseDb(db)
	return CreateDb()
}

// MigrateAll :
// テスト用DBでマイグレーション
func MigrateAll(db *gorm.DB) {

}
