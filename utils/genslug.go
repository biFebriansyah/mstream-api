package utils

import "regexp"

func Slug(name string) string {
	spaceRemove := regexp.MustCompile(`\s+`).ReplaceAll([]byte(name), []byte("-"))
	nonWordRemove := regexp.MustCompile(`[^\w-]+`).ReplaceAll(spaceRemove, []byte("-"))
	specialCharRemove := regexp.MustCompile(`-+`).ReplaceAll(nonWordRemove, []byte("-"))
	return string(specialCharRemove)
}
