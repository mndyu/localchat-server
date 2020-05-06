package apiv1

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database"
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

var (
	userResultJsonFields = []string{"ID", "Name", "IPAddress", "PCName"}
	userPostJsonFields   = []string{"Name", "IPAddress", "PCName"}
)

func getJsonFields() {

}

func getUserById(db *gorm.DB, c echo.Context) (database.User, error) {
	var user database.User

	id := c.Param("id")
	if id == "" {
		return user, fmt.Errorf("user param not found")
	}
	if db.First(&user, id).Error != nil {
		return user, fmt.Errorf("user not found")
	}
	return user, nil
}

// PostUsers POST /users
func PostUsers(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData userPostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	newUser := database.User{
		Name:      postData.Name,
		IPAddress: postData.IPAddress,
		PCName:    postData.PCName,
	}
	// newUser, ok := filterUserData(baseUser, newUser, userPostJsonFields).(database.User)
	if err := db.Create(&newUser).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user not found: %v", err))
	}

	// output
	jsonData := userResultJson{
		ID:        newUser.ID,
		Name:      newUser.Name,
		IPAddress: newUser.IPAddress,
		PCName:    newUser.PCName,
	}
	// jsonData := filterForJson(newUser, userResultJsonFields)
	return c.JSON(http.StatusOK, jsonData)
}

// GetUsers GET /users
func GetUsers(context *Context, c echo.Context) error {
	db := context.DB

	// db
	var users []database.User
	if db.Find(&users).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// output
	jsonData := []userResultJson{}
	for _, u := range users {
		jd := userResultJson{
			ID:        u.ID,
			Name:      u.Name,
			IPAddress: u.IPAddress,
			PCName:    u.PCName,
		}
		jsonData = append(jsonData, jd)
	}
	return c.JSON(http.StatusOK, jsonData)
}

// PutUserByID PUT /users/:id
func PutUserByID(context *Context, c echo.Context) error {
	db := context.DB

	// input
	var postData userPostJson
	if err := c.Bind(&postData); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// db
	user, err := getUserById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	newUser := database.User{
		Name:      postData.Name,
		IPAddress: postData.IPAddress,
		PCName:    postData.PCName,
	}
	newUser, ok := filterUserData(user, newUser, userPostJsonFields).(database.User)
	if !ok || db.Save(&user).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// output
	jsonData := userResultJson{
		ID:        newUser.ID,
		Name:      newUser.Name,
		IPAddress: newUser.IPAddress,
		PCName:    newUser.PCName,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// GetUserByID GET /users/:id
func GetUserByID(context *Context, c echo.Context) error {
	db := context.DB

	// db
	user, err := getUserById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// output
	jsonData := userResultJson{
		ID:        user.ID,
		Name:      user.Name,
		IPAddress: user.IPAddress,
		PCName:    user.PCName,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// DeleteUsersByID DELETE /users/:id
func DeleteUsersByID(context *Context, c echo.Context) error {
	db := context.DB

	// db
	user, err := getUserById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if db.Delete(&user).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user delete failed")
	}

	// output
	jsonData := userResultJson{
		ID:        user.ID,
		Name:      user.Name,
		IPAddress: user.IPAddress,
		PCName:    user.PCName,
	}
	return c.JSON(http.StatusOK, jsonData)
}

// PostUserMessages POST /users/:id/messages
func PostUserMessages(context *Context, c echo.Context) error {
	return nil
}

// GetUserMessages GET /users/:id/messages
func GetUserMessages(context *Context, c echo.Context) error {
	db := context.DB
	var msgs []database.Message
	id := c.Param("id")
	if db.Find(&msgs, id).Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}
	return c.JSON(http.StatusOK, msgs)
}

// PostUserGroups POST /users/:id/groups
func PostUserGroups(context *Context, c echo.Context) error {
	// db := context.DB
	return nil
}

// GetUserGroups GET /users/:id/groups
func GetUserGroups(context *Context, c echo.Context) error {
	db := context.DB
	user, err := getUserById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var groups []database.Group
	if err := db.Model(&user).Related(&groups).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, groups)
}

// PostUserChannels POST /users/:id/channels
func PostUserChannels(context *Context, c echo.Context) error {
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

// GetUserChannels GET /users/:id/channels
func GetUserChannels(context *Context, c echo.Context) error {
	db := context.DB
	user, err := getUserById(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	var channels []database.Channel
	if err := db.Model(&user).Related(&channels).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, channels)
}
