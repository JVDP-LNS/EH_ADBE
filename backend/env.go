package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)


func init() {
	wd, _ := os.Getwd()
	envPath := filepath.Join(wd, "../.env")
	loadDotEnv(envPath)
}


func loadDotEnv(path string) {
	f, _ := os.Open(path)
	defer f.Close()
	s := bufio.NewScanner(f)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line != "" && line[0] != '#' {
			k, v, _ := strings.Cut(line, "=")
			os.Setenv(k, v)
		}
	}
}
