package config

import (
	"fmt"
	"os"
)

func replaceIfEmpty(s string, replace string) string {
	if s == "" {
		return replace
	}
	return s
}

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
	Address      string = ":1323"
	PublicPrefix string = "/file/"
)

var (
	PublicDirectory string = replaceIfEmpty(os.Getenv("WEB_PUBLIC_DIRECTORY"), "/home/app/web/public")
	LogFile         string = replaceIfEmpty(os.Getenv("SERVER_LOG_FILE"), "/home/app/logs/default.log")
	SeedFile        string = replaceIfEmpty(os.Getenv("API_SERVER_SEED_FILE"), "/home/app/seeds/default.json")
)

// GetConnectionURL :
// DB 接続用 URL を生成
func GetConnectionURL() string {
	if dbURL != "" {
		return dbURL
	}
	if dbAddress == "" {
		return ""
	}
	connectTemplate := "postgres://%s:%s@%s/%s%s"
	connect := fmt.Sprintf(connectTemplate, dbUser, dbPass, dbAddress, dbName, dbConfig)
	return connect
}
