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

// generic
var (
	Mode          string = replaceIfEmpty(os.Getenv("APP_MODE"), "development")
	IsDevelopment bool   = (Mode == "development")
	IsProduction  bool   = (Mode == "production")
)

// api server
var (
	// Address : サーバのアドレス & ポート番号
	Address      string = replaceIfEmpty(os.Getenv("API_SERVER_ADDRESS"), "0.0.0.0:1323")
	Host         string = replaceIfEmpty(os.Getenv("API_SERVER_HOST"), "0.0.0.0")
	Port         string = replaceIfEmpty(os.Getenv("API_SERVER_PORT"), "1323")
	PublicPrefix string = "/file/"
)

var (
	PublicDirectory string = replaceIfEmpty(os.Getenv("WEB_PUBLIC_DIRECTORY"), "/home/app/web/public")
	LogFile         string = replaceIfEmpty(os.Getenv("API_SERVER_LOG_FILE"), "/home/app/logs/default.log")
	LogDirectory    string = replaceIfEmpty(os.Getenv("API_SERVER_LOG_DIR"), "/home/app/logs")
	SeedFile        string = replaceIfEmpty(os.Getenv("API_SERVER_SEED_FILE"), "/home/app/seeds/default.json")
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

func GetServerAddress() string {
	if Address != "" {
		return Address
	}
	return fmt.Sprintf("%s:%s", Host, Port)
}
