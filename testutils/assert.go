package testutils

import (
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// AssertObj :
// obj 同士を比較, assert
func AssertObj(t *testing.T, expected, actual interface{}) bool {
	b, err := CompareObj(expected, actual)
	if err != nil {
		return assert.Fail(t, err.Error())
	}
	return b
}

// AssertObjFields :
// obj 同士を比較, assert (フィールド指定)
func AssertObjFields(t *testing.T, expected, actual interface{}, fields []string) bool {
	b, err := CompareObjFields(expected, actual, fields)
	if err != nil {
		return assert.Fail(t, err.Error())
	}
	return b
}

// setMaxParam :
// echo サーバの param 数制限設定
// echo v4 では maxParam の制限があるため無理やり拡張
func SetMaxParam(e *echo.Echo, maxParam int) {
	p := echo.NewRouter(e)
	p.Add("GET", strings.Repeat("/:d", maxParam), func(c echo.Context) error { return nil })
}
