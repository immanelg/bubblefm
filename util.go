package main

import (
	"log"
	"os"
	"strings"
)

func modulo(a, b int) int {
	return (a%b + b) % b
}

func filter[T any](arr *[]T, fn func(v T) bool) (result []T) {
	for _, v := range *arr {
		if fn(v) {
			result = append(result, v)
		}
	}
	return
}

func expandHome(s string) string {
	user, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot read user home dir")
		return s
	}
	return strings.Replace(s, "~", user, 1)
}

func withTilde(s string) string {
	user, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot read user home dir")
		return s
	}
	return strings.Replace(s, user, "~", 1)
}
