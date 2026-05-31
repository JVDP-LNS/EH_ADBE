package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	if err := os.MkdirAll(genDir(), 0755); err != nil {
		fmt.Println("failed to create gen dir:", err)
		return
	}

	fmt.Println("starting server")
	up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		fmt.Println("client connected")
		defer fmt.Println("client disconnected")

		session := newWSSession(c)
		defer session.close()

		msg, err := session.recvText()
		if err != nil {
			fmt.Println("ready error:", err)
			return
		}
		fmt.Println(msg)

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter prompt (or 'exit'): ")
			if !scanner.Scan() {
				break
			}
			q := strings.TrimSpace(scanner.Text())
			if q == "" {
				continue
			}
			if q == "exit" {
				session.sendText("DONE")
				break
			}

			if err := session.sendText(q); err != nil {
				fmt.Println("send error:", err)
				return
			}
			fmt.Printf("\n--- query: %s ---\n", q)

			paths, err := recvImages(session)
			if err != nil {
				fmt.Println("recv error:", err)
				return
			}
			fmt.Printf("received %d image(s)\n", len(paths))
		}
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
