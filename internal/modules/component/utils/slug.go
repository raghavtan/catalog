package utils

import "strings"

func GetSlug(name, componentType string) string {
	var shortType string
	switch strings.ToUpper(componentType) {
	case "SERVICE":
		shortType = "svc"
	case "CLOUD_RESOURCE":
		shortType = "cr"
	case "WEBSITE":
		shortType = "web"
	case "APPLICATION":
		shortType = "app"
	default:
		shortType = "unknown"
	}
	return shortType + "-" + name
}
