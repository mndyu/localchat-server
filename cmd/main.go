package main

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/mndyu/localchat-server/apiv1"
	"github.com/mndyu/localchat-server/config"
	"github.com/mndyu/localchat-server/database"
	"github.com/mndyu/localchat-server/database/schema"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	log "github.com/sirupsen/logrus"
)

// 例: 202005302307
func getLogFileName() string {
	now := time.Now()
	return now.Format("200612150405.log")
}

type LogWriter struct {
	f    *os.File
	date time.Time
}

func NewLogWriter(filepath string) (*LogWriter, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &LogWriter{f: f, date: time.Now()}, nil
}
func (l LogWriter) isLogfileOld() bool {
	ld := l.date
	y, m, d := time.Now().Date()
	return y != ld.Year() && m != ld.Month() && d != ld.Day()
}
func (l LogWriter) Write(p []byte) (n int, err error) {
	if l.isLogfileOld() {
		newFilePath := path.Join(config.LogDirectory, getLogFileName())
		nlw, err := NewLogWriter(newFilePath)
		if err != nil {
			log.Errorf("failed to open log file %s: %s", newFilePath, err.Error())
			return 0, err
		}
		l = *nlw
	}
	n, err = os.Stdout.Write(p)
	n, err = l.f.Write(p)
	return
}
func (l LogWriter) Close() {
	if l.f != nil {
		l.f.Close()
	}
	return
}

var defaultLogWriter *LogWriter = &LogWriter{}

func init() {
	var err error
	defaultLogWriter, err = NewLogWriter(config.LogDirectory)
	if err != nil {
		log.Errorf("failed to open log file %s: %s", config.LogDirectory, err.Error())
		return
	}
	log.SetOutput(*defaultLogWriter)
}

func main() {
	runServer()
}

func runServer() {
	// logger
	defer func() {
		defaultLogWriter.Close()
	}()

	// db: 接続
	db, err := database.Connect(config.SQLType, config.GetConnectionURL(), 5)
	if err != nil {
		log.Fatalf("failed to connect to db %s", err.Error())
	}
	defer db.Close()

	// db: 初期化
	database.SetupDatabase(db)

	// TODO: 後で消す
	group := schema.Group{}
	group.ID = 1
	group.Name = "default"
	err = db.Create(&group).Error
	if err != nil {
		panic(err)
	}

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
	e.Logger.Fatal(e.Start(config.GetServerAddress()))

}
