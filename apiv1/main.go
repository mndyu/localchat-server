package apiv1

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database"
)

type (
	Context struct {
		DB *gorm.DB
	}
	HandlerFunc func(context *Context, c echo.Context) error
	jsonmap     map[string]interface{}
)

var (
	messageEditableFields = []string{"Name", "IPAddress", "PCName"}
	groupEditableFields   = []string{"Name", "IPAddress", "PCName"}
	channelEditableFields = []string{"Name", "IPAddress", "PCName"}
)

func (ctx *Context) bind(f HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return f(ctx, c)
	}
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

func getClientUser(context *Context, c echo.Context) (database.User, error) {
	db := context.DB
	ip := c.RealIP()

	var user database.User
	if db.First(&user, "ip_address = ?", ip).Error != nil {
		return user, fmt.Errorf("user not found: %s")
	}
	return user, nil
}

/*
search: g\.(.*)\("(.*)", bindDB\((.*), db\)\)
replcae: // $3 $1 $2\nfunc $3(context *Context, c echo.Context) error {
  db := context.DB\n}\n
*/

// SetupRoutes setups routes
func SetupRoutes(g *echo.Group, db *gorm.DB) {
	c := Context{db}
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
	g.DELETE("/groups/:id/members/:id", c.bind(DeleteGroupMemberByID))
	g.GET("/groups/:id/channels", c.bind(GetGroupChannels))
	g.DELETE("/groups/:id/channels/:id", c.bind(DeleteGroupChannelByID))
	g.POST("/groups/:id/messages", c.bind(PostGroupMessages))

	g.POST("/channels", c.bind(PostChannels))
	g.GET("/channels", c.bind(GetChannels))
	g.PUT("/channels/:id", c.bind(PutChannelByID))
	g.GET("/channels/:id", c.bind(GetChannelByID))
	g.DELETE("/channels/:id", c.bind(DeleteChannelByID))
	g.GET("/channels/:id/members", c.bind(GetChannelMembers))
	g.DELETE("/channels/:id/members/:id", c.bind(DeleteChannelMemberByID))
}
