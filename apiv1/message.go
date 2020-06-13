package apiv1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mndyu/localchat-server/database/schema"
	"github.com/mndyu/localchat-server/utils"
)

type messagePostJson struct {
	GroupID uint   `json:"group_id"`
	To      []uint `json:"to"`
	Body    string `json:"body"`
}

type messageResultJson struct {
	ID     uint `json:"id" gorm:"primary_key"`
	Author struct {
		ID        uint   `json:"id" gorm:"primary_key"`
		Name      string `json:"name"`
		IPAddress string `json:"ip_address"`
		PCName    string `json:"pc_name"`
	} `json:"author"`
	ThreadID *uint      `json:"thread"`
	GroupID  uint       `json:"group_id"`
	Body     string     `json:"body"`
	SentAt   time.Time  `json:"sent_at"`
	EditedAt *time.Time `json:"edited_at"`
}

// PostMessages POST /messages
func PostMessages(context *Context, c echo.Context) error {
	db := context.DB

	// input
	user, err := getClientUser(context, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("auth failed for IP address (%s)", c.RealIP()))
	}

	// var postData = jsonmap{}
	// if err := c.Bind(&postData); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "invalid request params")
	// }
	var filteredPostData messagePostJson
	if err := c.Bind(&filteredPostData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request params: %s", err.Error()))
	}
	// filteredPostData := filterJsonmapWithStruct(postData, messagePostJson{})

	if filteredPostData.To == nil || len(filteredPostData.To) == 0 || len(filteredPostData.To) > 10 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request params: to")
	}

	// db
	var newItem = schema.Message{
		AuthorID: user.ID,
		GroupID:  filteredPostData.GroupID,
		Body:     filteredPostData.Body,
		SentAt:   time.Now(),
	}
	if err := db.Create(&newItem).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("message create found: %v", err))
	}

	for _, u := range filteredPostData.To {
		n := schema.MessageUsers{
			MessageID: newItem.ID,
			UserID:    u,
		}
		if err := db.Create(&n).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("message create found: %v", err))
		}
	}

	// output
	// rows, err := db.Model(schema.Message{}).
	// 	Joins("join user on message.author_id = user.id").
	// 	Select("message.id, user.id, user.name, user.ip_address, user.pc_name, message.thread_id, message.group_id, message.body, message.sent_at, message.edited_at").
	// 	Rows()
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request params: %s", err.Error()))
	// }
	// var jsonData = []messageResultJson{}
	// err = utils.MapRows(&jsonData, rows)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request params: %s", err.Error()))
	// }
	// if len(jsonData) != 1 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request params: %s", err.Error()))
	// }
	// return c.JSON(http.StatusOK, jsonData[0])

	var us schema.User
	db.Model(newItem).Related(&us, "Author")

	var jsonData messageResultJson
	assignJSONFields(&jsonData, newItem)
	assignJSONFields(&jsonData.Author, us)

	return c.JSON(http.StatusOK, jsonData)
}

// GetMessages GET /messages
func GetMessages(context *Context, c echo.Context) error {
	db := context.DB

	// input
	// limit := getLimit(c)
	// offset := getOffset(c)

	// db
	// var msgs []schema.Message
	// if db.Limit(limit).Offset(offset).Find(&msgs).Error != nil {
	// 	return echo.NewHTTPError(http.StatusNotFound, "message not found")
	// }

	rows, err := db.Model(schema.Message{}).
		Joins("join user on message.author_id = user.id").
		Select("message.id, user.id, user.name, user.ip_address, user.pc_name, message.thread_id, message.group_id, message.body, message.sent_at, message.edited_at").
		Rows()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request params: %s", err.Error()))
	}
	var jsonData = []messageResultJson{}
	err = utils.MapRows(&jsonData, rows)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request params: %s", err.Error()))
	}

	// output
	// jsonData := []messageResultJson{}
	// assignJSONArrayFields(&jsonData, msgs)
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
	var us schema.User
	db.Model(msg).Related(&us, "Author")

	var jsonData messageResultJson
	assignJSONFields(&jsonData, msg)
	assignJSONFields(&jsonData.Author, us)
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
	var us schema.User
	db.Model(msg).Related(&us, "Author")

	var jsonData messageResultJson
	assignJSONFields(&jsonData, msg)
	assignJSONFields(&jsonData.Author, us)
	return c.JSON(http.StatusOK, jsonData)
}
