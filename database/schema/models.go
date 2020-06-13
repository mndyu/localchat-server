package schema

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key;auto_increment:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

// All :
// 全スキーマ
var All = []interface{}{
	// テスト用
	Tee{},
	// 本体
	UserReadMessages{},
	UserGroupSettings{},
	MessageUsers{},
	User{},
	Notification{},
	Message{},
	Thread{},
	File{},
	Group{},
	Channel{},
}

// テスト用
type Tee struct {
	gorm.Model
	Name string
}

// 本体
type User struct {
	gorm.Model
	Name         string     `json:"name"`
	IPAddress    string     `json:"ip_address"`
	PCName       string     `json:"pc_name"`
	IsOnline     bool       `json:"is_online"`
	LastLogin    time.Time  `json:"last_login"`
	Messages     []Message  `json:"messages"`
	ReadMessages []*Message `gorm:"many2many:user_read_messages;"`
	Groups       []*Group   `gorm:"many2many:user_groups;"`
	Channels     []*Channel `gorm:"many2many:user_channels;"`
}
type Notification struct {
	gorm.Model
	Type        uint   `json:"type"`
	RecipientID uint   `json:"recipient_id"`
	SenderID    string `json:"sender_id"`
	URL         string `json:"url"`
	Read        bool   `json:"read"`
}

type Message struct {
	gorm.Model
	AuthorID  uint       `json:"author_id"`
	ChannelID *uint      `json:"channel_id"`
	GroupID   uint       `json:"group_id"`
	ThreadID  *uint      `json:"thread_id"`
	Body      string     `json:"body"`
	SentAt    time.Time  `json:"sent_at"`
	EditedAt  *time.Time `json:"edited_at"`
	To        []User     `json:"to" gorm:"many2many:message_users;"`
	ReadAt    []*User    `json:"read_at" gorm:"many2many:user_read_messages;"`
	Files     []*File    `gorm:"many2many:message_files;"`
	Author    User       `gorm:"foreignkey:AuthorID"`
}
type Thread struct {
	gorm.Model
	MessageID uint `json:"message_id"`
}
type Group struct {
	gorm.Model
	AuthorID uint      `json:"author_id"`
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

type File struct {
	gorm.Model
	AuthorID    uint   `json:"author_id"`
	MessageID   uint   `json:"message_id"`
	Size        uint   `json:"size"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FileName    string `json:"file_name"`
}

type UserReadMessages struct {
	UserID    uint      `json:"user_id" gorm:"primary_key;auto_increment:false"`
	MessageID uint      `json:"message_id" gorm:"primary_key;auto_increment:false"`
	Read      bool      `json:"read"`
	ReadAt    time.Time `json:"read_at"`
}
type UserGroupSettings struct {
	UserID  uint   `json:"user_id" gorm:"primary_key;auto_increment:false"`
	GroupID uint   `json:"group_id" gorm:"primary_key;auto_increment:false"`
	Body    string `json:"body"`
}
type MessageUsers struct {
	MessageID uint `json:"message_id" gorm:"primary_key;auto_increment:false"`
	UserID    uint `json:"user_id" gorm:"primary_key;auto_increment:false"`
}
