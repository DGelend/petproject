package controller

import (
	"strconv"
	"strings"
	//"fmt"
	//"log"
)

func checkFormatPhone(s string) bool {
	if !(strings.HasPrefix(s, "79") && len([]rune(s)) == 11) {
		return false
	}
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true

}
