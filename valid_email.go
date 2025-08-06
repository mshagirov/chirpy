package main

import "strings"

func validEmail(email string) bool {
	if !strings.ContainsRune(email, '@') {
		return false
	}
	if strings.ContainsRune(email, ' ') {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) < 1 {
		return false
	}
	domain := strings.Split(parts[1], ".")
	if len(domain) != 2 {
		return false
	}
	return true
}
