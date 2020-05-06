package apiv1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database"
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
