package apiv1

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database/schema"
	"github.com/mndyu/localchat-server/utils"
)

type (
	Context struct {
		DB *gorm.DB
	}
	HandlerFunc func(context *Context, c echo.Context) error
	jsonmap     map[string]interface{}
)

func (ctx *Context) bind(f HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return f(ctx, c)
	}
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

func compareJsons(expected interface{}, target interface{}) (bool, error) {
	if expected == nil && target == nil {
		return true, nil
	}
	expectedType := reflect.TypeOf(expected)
	targetType := reflect.TypeOf(target)
	if targetType == nil || expectedType != targetType && !expectedType.ConvertibleTo(targetType) {
		return false, fmt.Errorf("types dont match: %s, %s", reflect.TypeOf(expected), reflect.TypeOf(target))
	}

	switch expected.(type) {
	case jsonmap, map[string]interface{}:
		var ejm jsonmap = interfaceToJsonmap(expected)
		var tjm jsonmap = interfaceToJsonmap(target)
		for k, ev := range ejm {
			tv := tjm[k]
			result, err := compareJsons(ev, tv)
			if !result || err != nil {
				return false, err
			}
		}
	case []interface{}:
		ear := expected.([]interface{})
		tar := target.([]interface{})
		if len(ear) != len(tar) {
			return false, nil
		}
		for i, item := range ear {
			r, err := compareJsons(item, tar[i])
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

func interfaceToJsonmap(a interface{}) jsonmap {
	switch v := a.(type) {
	case jsonmap:
		return v
	case map[string]interface{}:
		return jsonmap(v)
	}
	return nil
}

func eachJSONStructField(a interface{}, f func(fieldName string, val reflect.Value)) {
	var uv reflect.Value
	if v, ok := a.(reflect.Value); ok {
		uv = v
	} else {
		uv = reflect.ValueOf(a)
	}
	if uv.Type().Kind() == reflect.Ptr {
		uv = uv.Elem()
	}
	ut := uv.Type()

	for i := 0; i < ut.NumField(); i++ {
		uf := ut.Field(i)
		ufv := uv.Field(i)
		if uf.Anonymous {
			eachJSONStructField(ufv, f)
		}
		if strings.ToLower(uf.Name) == uf.Name {
			continue
		}
		fieldName := uf.Tag.Get("json")
		if fieldName == "" {
			fieldName = utils.ToSnakeCase(uf.Name)
			// continue
		}
		f(fieldName, ufv)
	}
}

func structToJsonmap(out *jsonmap, a interface{}) error {
	eachJSONStructField(a, func(fieldName string, val reflect.Value) {
		(*out)[fieldName] = val.Interface()
	})
	return nil
}
func assignMapToStruct(out interface{}, a jsonmap) error {
	eachJSONStructField(out, func(fieldName string, val reflect.Value) {
		v := val //reflect.ValueOf(i)
		if !v.IsValid() || !v.CanSet() {
			// fmt.Println("cant set:", uf.Name)
			return
		}
		fieldValue := a[fieldName]
		if fieldValue == nil {
			return
		}
		fv := reflect.ValueOf(fieldValue)
		err := utils.AssignJSONValue(v, fv)
		if err != nil {
			t1 := fv.Type().String()
			t2 := v.Type().String()
			fmt.Println(err, t1, t2)
			return
			// return fmt.Errorf("field %s of type %s is not assignable to type %s", uf.Name, ufv.Type().Name(), fv.Type().Name())
		}
	})
	return nil
}

func filterJsonmapWithStruct(jm jsonmap, i interface{}) jsonmap {
	newMap := jsonmap{}
	eachJSONStructField(i, func(fieldName string, val reflect.Value) {
		if jm[fieldName] != nil {
			newMap[fieldName] = jm[fieldName]
		}
	})
	return newMap
}

// assignJSONFields
// *out = userStruct
func assignJSONFields(out interface{}, userStruct interface{}) error {
	// fmt.Println("assignJSONFields:", reflect.TypeOf(out).Name(), reflect.TypeOf(userStruct).Name())

	var inMap jsonmap = interfaceToJsonmap(userStruct)
	if inMap == nil {
		inMap = jsonmap{}
		structToJsonmap(&inMap, userStruct)
	}
	var outMap jsonmap = interfaceToJsonmap(out)
	if outMap == nil {
		assignMapToStruct(out, inMap)
	} else {
		for k, v := range inMap {
			outMap[k] = v
		}
	}

	return nil
}
func assignJSONArrayFields(out interface{}, userArray interface{}) error {
	// fmt.Println("assignJSONArrayFields:", reflect.TypeOf(out).Name(), reflect.TypeOf(userArray).Name())
	uv := reflect.ValueOf(userArray)
	dv := reflect.Indirect(reflect.ValueOf(out))
	dElemType := dv.Type().Elem()

	for i := 0; i < uv.Len(); i++ {
		item := uv.Index(i)
		newItem := reflect.New(dElemType)
		err := assignJSONFields(newItem.Interface(), item.Interface())
		if err != nil {
			return err
		}
		dv.Set(reflect.Append(dv, reflect.Indirect(newItem)))
	}
	return nil
}

func scanStructs(db *gorm.DB, structs ...interface{}) error {
	rows, err := db.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		for _, s := range structs {
			err := db.ScanRows(rows, s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getRowByID(db *gorm.DB, c echo.Context, s interface{}, id string) (interface{}, error) {
	if id == "" {
		return s, fmt.Errorf("param not found")
	}
	if err := db.First(s, id).Error; err != nil {
		return s, fmt.Errorf("not found: %s", err.Error())
	}
	return s, nil
}

func getLimit(c echo.Context) interface{} {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return 10 // default
	}
	return limit
}
func getOffset(c echo.Context) interface{} {
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return nil // default
	}
	return offset
}
func parseTimeParam(c echo.Context, timeString string, defaultTime time.Time) time.Time {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return defaultTime
	}
	return t
}
func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}
func tomorrow() time.Time {
	now := time.Now().AddDate(0, 0, 1)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}
func getFrom(c echo.Context) time.Time {
	return parseTimeParam(c, c.QueryParam("from"), today())
}
func getTo(c echo.Context) time.Time {
	return parseTimeParam(c, c.QueryParam("to"), tomorrow())
}

func filterUserData(orig interface{}, ud interface{}, fields []string) interface{} {
	var origvalue = reflect.ValueOf(orig)
	var uvalue = reflect.ValueOf(ud)
	var newData = reflect.Indirect(reflect.New(reflect.TypeOf(ud)))

	for i := 0; i < reflect.TypeOf(ud).NumField(); i++ {
		from := origvalue.Field(i)
		to := newData.Field(i)
		to.Set(from)
	}

	for _, fn := range fields {
		from := uvalue.FieldByName(fn)
		to := newData.FieldByName(fn)
		to.Set(from)
	}
	return newData.Interface()
}

func filterForJson(data interface{}, fields []string) jsonmap {
	ret := jsonmap{}

	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	for _, fieldName := range fields {
		if strings.ToLower(fieldName) == fieldName {
			continue
		}
		vf := v.FieldByName(fieldName)
		tf, found := t.FieldByName(fieldName)
		if found {
			ret[tf.Name] = vf.Interface()
		}
	}
	return ret
}

func getRowById(db *gorm.DB, c echo.Context, s interface{}, id string) (interface{}, error) {
	if id == "" {
		return s, fmt.Errorf("param not found")
	}
	if err := db.First(s, id).Error; err != nil {
		return s, fmt.Errorf("not found: %s", err.Error())
	}
	return s, nil
}
func getClientUser(context *Context, c echo.Context) (schema.User, error) {
	db := context.DB
	ip := c.RealIP()

	var user schema.User
	if err := db.First(&user, "ip_address = ?", ip).Error; err != nil {
		return user, fmt.Errorf("user not found: %s", err.Error())
	}
	return user, nil
}

// ----------------------------

/*
search: g\.(.*)\("(.*)", bindDB\((.*), db\)\)
replcae: // $3 $1 $2\nfunc $3(context *Context, c echo.Context) error {
  db := context.DB\n}\n
*/

// SetupRoutes setups routes
func SetupRoutes(g *echo.Group, db *gorm.DB) {
	c := Context{db}

	g.GET("/ping", c.bind(GetPing))
	g.GET("/ws", c.bind(GetWs))
	g.POST("/ws", c.bind(PostWs))
	g.POST("/upload", c.bind(PostUpload))

	//
	g.GET("/profile", c.bind(GetProfile))
	g.POST("/users", c.bind(PostUsers))
	g.GET("/users", c.bind(GetUsers))
	g.PUT("/users/:id", c.bind(PutUserByID))
	g.GET("/users/:id", c.bind(GetUserByID))
	g.DELETE("/users/:id", c.bind(DeleteUsersByID))
	// g.POST("/users/:id/messages", c.bind(PostUserMessages))
	g.GET("/users/:id/messages", c.bind(GetUserMessages))
	g.GET("/users/:id/groups", c.bind(GetUserGroups))
	g.GET("/users/:id/channels", c.bind(GetUserChannels))

	g.POST("/messages", c.bind(PostMessages))
	g.GET("/messages", c.bind(GetMessages))
	g.GET("/messages/:id", c.bind(GetMessageByID))
	g.PUT("/messages/:id", c.bind(PutMessageByID))
	g.DELETE("/messages/:id", c.bind(DeleteMessageByID))

	g.POST("/groups", c.bind(PostGroups))
	g.GET("/groups", c.bind(GetGroups))
	g.GET("/groups/:id", c.bind(GetGroupByID))
	g.PUT("/groups/:id", c.bind(PutGroupByID))
	g.DELETE("/groups/:id", c.bind(DeleteGroupByID))
	g.POST("/groups/:id/members", c.bind(PostGroupMembers))
	g.GET("/groups/:id/members", c.bind(GetGroupMembers))
	g.DELETE("/groups/:id/members/:user_id", c.bind(DeleteGroupMemberByID))
	g.GET("/groups/:id/channels", c.bind(GetGroupChannels))
	g.DELETE("/groups/:id/channels/:channel_id", c.bind(DeleteGroupChannelByID))
	g.POST("/groups/:id/messages", c.bind(PostGroupMessages))

	g.POST("/channels", c.bind(PostChannels))
	g.GET("/channels", c.bind(GetChannels))
	g.PUT("/channels/:id", c.bind(PutChannelByID))
	g.GET("/channels/:id", c.bind(GetChannelByID))
	g.DELETE("/channels/:id", c.bind(DeleteChannelByID))
	g.GET("/channels/:id/members", c.bind(GetChannelMembers))
	g.DELETE("/channels/:id/members/:user_id", c.bind(DeleteChannelMemberByID))
}
