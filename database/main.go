package database

import "github.com/jinzhu/gorm"

func SetupDatabase(db *gorm.DB) {
	db.AutoMigrate(
		&User{},
		&Message{},
		&Group{},
		&Channel{},
	)
}

// func CreateMockDB() *gorm.DB {
// 	db, err = gorm.Open("sqlite3", ":memory:") // テスト用インメモリDB

// }
