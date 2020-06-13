package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// case 変換のパターン
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchUnderscore = regexp.MustCompile("_([A-Za-z])")

// ToSnakeCase :
// camel case -> snake case に変換
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToCamelCase :
// snake case -> lower camel case に変換
func ToCamelCase(str string) string {
	if str[0] == '_' {
		str = str[1:]
	}
	camel := matchUnderscore.ReplaceAllStringFunc(str, func(m string) string {
		return strings.ToUpper(m)
	})
	camel = matchUnderscore.ReplaceAllString(camel, "${1}")
	return camel
}

// ToUpperCamelCase :
// snake case -> upper camel case に変換
func ToUpperCamelCase(str string) string {
	str = strings.ToUpper(string(str[0])) + str[1:]
	camel := matchUnderscore.ReplaceAllStringFunc(str, func(m string) string {
		return strings.ToUpper(m)
	})
	camel = matchUnderscore.ReplaceAllString(camel, "${1}")
	return camel
}

// MeasureTime : 処理時間計測
func MeasureTime(f func()) {
	start := time.Now()

	f()

	elapsed := time.Since(start)
	fmt.Printf("time: %s\n", elapsed)
}
