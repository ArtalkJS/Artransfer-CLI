package lib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/araddon/dateparse"
	log "github.com/sirupsen/logrus"
)

// 任何类型转 String
//	(bool) true => (string) "true"
//	(int) 0 => (string) "0"
func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func ParseVersion(s string) int64 {
	strList := strings.Split(s, ".")
	format := fmt.Sprintf("%%s%%0%ds", len(strList))
	v := ""
	for _, value := range strList {
		v = fmt.Sprintf(format, v, value)
	}
	var result int64
	var err error
	if result, err = strconv.ParseInt(v, 10, 64); err != nil {
		log.Error("ugh: parseVersion(%s): error=%s", s, err)
		return 0
	}
	return result
}

func ParseDate(s string) time.Time {
	denverLoc, _ := time.LoadLocation("Local") // 时区
	time.Local = denverLoc
	t, _ := dateparse.ParseIn(s, denverLoc)

	return t
}

func PrintEncodeData(dataType string, val interface{}) {
	print(SprintEncodeData(dataType, val))
}

func SprintEncodeData(dataType string, val interface{}) string {
	return fmt.Sprintf("[%s]\n\n   %#v\n\n", dataType, val)
}

func PrintTable(rows [][]interface{}) {
	println("-------------------------")
	for _, row := range rows {
		l := len(row)
		print(" + ")
		for i, col := range row {
			print(col)
			if i < l-1 {
				print(": ")
			}
		}
		println()
	}
	println("-------------------------")
}

func HideJsonLongText(key string, text string) string {
	r := regexp.MustCompile(key + `:"(.+?)"`)
	sm := r.FindStringSubmatch(text)
	postText := ""
	if len(sm) > 0 {
		postText = sm[1]
	}

	text = r.ReplaceAllString(text, fmt.Sprintf(key+": <!-- 省略 %d 个字符 -->", utf8.RuneCountInString(postText)))
	return text
}
