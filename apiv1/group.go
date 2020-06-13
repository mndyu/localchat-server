package apiv1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database/schema"
)

type groupPostJson struct {
	Name string `json:"name"`
}
type groupResultJson struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type groupMemberPostJson struct {
	AuthorID uint `json:"user_id"`
	Myself   bool `json:"myself"`
}

// PostGroups POST /groups
func PostGroups(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData = jsonmap{}
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	// db
	var newItem schema.Group
	filteredPostData := filterJsonmapWithStruct(postData, groupPostJson{})
	assignJSONFields(&newItem, filteredPostData)
	if db.Create(&newItem).Error != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "group update failed:")
	}

	// output
	var jsonData groupResultJson
	assignJSONFields(&jsonData, newItem)
	return c.JSON(http.StatusOK, jsonData)
}

// GetGroups GET /groups
func GetGroups(context *Context, c echo.Context) error {
	db := context.DB

	// input
	limit := getLimit(c)
	offset := getOffset(c)

	// db
	var groups []schema.Group
	if db.Limit(limit).Offset(offset).Find(&groups).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}

	// output
	jsonData := []groupResultJson{}
	assignJSONArrayFields(&jsonData, groups)
	return c.JSON(http.StatusOK, jsonData)
}

// GetGroupByID GET /groups/:id
func GetGroupByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var group schema.Group
	if err := db.Find(&group, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData groupResultJson
	assignJSONFields(&jsonData, group)
	return c.JSON(http.StatusOK, jsonData)
}

// PutGroupByID PUT /groups/:id
func PutGroupByID(context *Context, c echo.Context) error {
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
	var group schema.Group
	if err := db.Find(&group, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	var newItem schema.Group
	filteredPostData := filterJsonmapWithStruct(postData, groupPostJson{})
	assignJSONFields(&newItem, filteredPostData)
	if err := db.Save(&newItem).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("message create found: %v", err))
	}

	// output
	var jsonData groupResultJson
	assignJSONFields(&jsonData, newItem)
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteGroupByID DELETE /groups/:id
func DeleteGroupByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var group schema.Group
	if err := db.Find(&group, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err := db.Delete(&group).Error; err != nil {
		echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("group delete failed: %s", err.Error()))
	}

	// output
	var jsonData groupResultJson
	assignJSONFields(&jsonData, group)
	return c.JSON(http.StatusOK, jsonData)
}

// PostGroupMembers GET /groups/:id/members
func PostGroupMembers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}
	var postData groupMemberPostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	// db
	var item schema.Group
	if err := db.Find(&item, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if item.Members == nil {
		item.Members = []*schema.User{}
	}
	var userID uint
	if postData.Myself {
		clientUser, err := getClientUser(context, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("unknown IP address: %s", c.RealIP()))
		}
		userID = clientUser.ID
	} else {
		userID = postData.AuthorID
	}
	var newMember schema.User
	if err := db.First(&newMember, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user not found: %d, %v", postData.AuthorID, postData.Myself))
	}
	item.Members = append(item.Members, &newMember)
	if err := db.Save(&item).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "failed updating group")
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, newMember)
	return c.JSON(http.StatusOK, jsonData)
}

// GetGroupMembers GET /groups/:id/members
func GetGroupMembers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var item schema.Group
	if err := db.Find(&item, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var members []schema.User
	if err := db.Model(&item).Related(&members, "Members").Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData []userResultJson
	assignJSONArrayFields(&jsonData, members)
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteGroupMemberByID DELETE /groups/:id/members/:id
func DeleteGroupMemberByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var item schema.Group
	if err := db.Find(&item, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var user schema.User
	if err := db.Find(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err := db.Model(&item).Association("Members").Delete(user).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, user)
	return c.JSON(http.StatusOK, jsonData)
}

// GetGroupChannels GET /groups/:id/channels
func GetGroupChannels(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var item schema.Group
	if err := db.Find(&item, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var channels []schema.Channel
	if err := db.Model(&item).Related(&channels, "Channels").Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData []channelResultJson
	assignJSONFields(&jsonData, channels)
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteGroupChannelByID DELETE /groups/:id/channels/:id
func DeleteGroupChannelByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}
	channelID, err := strconv.Atoi(c.Param("channel_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var item schema.Group
	if err := db.Find(&item, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var channel schema.Channel
	if err := db.Find(&channel, channelID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err := db.Model(&item).Association("Channels").Delete(channel).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData channelResultJson
	assignJSONFields(&jsonData, channel)
	return c.JSON(http.StatusOK, jsonData)
}

// PostGroupMessages POST /groups/:id/messages
func PostGroupMessages(context *Context, c echo.Context) error {
	db := context.DB

	clientUser, err := getClientUser(context, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("unknown IP address: %s", c.RealIP()))
	}

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}
	var postData messagePostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	var group schema.Group
	if err := db.Find(&group, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	var channel schema.Channel
	if group.Channels == nil {
		newChannel := schema.Channel{}
		if err := db.Create(&newChannel).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "failed creating default channel")
		}
		group.Channels = []schema.Channel{
			newChannel,
		}
		channel = newChannel
	} else {
		channel = group.Channels[0]
	}

	if channel.Messages == nil {
		channel.Messages = []schema.Message{}
	}
	newMessage := schema.Message{
		Body:      postData.Body,
		ChannelID: &channel.ID,
		AuthorID:  clientUser.ID,
	}
	channel.Messages = append(channel.Messages, newMessage)
	if err := db.Save(&channel).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "failed updating channel")
	}

	// output
	var jsonData messageResultJson
	assignJSONFields(&jsonData, newMessage)
	return c.JSON(http.StatusOK, jsonData)
}
