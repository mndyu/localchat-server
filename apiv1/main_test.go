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

	"github.com/mndyu/localchat-server/database"
	utils "github.com/mndyu/localchat-server/testutils"
)

var mockDB *gorm.DB

// TestMain 全テストの実行
func TestMain(m *testing.M) {
	defer utils.CloseDb(mockDB)
	var s = m.Run()
	os.Exit(s)
}

func getNewMockDB() *gorm.DB {
	mockDB := utils.RecreateDb(mockDB)
	initDb(mockDB)
	loadDummyData(mockDB)
	return mockDB
}

// initDb テスト用 DB を初期化する関数
func initDb(db *gorm.DB) {
	// schemass
	for _, s := range database.All {
		db.AutoMigrate(s)
	}
}

// loadDummyData :
// ダミーデータ読み込み
func loadDummyData(db *gorm.DB) {
	all := []interface{}{
		database.User{},
		database.Message{},
		database.Group{},
		database.Channel{},
	}
	utils.ReadJSON("../testdata/dummy.json", all, func(a interface{}) {
		fmt.Println(reflect.TypeOf(a))
		fmt.Printf("%d", reflect.TypeOf(a).Kind())
		fmt.Printf("%d", reflect.ValueOf(a).Type().Kind())
		db.Create(a)
	})
}

//
func indirect(a interface{}) interface{} {
	if reflect.TypeOf(a).Kind() == reflect.Ptr {
		return reflect.Indirect(reflect.ValueOf(a)).Interface()
	}
	return a
}
func convert(a interface{}, t reflect.Type) interface{} {
	v := reflect.ValueOf(a)
	if v.Type().ConvertibleTo(t) == false {
		return nil
	}
	return v.Convert(t).Interface()
}
func compareJsonMaps(expected jsonmap, target jsonmap) bool {
	for k, v := range expected {
		tv := target[k]
		if reflect.TypeOf(v) != reflect.TypeOf(tv) {
			tv = convert(tv, reflect.TypeOf(v))
			if tv == nil {
				return false
			}
		}
		switch v.(type) {
		case jsonmap:
			if compareJsonMaps(v.(jsonmap), tv.(jsonmap)) == false {
				return false
			}
		default:
			if v != tv {
				return false
			}
		}
	}
	return true
}
func compareJsons(expected interface{}, target interface{}) (bool, error) {
	if reflect.TypeOf(expected) != reflect.TypeOf(target) {
		return false, fmt.Errorf("types dont match: %s, %s", reflect.TypeOf(expected), reflect.TypeOf(target))
	}

	switch v := expected.(type) {
	case jsonmap:
		jov := target.(jsonmap)
		if compareJsonMaps(v, jov) == false {
			return false, nil
		}
	case []interface{}:
		jov := target.([]interface{})
		for i, item := range v {
			r, err := compareJsons(item, jov[i])
			if !r || err != nil {
				return false, err
			}
		}
	default:
		if expected != target {
			return false, nil
		}
	}
	return true, nil
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
