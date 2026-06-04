package main

import "fmt"

func logInfo(message string) {
	fmt.Println("[INFO] \t\t" + message)
}

func logWarning(message string) {
	fmt.Println("[WARNING]\t" + message)
}

func logError(message string) {
	fmt.Println("[ERROR]\t\t" + message)
}