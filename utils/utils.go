package utils

import (
	"strconv"
	"bytes"
	"log"
	"regexp"
)

var errorCodeRegex *regexp.Regexp = regexp.MustCompile(`^Error (\d+):.*`)

func FatalErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func PrintErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}

func ConvertInt64ArrayToString(int64Array []int64) string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	for index, id := range int64Array {
		buffer.WriteString(strconv.FormatInt(id, 10))
		if index < len(int64Array) - 1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(")")
	return buffer.String()
}

func ParseSQLError(err error) string {
	match := errorCodeRegex.FindStringSubmatch(err.Error())
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
