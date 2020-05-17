package apiv1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PostUsers POST /users
func TestPostUsers(t *testing.T) {
}

// GetUsers GET /users
func TestGetUsers(t *testing.T) {
	db := getNewMockDB()
	api := testapi{
		context: &Context{db},
	}

	//
	res, err := api.Get(GetUsers, nil, nil)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	res.assertStatus(t, 200)
	res.assertBody(t, []interface{}{jsonmap{
		"name": "unkunkdo",
	}})

	b, _ := json.Marshal(res.body)
	fmt.Println("data: ", string(b))
}

// PutUserByID PUT /users/:id
func TestPutUserByID(t *testing.T) {

}

// GetUserByID GET /users/:id
func TestGetUserByID(t *testing.T) {
}

// DeleteUsersByID DELETE /users/:id
func TestDeleteUsersByID(t *testing.T) {
}

// PostUserMessages POST /users/:id/messages
func TestPostUserMessages(t *testing.T) {
}

// GetUserMessages GET /users/:id/messages
func TestGetUserMessages(t *testing.T) {
}

// PostUserGroups POST /users/:id/groups
func TestPostUserGroups(t *testing.T) {
}

// GetUserGroups GET /users/:id/groups
func TestGetUserGroups(t *testing.T) {
}

// PostUserChannels POST /users/:id/channels
func TestPostUserChannels(t *testing.T) {
}

// GetUserChannels GET /users/:id/channels
func TestGetUserChannels(t *testing.T) {
	db := getNewMockDB()
	api := testapi{
		context: &Context{db},
	}

	// cards := []database.Card{
	// 	database.Card{ID: 2},
	// }
	// db.Create()

	res, err := api.Get(GetUserChannels, "id=2", "name=2&id=23")
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	res.assertStatus(t, 200)
}
