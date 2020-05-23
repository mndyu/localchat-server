package apiv1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database/schema"
)

type channelPostJson struct {
	Name    string `json:"name"`
	GroupID uint   `json:"group_id"`
}

type channelResultJson struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	GroupID uint   `json:"group_id"`
}

type channelMemberPostJson struct {
	UserID uint `json:"user_id"`
	Myself bool `json:"myself"`
}

// PostChannels POST /channels
func PostChannels(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData jsonmap
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	// db
	var newItem schema.Channel
	filteredPostData := filterJsonmapWithStruct(postData, channelPostJson{})
	assignJSONFields(&newItem, filteredPostData)
	if db.Create(&newItem).Error != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "channel update failed:")
	}

	// output
	var jsonData channelResultJson
	assignJSONFields(&jsonData, newItem)
	return c.JSON(http.StatusOK, jsonData)
}

// GetChannels GET /channels
func GetChannels(context *Context, c echo.Context) error {
	db := context.DB

	// input
	limit := getLimit(c)
	offset := getOffset(c)

	// db
	var channels []schema.Channel
	if db.Limit(limit).Offset(offset).Find(&channels).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "channels not found")
	}

	// output
	jsonData := []channelResultJson{}
	assignJSONArrayFields(&jsonData, channels)
	return c.JSON(http.StatusOK, jsonData)
}

// PutChannelByID PUT /channels/:id
func PutChannelByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}
	var postData jsonmap
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	var channel schema.Channel
	if err := db.Find(&channel, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	var newItem schema.Channel
	filteredPostData := filterJsonmapWithStruct(postData, channelPostJson{})
	assignJSONFields(&newItem, filteredPostData)
	if err := db.Save(&newItem).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("message create found: %v", err))
	}

	// output
	var jsonData groupResultJson
	assignJSONFields(&jsonData, channel)
	return c.JSON(http.StatusOK, jsonData)
}

// GetChannelByID GET /channels/:id
func GetChannelByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var channel schema.Channel
	if err := db.Find(&channel, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData groupResultJson
	assignJSONFields(&jsonData, channel)
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteChannelByID DELETE /channels/:id
func DeleteChannelByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var channel schema.Channel
	if err := db.Find(&channel, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err := db.Delete(&channel).Error; err != nil {
		echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("channel delete failed: %s", err.Error()))
	}

	// output
	var jsonData groupResultJson
	assignJSONFields(&jsonData, channel)
	return c.JSON(http.StatusOK, jsonData)
}

// PostChannelMembers POST /channels/:id/members
func PostChannelMembers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}
	var postData channelMemberPostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	// db
	var item schema.Channel
	if err := db.Find(&item, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if item.Members == nil {
		item.Members = []schema.User{}
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
	var newMember schema.User
	if err := db.First(&newMember, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user not found: %d, %v", postData.UserID, postData.Myself))
	}
	item.Members = append(item.Members, newMember)
	if err := db.Save(&item).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "failed updating group")
	}

	// output
	var jsonData userResultJson
	assignJSONFields(&jsonData, newMember)
	return c.JSON(http.StatusOK, jsonData)
}

// GetChannelMembers GET /channels/:id/members
func GetChannelMembers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", c.Param("id")))
	}

	// db
	var item schema.Channel
	if err := db.Find(&item, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var members []schema.User
	if err := db.Model(&item).Related(&members, "Members").Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	var jsonData []userResultJson
	assignJSONFields(&jsonData, members)
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteChannelMemberByID DELETE /channels/:id/members/:id
func DeleteChannelMemberByID(context *Context, c echo.Context) error {
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
	var item schema.Channel
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
