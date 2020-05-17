package apiv1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/log"

	"github.com/mndyu/localchat-server/database"
	utils "github.com/mndyu/localchat-server/test/utils"
)

var mockDB *gorm.DB

// TestMain 全テストの実行
func TestMain(m *testing.M) {
	defer utils.CloseDb(mockDB)
	db := getNewMockDB()
	fmt.Print(db != nil)
	var s = m.Run()
	os.Exit(s)
}

func getNewMockDB() *gorm.DB {
	mockDB = utils.ResetDB(mockDB)
	initDb(mockDB)
	loadDummyData(mockDB)
	return mockDB
}

// initDb テスト用 DB を初期化する関数
func initDb(db *gorm.DB) {
	// schemass
	for _, s := range database.AllSchemas {
		db.AutoMigrate(s)
		t := reflect.TypeOf(s).Name()
		log.Errorf("auto migrate %s failed: %s", t, err.Error())
	}
}

// loadDummyData :
// ダミーデータ読み込み
func loadDummyData(db *gorm.DB) {
	utils.ReadJSON("../test/testdata/dummy.json", database.AllSchemas, func(a interface{}) {
		//fmt.Printf("%v", a)
		if err := db.Create(a).Error; err != nil {
			t := reflect.TypeOf(a).Name()
			log.Errorf("loadDummyData %s failed: %s", t, err.Error())
		}
	})
}

type apiresult struct {
	status int
	body   interface{}
	record *httptest.ResponseRecorder
	err    error
}

func (r *apiresult) assertStatus(t *testing.T, expectedStatus int) bool {
	return assert.Equal(t, r.status, expectedStatus)
}
func (r *apiresult) assertBody(t *testing.T, expectedBody interface{}) bool {
	if str, ok := expectedBody.(string); ok {
		var err error
		expectedBody, err = utils.DecodeJSON(str)
		if err != nil {
			panic(fmt.Sprintf("expectedBody parse error %v", str))
		}
	}
	result, err := compareJsons(expectedBody, r.body)
	if err != nil {
		assert.Fail(t, err.Error())
	} else if result == false {
		assert.Fail(t, fmt.Sprintf("body not equal: %v != %v", expectedBody, r.body))
	}
	return result
}

//
type testapi struct {
	context *Context
}

func (a *testapi) Get(handlerFunc HandlerFunc, params interface{}, queryParams interface{}) (*apiresult, error) {
	// query string
	var err error
	var queryString string
	if queryParams != nil {
		var ok bool
		queryString, ok = queryParams.(string)
		if ok == false {
			uv, ok := queryParams.(url.Values)
			if !ok {
				return nil, fmt.Errorf("query string error: %v", queryParams)
			}
			queryString = uv.Encode()
		}
	}

	// url params
	var paramMap map[string]string
	if params != nil {
		switch v := params.(type) {
		case map[string]string:
			paramMap = v
		case string:
			paramValue, err := url.ParseQuery(v)
			if err != nil {
				return nil, fmt.Errorf("url params error: %v", params)
			}
			paramMap = map[string]string{}
			for key, value := range paramValue {
				paramMap[key] = value[0]
			}
		}
	}

	// ダミーのリクエスト作成
	e := echo.New()
	if paramMap != nil {
		utils.SetMaxParam(e, len(paramMap)) // echo v4 のみ必要
	}
	req := httptest.NewRequest(http.MethodGet, "/?"+queryString, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if paramMap != nil {
		for key, value := range paramMap {
			c.SetParamNames(key)
			c.SetParamValues(value)
		}
	}

	// handler 実行
	handlerErr := handlerFunc(a.context, c)
	if _, ok := handlerErr.(*echo.HTTPError); ok {
		e.DefaultHTTPErrorHandler(handlerErr, c)
	}

	// result
	body, err := utils.DecodeJSON(rec.Body.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to decode json string: %s, %s", rec.Body.String(), err.Error())
	}

	return &apiresult{rec.Code, body, rec, handlerErr}, nil
}

func (a *testapi) Post(handlerFunc HandlerFunc, params interface{}, reqBody interface{}) (*apiresult, error) {
	var err error

	// url params
	var paramMap map[string]string
	if params != nil {
		switch v := params.(type) {
		case map[string]string:
			paramMap = v
		case string:
			paramValue, err := url.ParseQuery(v)
			if err != nil {
				return nil, fmt.Errorf("url params error: %v", params)
			}
			paramMap = map[string]string{}
			for key, value := range paramValue {
				paramMap[key] = value[0]
			}
		}
	}

	// body
	reqString, ok := reqBody.(string)
	if ok == false {
		reqString, err = utils.ToJSON(reqBody)
		if err != nil {
			return nil, fmt.Errorf("reqbody parse error %v", reqBody)
		}
	}

	// ダミーのリクエスト作成
	e := echo.New()
	if paramMap != nil {
		utils.SetMaxParam(e, len(paramMap)) // echo v4 のみ必要
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if paramMap != nil {
		for key, value := range paramMap {
			c.SetParamNames(key)
			c.SetParamValues(value)
		}
	}

	// handler 実行
	handlerErr := handlerFunc(a.context, c)
	if _, ok := handlerErr.(*echo.HTTPError); ok {
		e.DefaultHTTPErrorHandler(handlerErr, c)
	}

	// result
	body, err := utils.DecodeJSON(rec.Body.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to decode json string: %s, %s", rec.Body.String(), err)
	}
	return &apiresult{rec.Code, body, rec, handlerErr}, nil
}
