package utils

import "strconv"

func StrToI64(str string) int64 {
	if len(str) <= 0 {
		return 0
	}
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return res
}

func I64ToStr(integer int64) string {
	return strconv.FormatInt(integer, 10)
}

func IToStr(integer int) string {
	return strconv.Itoa(integer)
}

func StrToI(str string) int {
	if len(str) <= 0 {
		return 0
	}
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return res
}
