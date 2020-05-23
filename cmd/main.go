package main

import (
	"net/http"
	"os"

	"github.com/mndyu/localchat-server/apiv1"
	"github.com/mndyu/localchat-server/config"
	"github.com/mndyu/localchat-server/database"

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

func NewLogWriter(filepath string) (*LogWriter, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &LogWriter{f: f}, nil
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
	var err error
	defaultLogWriter, err = NewLogWriter(config.LogFile)
	if err != nil {
		log.Errorf("failed to open log file %s", config.LogFile)
		return
	}
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
	db, err := database.Connect(config.SQLType, config.GetConnectionURL(), 5)
	if err != nil {
		log.Fatalf("failed to connect to db %s", err.Error())
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
