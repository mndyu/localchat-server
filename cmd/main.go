package main

import (
	"net/http"
	"os"
	"time"

	"github.com/mndyu/localchat-server/apiv1"
	"github.com/mndyu/localchat-server/config"
	"github.com/mndyu/localchat-server/database"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	log "github.com/sirupsen/logrus"
)

type LogWriter struct {
	f *os.File
}

func NewLogWriter(filename string) *LogWriter {
	logFilePath := "./log/haha.txt"
	// logFilePath := "/Users/mon/dev/gw/localchat/localchat-server/log/haha.txt"
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return &LogWriter{f: f}
}
func (l LogWriter) Write(p []byte) (n int, err error) {
	n, err = os.Stdout.Write(p)
	n, err = l.f.Write(p)
	return
}
func (l LogWriter) Close() {
	l.f.Close()
	return
}

var defaultLogWriter *LogWriter

func init() {
	defaultLogWriter = NewLogWriter("")
	log.SetOutput(*defaultLogWriter)
}

func main() {
	runServer()
}

var a = middleware.BasicAuth

func runServer() {
	// logger
	defer defaultLogWriter.Close()

	// db: 接続
	var db *gorm.DB
	var err error
	retrySec := 5
	for {
		log.Infof("DB connection: %s %s", config.SQLType, config.GetConnectionURL())
		// db, err = gorm.Open("sqlite3", ":memory:") // テスト用インメモリDB
		db, err = gorm.Open(config.SQLType, config.GetConnectionURL()) // DBMS
		if err == nil {
			// panic("failed to connect database")
			break
		}
		log.Errorf("DB connection failed: %s", err.Error())
		log.Infof("DB connection: retrying in %d seconds ...", retrySec)
		time.Sleep(time.Duration(retrySec) * time.Second)
	}
	defer db.Close()

	// db: 初期化
	database.SetupDatabase(db)

	// echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Pre(middleware.HTTPSRedirect())
	// e.Use(middleware.CORS())
	// e.Use(middleware.CSRF())
	// e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	// static
	e.Static("/file", config.PublicDirectory)

	// echo: routes
	api := e.Group("/api")
	ver := api.Group("/v1")
	ver.GET("/poop", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"hid": "jo"})
	})
	apiv1.SetupRoutes(ver, db)

	// echo: start
	e.Logger.Fatal(e.Start(config.Address))

}
