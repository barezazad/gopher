package helper

import (
	"strconv"
	"strings"
	"time"
)

func StrToInt(value string) (result int, err error) {
	result, err = strconv.Atoi(value)
	return
}

func IntToStr(value int) (result string) {
	result = strconv.Itoa(value)
	return
}

func UintToStr(value uint) (result string) {
	result = strconv.FormatUint(uint64(value), 10)
	return
}

func StrToUint(value string) (result uint, err error) {
	tmpID, err := strconv.ParseUint(value, 10, 64)
	result = uint(tmpID)
	return
}

func StrToUint64(value string) (result uint64, err error) {
	result, err = strconv.ParseUint(value, 10, 64)
	return
}

func StrToFloat(value string) (result float64, err error) {
	result, err = strconv.ParseFloat(value, 64)
	return
}

func IntToPointer(value int) *int {
	return &value
}

func StrToByte(value string) []byte {
	return []byte(value)
}

func ByteToStr(value []byte) string {
	return string(value)
}

func StrToDuration(value string) time.Duration {
	num, _ := StrToUint64(value)
	return time.Duration(num)
}

func StrToBool(value string) bool {
	switch strings.ToLower(value) {
	case "0", "false":
		return false
	case "1", "true":
		return true
	default:
		return false
	}
}
