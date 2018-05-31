package util

import "strings"

func CheckCommand(s string) string {
	return strings.Split(s, " ")[0]
}
