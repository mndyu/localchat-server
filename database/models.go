package database

import "github.com/jinzhu/gorm"

// All :
// 全スキーマ
var All = []interface{}{
	Tee{},
	User{},
	Message{},
	Group{},
	Channel{},
}

type Tee struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Name         string     `json:"name"`
	Messages     []Message  `json:"messages"`
	ReadMessages []*Message `json:"read_messages" gorm:"many2many:user_read_messages;"`
	Groups       []*Group   `json:"groups" gorm:"many2many:user_groups;"`
	Channels     []*Channel `json:"channels" gorm:"many2many:user_channels;"`
	IPAddress    string     `json:"ip_address"`
	PCName       string     `json:"pc_name"`
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
}
