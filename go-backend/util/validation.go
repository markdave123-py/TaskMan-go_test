package util

import "strings"

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func IsValidTaskStatus(status string) bool {
	switch status {
	case "pending", "in-progress", "completed":
		return true
	default:
		return false
	}
}
