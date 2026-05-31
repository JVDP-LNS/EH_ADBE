package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func genDir() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "../gen")
}

func saveImage(gpuID int, data []byte) (string, error) {
	dir := genDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d_gpu%d.png", time.Now().UnixNano(), gpuID)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func recvImages(s *wsSession) ([]string, error) {
	var paths []string
	expectBin := false
	gpuID := 0

	for {
		select {
		case msg := <-s.textCh:
			if msg == "END" {
				return paths, nil
			}
			if !strings.HasPrefix(msg, "IMG:") {
				continue
			}
			id, err := strconv.Atoi(strings.TrimPrefix(msg, "IMG:"))
			if err != nil {
				return paths, fmt.Errorf("invalid IMG header: %s", msg)
			}
			gpuID = id
			expectBin = true

		case data := <-s.binCh:
			if !expectBin {
				continue
			}
			path, err := saveImage(gpuID, data)
			if err != nil {
				return paths, err
			}
			paths = append(paths, path)
			fmt.Println("saved:", path)
			expectBin = false

		case err := <-s.errCh:
			return paths, err
		case <-s.done:
			return paths, fmt.Errorf("connection closed")
		}
	}
}
