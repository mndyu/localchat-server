package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

// All :
// 全スキーマ
var AllSchemas = []interface{}{
	Tee{},
	User{},
	Message{},
	Group{},
	Channel{},
}

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key;auto_increment:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

type Tee struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Name         string     `json:"name"`
	IPAddress    string     `json:"ip_address"`
	PCName       string     `json:"pc_name"`
	Messages     []Message  `json:"messages"`
	ReadMessages []*Message `json:"read_messages" gorm:"many2many:user_read_messages;"`
	Groups       []*Group   `json:"groups" gorm:"many2many:user_groups;"`
	Channels     []*Channel `json:"channels" gorm:"many2many:user_channels;"`
}
type Message struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	ChannelID uint    `json:"channel_id"`
	Text      string  `json:"text"`
	ReadAt    []*User `json:"read_at" gorm:"many2many:user_read_messages;"`
}

type Group struct {
	gorm.Model
	Name     string    `json:"name"`
	Members  []*User   `json:"members" gorm:"many2many:user_groups;"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
	gorm.Model
	Name     string    `json:"name"`
	GroupID  uint      `json:"group_id"`
	Members  []User    `json:"members" gorm:"many2many:user_channels;"`
	Messages []Message `json:"messages"`
	Group    Group
}
