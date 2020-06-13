package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func createIMDB(tables []interface{}, rows []interface{}) (*gorm.DB, error) {
	// db connection
	db, err := gorm.Open("sqlite3", ":memory:") // DBMS
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// create tables
	for _, t := range tables {
		db.AutoMigrate(t)
	}

	// insert rows
	for _, r := range rows {
		db.Create(r)
	}
	return db, nil
}

func toJSONStruct(a interface{}) (interface{}, error) {
	if _, ok := a.(string); ok == false {
		b, err := json.Marshal(&a)
		if err != nil {
			return nil, err
		}
		a = string(b)
	}

	var m interface{}
	err := json.Unmarshal([]byte(a.(string)), &m)
	return m, err
}
func isEqual(a interface{}, b interface{}) (bool, error) {
	m1, err := toJSONStruct(a)
	if err != nil {
		return false, err
	}
	m2, err := toJSONStruct(b)
	if err != nil {
		return false, err
	}
	return cmp.Equal(m1, m2), nil
}
func printJSON(a interface{}) {
	b, err := json.Marshal(a)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println(string(b))
}

type messageResultJson struct {
	ID   uint `json:"id" gorm:"primary_key"`
	User struct {
		ID        uint   `json:"id" gorm:"primary_key"`
		Name      string `json:"name"`
		IPAddress string `json:"ip_address"`
		PCName    string `json:"pc_name"`
	} `json:"user"`
	ThreadID *uint      `json:"thread"`
	GroupID  uint       `json:"group_id"`
	Body     string     `json:"body"`
	SentAt   time.Time  `json:"sent_at"`
	EditedAt *time.Time `json:"edited_at"`
}

func TestEachQueryStructField(t *testing.T) {
	a := uint(9982)
	te := time.Now()
	u := messageResultJson{
		123,
		struct {
			ID        uint   `json:"id" gorm:"primary_key"`
			Name      string `json:"name"`
			IPAddress string `json:"ip_address"`
			PCName    string `json:"pc_name"`
		}{333, "ji", "ddd", "apap"},
		&a,
		1323,
		"esofse",
		time.Now(),
		&te,
	}
	fields := []reflect.Value{}
	EachQueryStructField(&u, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		fields = append(fields, val)
		// val.Set(reflect.Zero(field.Type))
	})
	for i, f := range fields {
		fmt.Println(i, f, f.Type())
	}
	fmt.Println(u)
	t.Error(0)
}

func TestGetTypedColumns(t *testing.T) {
	u := []messageResultJson{}
	var columns []interface{}
	GetTypedColumns(&u, &columns, nil)
	// t.Error(columns)
}

func TestCreateOrUpdateRow(t *testing.T) {
	type User struct {
		ID      uint `gorm:"primary_key"`
		Name    string
		Desc    string
		Profile string
		Heh     *bool
	}
	type Message struct {
		ID     uint `gorm:"primary_key"`
		Title  string
		Body   string
		UserID uint
	}
	tables := []interface{}{
		User{},
		Message{},
	}
	boolean := true
	rows := []interface{}{
		User{1, "Confucius", "Han", "smart boi", nil},
		User{2, "Ogedei", "Mongolian", "warlord boi", nil},
		User{3, "Panini", "Aryan", "linguist boi", &boolean},
		Message{1, "taitor", "bodeli-", 1},
		Message{2, "haahe", "fef-", 1},
		Message{3, "ooopoe", "hage", 3},
	}

	// db connection
	db, err := createIMDB(tables, rows) // DBMS
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//
	var testCases = []struct {
		query       string
		queryStruct interface{}
		expected    string
	}{
		{
			"select * from users left join messages on users.id = messages.user_id",
			&[]struct {
				User
				Messages []Message
			}{},
			`[{"ID":1,"Name":"Confucius","Desc":"Han","Profile":"smart boi","Heh":null,"Messages":[{"ID":2,"Title":"haahe","Body":"fef-","UserID":1},{"ID":1,"Title":"taitor","Body":"bodeli-","UserID":1}]},{"ID":2,"Name":"Ogedei","Desc":"Mongolian","Profile":"warlord boi","Heh":null,"Messages":null},{"ID":3,"Name":"Panini","Desc":"Aryan","Profile":"linguist boi","Heh":true,"Messages":[{"ID":3,"Title":"ooopoe","Body":"hage","UserID":3}]}]`,
		},
	}

	for _, tc := range testCases {
		queryRows, err := db.Raw(tc.query).Rows()
		if err != nil {
			t.Fatal(err)
		}
		err = MapRows(tc.queryStruct, queryRows)
		if err != nil {
			t.Fatal(err)
		}

		e, err := isEqual(tc.queryStruct, tc.expected)
		if err != nil {
			t.Fatal(err)
		}
		if !e {
			t.Errorf("not equal: %v != %v", tc.queryStruct, tc.expected)
		}
		printJSON(tc.queryStruct)
	}
}
