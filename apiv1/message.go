package apiv1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database/schema"
)

type messagePostJson struct {
	Text string `json:"text"`
}

type messageResultJson struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ChannelID uint   `json:"channel_id"`
	Text      string `json:"text"`
}

// PostMessages POST /messages
func PostMessages(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData = jsonmap{}
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}

	// db
	var newItem schema.Message
	filteredPostData := filterJsonmapWithStruct(postData, messagePostJson{})
	assignJSONFields(&newItem, filteredPostData)
	if err := db.Create(&newItem).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("message create found: %v", err))
	}

	// output
	var jsonData messageResultJson
	assignJSONFields(&jsonData, newItem)
	return c.JSON(http.StatusOK, jsonData)
}

// GetMessages GET /messages
func GetMessages(context *Context, c echo.Context) error {
	db := context.DB

	// input
	limit := getLimit(c)
	offset := getOffset(c)

	// db
	var msgs []schema.Message
	if db.Limit(limit).Offset(offset).Find(&msgs).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}

	// output
	jsonData := []messageResultJson{}
	assignJSONArrayFields(&jsonData, msgs)
	return c.JSON(http.StatusOK, jsonData)
}

// GetMessageByID GET /messages/:id
func GetMessageByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var msg schema.Message
	if err := db.Find(&msg, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData messageResultJson
	assignJSONFields(&jsonData, msg)
	return c.JSON(http.StatusOK, jsonData)
}

// PutMessageByID PUT /messages/:id
func PutMessageByID(context *Context, c echo.Context) error {
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
	var msg schema.Message
	if err := db.Find(&msg, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	filteredPostData := filterJsonmapWithStruct(postData, messagePostJson{})
	assignJSONFields(&msg, filteredPostData)
	if err := db.Save(&msg).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user not found: %s", err.Error()))
	}

	// output
	var jsonData messageResultJson
	assignJSONFields(&jsonData, msg)
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteMessageByID DELETE /messages/:id
func DeleteMessageByID(context *Context, c echo.Context) error {
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
	var msg schema.Message
	if err := db.Find(&msg, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if db.Delete(&msg).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "msg delete failed")
	}

	// output
	var jsonData messageResultJson
	assignJSONFields(&jsonData, msg)
	return c.JSON(http.StatusOK, jsonData)
}
