package config

import (
	"fmt"
	"os"
)

// db
var (
	// SQLType : DBMS の種類
	SQLType string = os.Getenv("DB_TYPE")

	// DBURL : DB 接続用 URL
	dbURL string = os.Getenv("DB_URL")

	// DBAddress : DB サーバのアドレス
	dbAddress string = os.Getenv("DB_ADDRESS")

	// DBUser : DB のユーザー名
	dbUser string = os.Getenv("DB_USER")

	// DBPass : DB のパスワード
	dbPass string = os.Getenv("DB_PASSWORD")

	// DBName : DB 名
	dbName string = os.Getenv("DB_DATABASE")

	// DB
	// dbHost string = "localhost" // "mysql", "localhost"
	// dbPort    string = "3306"

	// DBProtocol : DB の接続プロトコル
	// dbProtocol string = "tcp(" + dbHost + ":" + dbPort + ")"
	dbProtocol string = "tcp(" + dbAddress + ")"

	// DBConfig : その他設定
	dbConfig string = "?sslmode=disable"

	// SQLfile : SQL ファイル？
	// SQLfile string = "creatorslab-ubuntu-googlegcp.sql"
)

// minio
var (
	accesssKey string
	secretKey  string
)

// api server
const (
	// Address : サーバのアドレス & ポート番号
	Address         string = ":1323"
	PublicDirectory string = "public"
	PublicPrefix    string = "/file/"
)

// GetConnectionURL :
// DB 接続用 URL を生成
func GetConnectionURL() string {
	if dbURL != "" {
		return dbURL
	}
	connectTemplate := "postgres://%s:%s@%s/%s%s"
	connect := fmt.Sprintf(connectTemplate, dbUser, dbPass, dbAddress, dbName, dbConfig)
	return connect
}
