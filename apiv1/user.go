package apiv1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database/schema"
)

type userPostJson struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	PCName    string `json:"pc_name"`
}

type userResultJson struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	PCName    string `json:"pc_name"`
}

// GetProfile GET /profile
func GetProfile(context *Context, c echo.Context) error {
	// db := context.DB

	// input
	user, err := getClientUser(context, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("auth failed for IP address (%s)", c.RealIP()))
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, user)
	return c.JSON(http.StatusOK, jsonData)
}

// GetUsers GET /users
func GetUsers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	limit := getLimit(c)
	offset := getOffset(c)

	// db
	var users []schema.User
	if db.Limit(limit).Offset(offset).Find(&users).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// output
	jsonData := []userResultJson{}
	assignJSONArrayFields(&jsonData, users)
	return c.JSON(http.StatusOK, jsonData)
}

// PostUsers POST /users
func PostUsers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData = jsonmap{}
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	var newUser schema.User
	filteredPostData := filterJsonmapWithStruct(postData, userPostJson{})
	assignJSONFields(&newUser, filteredPostData)
	if err := db.Create(&newUser).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user not found: %v", err))
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, newUser)
	return c.JSON(http.StatusOK, jsonData)
}

// PutUserByID PUT /users/:id
func PutUserByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}
	var postData = jsonmap{}
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	var user schema.User
	if err := db.Find(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	filteredPostData := filterJsonmapWithStruct(postData, userPostJson{})
	assignJSONFields(&user, filteredPostData)
	if err := db.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user not found: %s", err.Error()))
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, user)
	return c.JSON(http.StatusOK, jsonData)
}

// GetUserByID GET /users/:id
func GetUserByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var user schema.User
	if err := db.Find(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, user)
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteUsersByID DELETE /users/:id
func DeleteUsersByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var user schema.User
	if err := db.Find(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if db.Delete(&user).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user delete failed")
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, user)
	return c.JSON(http.StatusOK, jsonData)
}

// GetUserMessages GET /users/:id/messages
func GetUserMessages(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var user schema.User
	if err := db.Find(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var msgs []schema.Message
	if err := db.Model(&user).Related(&msgs, "Messages").Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	jsonData := []messageResultJson{}
	assignJSONArrayFields(&jsonData, msgs)
	return c.JSON(http.StatusOK, jsonData)
}

// GetUserGroups GET /users/:id/groups
func GetUserGroups(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var user schema.User
	if err := db.Find(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var groups []schema.Group
	if err := db.Model(&user).Related(&groups, "Groups").Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	jsonData := []groupResultJson{}
	assignJSONArrayFields(&jsonData, groups)
	return c.JSON(http.StatusOK, jsonData)
}

// GetUserChannels GET /users/:id/channels
func GetUserChannels(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var user schema.User
	if err := db.Find(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var channels []schema.Channel
	if err := db.Model(&user).Related(&channels, "Channels").Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	jsonData := []channelResultJson{}
	assignJSONArrayFields(&jsonData, channels)
	return c.JSON(http.StatusOK, jsonData)
}
