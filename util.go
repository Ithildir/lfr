package main

import "os"

func fileExists(f string) bool {
	_, err := os.Stat(f)

	if os.IsNotExist(err) {
		return false
	}

	return true
}
