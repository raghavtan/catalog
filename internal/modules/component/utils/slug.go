package utils

import "strings"

func GetSlug(name, componentType string) string {
	var shortType string
	switch strings.ToLower(componentType) {
	case "service":
		shortType = "svc"
	case "cloud-resource":
		shortType = "cr"
	case "website":
		shortType = "web"
	case "application":
		shortType = "app"
	default:
		shortType = "unknown"
	}
	return shortType + "-" + name
}
