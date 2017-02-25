package utils

import (
	"strconv"
	"bytes"
)

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
