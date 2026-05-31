package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

var queries = []string{
	"What is 2+2? Answer in one short sentence.",
	"Name one primary color. One word only.",
}

func main() {
	fmt.Println("starting server")
	up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, _ := up.Upgrade(w, r, nil)
		fmt.Println("client connected")
		defer fmt.Println("client disconnected")

		_, b, _ := c.ReadMessage()
		fmt.Println(string(b))

		for _, q := range queries {
			c.WriteMessage(websocket.TextMessage, []byte(q))
			fmt.Printf("\n--- query: %s ---\n", q)
			for {
				_, b, err := c.ReadMessage()
				if err != nil {
					return
				}
				msg := string(b)
				if msg == "END" {
					fmt.Println()
					break
				}
				fmt.Print(msg)
			}
		}
		c.WriteMessage(websocket.TextMessage, []byte("DONE"))
	})
	go http.ListenAndServe(":8080", nil)

	fmt.Println("starting ssh")
	cmd := exec.Command("ssh", "-nT", "-R", "80:localhost:8080", "nokey@localhost.run")
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	re := regexp.MustCompile(`tunneled with tls termination,\s*(https://\S+)`)
	found := make(chan string, 1)
	scan := func(r io.Reader) {
		for s := bufio.NewScanner(r); s.Scan(); {
			if m := re.FindStringSubmatch(s.Text()); m != nil {
				found <- m[1]
				return
			}
		}
	}
	go scan(stdout)
	go scan(stderr)
	https := <-found
	wsURL := strings.Replace(https, "https://", "wss://", 1) + "/ws"
	fmt.Println("websocket URL:", wsURL)

	pushKernel("agent", wsURL)

	select {}
}
