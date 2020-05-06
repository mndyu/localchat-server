package testutils

import (
	"encoding/json"
)

// EncodeObj :
// obj -> JSON
func EncodeObj(m interface{}) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DecodeJSON :
// JSON -> obj
func DecodeJSON(jsonString string) (interface{}, error) {
	var m interface{}
	b := []byte(jsonString)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// FormatJSON :
// JSON を再変換してフォーマット
func FormatJSON(jsonString string) (string, error) {
	m, err := DecodeJSON(jsonString)
	if err != nil {
		return "", err
	}
	return EncodeObj(m)
}

// RepackObj :
// obj を再変換してフォーマット
func RepackObj(m interface{}) (interface{}, error) {
	var r interface{}
	jstr, err := ToJSON(m)
	if err != nil {
		return r, err
	}
	r, err = DecodeJSON(jstr)
	if err != nil {
		return r, err
	}
	return r, nil
}

// ToJSON :
// obj -> JSON (フォーマット済み)
func ToJSON(m interface{}) (string, error) {
	s, err := EncodeObj(m)
	if err != nil {
		return "", err
	}
	return FormatJSON(s)
}
