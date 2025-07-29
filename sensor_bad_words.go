package main

import (
	"slices"
	"strings"
)

func censorBadWords(s string) string {
	isBadWord := func(w string) bool {
		badWords := []string{"kerfuffle", "sharbert", "fornax"}
		return slices.Contains(badWords, strings.ToLower(w))
	}
	words := strings.Split(s, " ")
	for idx, val := range words {
		if isBadWord(val) {
			words[idx] = "****"
		}
	}
	return strings.Join(words, " ")
}
