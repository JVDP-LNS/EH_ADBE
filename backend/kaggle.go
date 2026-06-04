package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func pushKernel(wsURL string) {
	username := strings.TrimSpace(os.Getenv("KAGGLE_USERNAME"))
	key := strings.TrimSpace(os.Getenv("KAGGLE_KEY"))

	metaPath := nbMetaDataPath
	metaBytes, _ := os.ReadFile(metaPath)

	var meta kernelMeta
	json.Unmarshal(metaBytes, &meta)
	meta.ID = fmt.Sprintf("%s/%s", username, meta.ID)

	injectionVars := map[string]string{
		"BACKEND_WS_URL": wsURL,
	}
	codePath := nbCodePath
	codeBytes, _ := os.ReadFile(codePath)
	text := injectNotebook(codeBytes, injectionVars)

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
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "kaggle-api/v1.7.0")

	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Kaggle push failed:", string(body))
		return
	}
	defer resp.Body.Close()
	
	var out struct {
		URL   string `json:"url"`
		Error string `json:"error"`
		Ref   string `json:"ref"`
	}
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &out)

	fmt.Println("Kaggle pushed successfully:", out.URL)
}


func formatInjectionCell(vars map[string]string) string {
	var b strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&b, "%s = %q\n", k, v)
	}
	return b.String()
}

func injectNotebook(raw []byte, vars map[string]string) string {
	var nb map[string]any
	json.Unmarshal(raw, &nb)

	body := formatInjectionCell(vars)
	injectionCell := nb["cells"].([]any)[1].(map[string]any)
	injectionCell["source"] = append(injectionCell["source"].([]any), body)

	b, _ := json.Marshal(nb)
	return string(b)
}
