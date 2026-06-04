package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	rootDir string
	backendDir string
	genDir  string
	agentDir string

	nbMetaDataPath string
	nbCodePath string
)

func init() {
	_, rootFile, _, _ := runtime.Caller(0)

	backendDir = filepath.Dir(rootFile)
	rootDir = filepath.Dir(backendDir)
	genDir = rootDir + "/gen"
	agentDir = rootDir + "/agent"

	nbMetaDataPath = agentDir + "/kernel-metadata.json"
	nbCodePath = agentDir + "/notebook.ipynb"
}

func dirSetup() {
	os.MkdirAll(genDir, 0755)
}

func saveImage(data []byte) string {
	name := fmt.Sprintf("%d.png", time.Now().UnixNano())
	path := filepath.Join(genDir, name)
	os.WriteFile(path, data, 0644)
	return name
}