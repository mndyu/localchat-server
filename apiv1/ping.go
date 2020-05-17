package apiv1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type pingResultJson struct {
	Time time.Time `json:"time"`
}

// GetPing GET /ping
func GetPing(context *Context, c echo.Context) error {
	jsonData := pingResultJson{
		Time: time.Now(),
	}
	return c.JSON(http.StatusOK, jsonData)
}
