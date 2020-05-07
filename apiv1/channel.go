package apiv1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database"
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

var (
	channelPostJsonFields   = []string{"Name", "GroupID"}
	channelResultJsonFields = []string{"ID", "Name", "GroupID"}
)

// PostChannels POST /channels
func PostChannels(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// GetChannels GET /channels
func GetChannels(context *Context, c echo.Context) error {
	db := context.DB
	var channels []database.Channel
	if db.Find(&channels).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "channel not found")
	}
	return c.JSON(http.StatusOK, channels)
}

// PutChannelByID PUT /channels/:id
func PutChannelByID(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// GetChannelByID GET /channels/:id
func GetChannelByID(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// DeleteChannelByID DELETE /channels/:id
func DeleteChannelByID(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// PostChannelMembers POST /channels/:id/members
func PostChannelMembers(context *Context, c echo.Context) error {
	// db := context.DB
	// user, err := getUserById(db, c)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusNotFound, err.Error())
	// }

	// if user.Channels == nil {
	// 	user.Channels = []*database.Channels{}
	// }

	// user.Channels = append(user.Channels, )

	// var channels []database.Channel
	// if err := db.Model(&user).Related(&channels).Error; err != nil {
	// 	return echo.NewHTTPError(http.StatusNotFound, err.Error())
	// }
	// return c.JSON(http.StatusOK, channels)
	return nil
}

// GetChannelMembers GET /channels/:id/members
func GetChannelMembers(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// DeleteChannelMemberByID DELETE /channels/:id/members/:id
func DeleteChannelMemberByID(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}
