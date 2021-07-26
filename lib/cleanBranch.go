package lib

import (
	"regexp"
	"strings"
)

func CleanBranch(name string) string {
	matchSlash := regexp.MustCompile(`/`)
	matchSpecial := regexp.MustCompile(`[^\w\s-]`)
	res := matchSlash.ReplaceAllString(name, "-")
	res = matchSpecial.ReplaceAllString(res, "")
	res = strings.ToLower(res)
	return res
}
