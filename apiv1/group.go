package apiv1

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database"
)

type groupPostJson struct {
	Name string `json:"name"`
}
type groupResultJson struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type groupMemberPostJson struct {
	UserID uint `json:"user_id"`
	Myself bool `json:"myself"`
}

var (
	groupPostJsonFields   = []string{"Name"}
	groupResultJsonFields = []string{"ID", "Name"}
)

func getGroupById(db *gorm.DB, c echo.Context) (database.Group, error) {
	var group database.Group

	id := c.Param("id")
	if id == "" {
		return group, fmt.Errorf("group param not found")
	}
	if db.First(&group, id).Error != nil {
		return group, fmt.Errorf("group not found")
	}
	return group, nil
}

// PostGroups POST /groups
func PostGroups(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData groupPostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	// db
	newItem := database.Group{
		Name: postData.Name,
	}
	if db.Create(&newItem).Error != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "group update failed:")
	}

	// output
	jsonData := groupResultJson{
		ID:   newItem.ID,
		Name: newItem.Name,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// GetGroups GET /groups
func GetGroups(context *Context, c echo.Context) error {
	db := context.DB

	// db
	var groups []database.Group
	if db.Find(&groups).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	// output
	jsonData := []groupResultJson{}
	for _, item := range jsonData {
		jd := groupResultJson{
			ID:   item.ID,
			Name: item.Name,
		}
		jsonData = append(jsonData, jd)
	}
	return c.JSON(http.StatusOK, jsonData)
}

// GetGroupByID GET /groups/:id
func GetGroupByID(context *Context, c echo.Context) error {
	db := context.DB

	group, err := getGroupById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	jsonData := groupResultJson{
		ID:   group.ID,
		Name: group.Name,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// PutGroupByID PUT /groups/:id
func PutGroupByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData messagePostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	item, err := getGroupById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	newItem := database.Group{
		Name: item.Name,
	}
	newItem, ok := filterUserData(item, newItem, groupPostJsonFields).(database.Group)
	if !ok || db.Save(&newItem).Error != nil {
		echo.NewHTTPError(http.StatusNotFound, "group update failed:")
	}

	// output
	jsonData := groupResultJson{
		ID:   newItem.ID,
		Name: newItem.Name,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteGroupByID DELETE /groups/:id
func DeleteGroupByID(context *Context, c echo.Context) error {
	db := context.DB

	group, err := getGroupById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err := db.Delete(&group).Error; err != nil {
		echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("group delete failed: %s", err.Error()))
	}

	// output
	jsonData := groupResultJson{
		ID:   group.ID,
		Name: group.Name,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// PostGroupMembers GET /groups/:id/members
func PostGroupMembers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData groupMemberPostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	// db
	item, err := getGroupById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if item.Members == nil {
		item.Members = []*database.User{}
	}
	var userID uint
	if postData.Myself {
		clientUser, err := getClientUser(context, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("unknown IP address: %s", c.RealIP()))
		}
		userID = clientUser.ID
	} else {
		userID = postData.UserID
	}
	var newMember database.User
	if err := db.First(&newMember, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user not found: %d, %v", postData.UserID, postData.Myself))
	}
	item.Members = append(item.Members, &newMember)
	if err := db.Save(&item).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "failed updating group")
	}

	// output
	return c.JSON(http.StatusOK, postData)
}

// GetGroupMembers GET /groups/:id/members
func GetGroupMembers(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// DeleteGroupMemberByID DELETE /groups/:id/members/:id
func DeleteGroupMemberByID(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// GetGroupChannels GET /groups/:id/channels
func GetGroupChannels(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// DeleteGroupChannelByID DELETE /groups/:id/channels/:id
func DeleteGroupChannelByID(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// PostGroupMessages POST /groups/:id/messages
func PostGroupMessages(context *Context, c echo.Context) error {
	db := context.DB

	clientUser, err := getClientUser(context, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("unknown IP address: %s", c.RealIP()))
	}

	// input
	var postData messagePostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	group, err := getGroupById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	var channel database.Channel
	if group.Channels == nil {
		newChannel := database.Channel{}
		if err := db.Create(&newChannel).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "failed creating default channel")
		}
		group.Channels = []database.Channel{
			newChannel,
		}
		channel = newChannel
	} else {
		channel = group.Channels[0]
	}

	if channel.Messages == nil {
		channel.Messages = []database.Message{}
	}
	newMessage := database.Message{
		Text:      postData.Text,
		ChannelID: channel.ID,
		UserID:    clientUser.ID,
	}
	channel.Messages = append(channel.Messages, newMessage)
	if err := db.Save(&channel).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "failed updating channel")
	}

	// output
	jsonData := messageResultJson{
		ID:        newMessage.ID,
		UserID:    newMessage.UserID,
		ChannelID: newMessage.ChannelID,
		Text:      newMessage.Text,
	}
	return c.JSON(http.StatusOK, jsonData)
}
