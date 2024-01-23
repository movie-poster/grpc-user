package utils

import (
	"net/url"
	"path"
	"strings"
)

func IsURL(content string) bool {
	_, err := url.ParseRequestURI(content)
	return err == nil && strings.HasPrefix(content, "https://")
}

func ExtractObjectName(fileURL string) string {
	urlPath, _ := url.ParseRequestURI(fileURL)
	return path.Base(urlPath.Path)
}
