package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

func serverSetup() {
	logInfo("Starting server")
	http.HandleFunc("/ws", handleConnection)
	go http.ListenAndServe(":8080", nil)
}


func handleConnection(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	c, _ := up.Upgrade(w, r, nil)
	defer c.Close()
	logInfo("Connected to client")
	defer logInfo("Disconnected from client")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter prompt (Ctrl+D to quit): ")
		if !scanner.Scan() {
			break
		}
		c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))

		for {
			_, data, _ := c.ReadMessage()
			if len(data) == 0 {
				break
			}
			saveImage(data)
			fmt.Println("Saved image")
		}
	}
}

func getWebsocketURL() string {
	logInfo("Starting ssh")

	cmd := exec.Command("ssh", "-nT", "-R", "80:localhost:8080", "nokey@localhost.run")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	prefix := "wss://"
	mark := ".lhr.life"
	suffix := "/ws"

	s := bufio.NewScanner(stdout); 
	s.Scan(); 
	s.Scan(); 
	line := s.Text()
	before, _, _ := strings.Cut(line, mark)
	wsURL := prefix + before + mark + suffix

	logInfo("Websocket URL: " + wsURL)
	return wsURL
}