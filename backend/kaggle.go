package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// KagglePushReq strictly mirrors the official Kaggle API JSON schema
type KagglePushReq struct {
	Slug               string   `json:"slug"`
	NewTitle           string   `json:"newTitle"`
	Text               string   `json:"text"`
	Language           string   `json:"language"`
	KernelType         string   `json:"kernelType"`
	IsPrivate          bool     `json:"isPrivate"`
	EnableGpu          bool     `json:"enableGpu"`
	EnableTpu          bool     `json:"enableTpu"`
	EnableInternet     bool     `json:"enableInternet"`
	MachineShape       string   `json:"machineShape,omitempty"`
	DatasetDataSources []string `json:"datasetDataSources"`
}

type kernelMeta struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	CodeFile       string   `json:"code_file"`
	Language       string   `json:"language"`
	KernelType     string   `json:"kernel_type"`
	IsPrivate      bool     `json:"is_private"`
	EnableGPU      bool     `json:"enable_gpu"`
	MachineShape   string   `json:"machine_shape"`
	EnableInternet bool     `json:"enable_internet"`
	DatasetSources []string `json:"dataset_sources"`
}

func pushKernel(agentDir string, wsURL string) {
	username := strings.TrimSpace(os.Getenv("KAGGLE_USERNAME"))
	key := strings.TrimSpace(os.Getenv("KAGGLE_KEY"))

	wd, _ := os.Getwd()
	metaPath := filepath.Join(wd, "../agent/kernel-metadata.json")
	metaBytes, _ := os.ReadFile(metaPath)

	var meta kernelMeta
	json.Unmarshal(metaBytes, &meta)
	meta.ID = fmt.Sprintf("%s/%s", username, meta.ID)

	codePath := filepath.Join(wd, "../agent", meta.CodeFile)
	codeBytes, _ := os.ReadFile(codePath)
	codeBytes = injectNotebook(codeBytes, map[string]string{
		"BACKEND_WS_URL": wsURL,
	})
	text := notebookForPush(codeBytes)
	if meta.MachineShape != "" {
		fmt.Println("machine shape:", meta.MachineShape)
	}

	// 4. Build the strict JSON request payload
	reqPayload := KagglePushReq{
		Slug:               meta.ID,
		NewTitle:           meta.Title,
		Text:               text,
		Language:           meta.Language,
		KernelType:         meta.KernelType,
		IsPrivate:          meta.IsPrivate,
		EnableGpu:          meta.EnableGPU,
		EnableTpu:          false,
		EnableInternet:     meta.EnableInternet,
		MachineShape:       meta.MachineShape,
		DatasetDataSources: meta.DatasetSources,
	}

	payloadBytes, _ := json.Marshal(reqPayload)

	kaggleAPI := "https://www.kaggle.com/api/v1/kernels/push"

	req, _ := http.NewRequest(http.MethodPost, kaggleAPI, bytes.NewReader(payloadBytes))
	// req.SetBasicAuth(username, key)
	// Send the KCAT_ token as a Bearer token instead of Basic Auth
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "kaggle-api/v1.7.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP Error %d: %s\n", resp.StatusCode, string(body))
		return
	}

	var out struct {
		URL   string `json:"url"`
		Error string `json:"error"`
		Ref   string `json:"ref"`
	}
	_ = json.Unmarshal(body, &out)

	if out.Error != "" {
		fmt.Println("Kaggle returned error:", out.Error)
	} else {
		fmt.Println("Kaggle pushed successfully:", out.URL)
	}
}

func notebookForPush(raw []byte) string {
	var nb map[string]any
	json.Unmarshal(raw, &nb)

	for _, c := range nb["cells"].([]any) {
		cell, _ := c.(map[string]any)
		if cell["cell_type"] == "code" {
			cell["outputs"] = []any{}
			cell["execution_count"] = nil
		}
	}

	meta, _ := nb["metadata"].(map[string]any)
	if meta == nil {
		meta = map[string]any{}
		nb["metadata"] = meta
	}
	if _, ok := meta["kernelspec"]; !ok {
		meta["kernelspec"] = map[string]any{
			"display_name": "Python 3",
			"language":     "python",
			"name":         "python3",
		}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(nb)
	return buf.String()
}

func formatInjectionCell(vars map[string]string) string {
	var b strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&b, "%s = %q\n", k, v)
	}
	return b.String()
}

// injectNotebook appends injected variables to the injection cell (index 2).
func injectNotebook(raw []byte, vars map[string]string) []byte {
	var nb map[string]any
	json.Unmarshal(raw, &nb)

	body := formatInjectionCell(vars)
	injectionCell := nb["cells"].([]any)[2].(map[string]any)
	injectionCell["source"] = append(injectionCell["source"].([]any), body)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(nb)
	return buf.Bytes()
}
