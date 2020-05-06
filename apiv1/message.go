package apiv1

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database"
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

var (
	messageResultJsonFields = []string{"ID", "UserID", "ChannelID", "Text"}
	messagePostJsonFields   = []string{"UserID", "ChannelID", "Text"}
)

func getMessageById(db *gorm.DB, c echo.Context) (database.Message, error) {
	var msg database.Message

	id := c.Param("id")
	if id == "" {
		return msg, fmt.Errorf("msg param not found")
	}
	if db.First(&msg, id).Error != nil {
		return msg, fmt.Errorf("msg not found")
	}
	return msg, nil
}

// PostMessages POST /messages
func PostMessages(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData messagePostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}

	// db
	newItem := database.Message{
		Text: postData.Text,
	}
	// newUser, ok := filterUserData(baseUser, newUser, userPostJsonFields).(database.User)
	if err := db.Create(&newItem).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("message create found: %v", err))
	}

	// output
	jsonData := messageResultJson{
		ID:        newItem.ID,
		Text:      newItem.Text,
		UserID:    newItem.UserID,
		ChannelID: newItem.ChannelID,
	}
	// jsonData := filterForJson(newUser, userResultJsonFields)
	return c.JSON(http.StatusOK, jsonData)
}

// GetMessages GET /messages
func GetMessages(context *Context, c echo.Context) error {
	db := context.DB

	// db
	var msgs []database.Message
	if db.Find(&msgs).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}

	// output
	jsonData := []messageResultJson{}
	for _, u := range msgs {
		jd := messageResultJson{
			ID:        u.ID,
			Text:      u.Text,
			UserID:    u.UserID,
			ChannelID: u.ChannelID,
		}
		jsonData = append(jsonData, jd)
	}
	return c.JSON(http.StatusOK, jsonData)
}

// GetMessageByID GET /messages/:id
func GetMessageByID(context *Context, c echo.Context) error {
	db := context.DB

	// db
	msg, err := getMessageById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	jsonData := messageResultJson{
		ID:        msg.ID,
		Text:      msg.Text,
		UserID:    msg.UserID,
		ChannelID: msg.ChannelID,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// PutMessageByID PUT /messages/:id
func PutMessageByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData messagePostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	item, err := getMessageById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	newItem := database.Message{
		Text:      item.Text,
		UserID:    item.UserID,
		ChannelID: item.ChannelID,
	}
	newItem, ok := filterUserData(item, newItem, messagePostJsonFields).(database.Message)
	if !ok || db.Save(&item).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}

	// output
	jsonData := messageResultJson{
		ID:        item.ID,
		Text:      item.Text,
		UserID:    item.UserID,
		ChannelID: item.ChannelID,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteMessageByID DELETE /messages/:id
func DeleteMessageByID(context *Context, c echo.Context) error {
	db := context.DB

	// db
	msg, err := getMessageById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if db.Delete(&msg).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "msg delete failed")
	}

	// output
	jsonData := messageResultJson{
		ID:        msg.ID,
		Text:      msg.Text,
		UserID:    msg.UserID,
		ChannelID: msg.ChannelID,
	}
	return c.JSON(http.StatusOK, jsonData)
}
