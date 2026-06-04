package main

import "fmt"

func main() {
	dirSetup()
	serverSetup()
	wsURL := getWebsocketURL()
	fmt.Println("Websocket URL:", wsURL)
	// pushKernel(wsURL)
	select {}
}
