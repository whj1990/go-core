package util

import (
	"fmt"
)

func GetLikeString(param string) string {
	return fmt.Sprintf("%%%s%%", param)
}

func GetSortString(sortField string, sortOrder int32) string {
	var suffix string
	if sortField == "" {
		return ""
	}
	switch sortOrder {
	case 1:
		suffix = "asc"
	case -1:
		suffix = "desc"
	}
	return fmt.Sprintf("%s %s", sortField, suffix)
}
