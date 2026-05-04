package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const nvidiaAPI = "https://integrate.api.nvidia.com/v1/chat/completions"
const nvidiaKey = "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"

func TestNvidiaMinimax(t *testing.T) {
	payload := map[string]interface{}{
		"model": "minimaxai/minimax-m2.7",
		"messages": []map[string]string{
			{"role": "user", "content": "test"},
		},
		"max_tokens": 50,
	}
	data, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", nvidiaAPI, strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+nvidiaKey)
	
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to connect to NVIDIA: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if choices, ok := result["choices"].([]interface{}); !ok || len(choices) == 0 {
		t.Fatal("No choices in NVIDIA response")
	}
	fmt.Println("✓ NVIDIA MiniMax OK")
}

func TestNvidiaDeepseek(t *testing.T) {
	payload := map[string]interface{}{
		"model": "deepseek-ai/deepseek-v3.2",
		"messages": []map[string]string{
			{"role": "user", "content": "test"},
		},
		"max_tokens": 50,
		"extra_body": map[string]interface{}{
			"chat_template_kwargs": map[string]bool{
				"thinking": true,
			},
		},
	}
	data, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", nvidiaAPI, strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+nvidiaKey)
	
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to connect to NVIDIA: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	// Debug
	bodyStr := string(body)
	if len(bodyStr) > 200 {
		bodyStr = bodyStr[:200]
	}
	t.Logf("DeepSeek status: %d, body: %s", resp.StatusCode, bodyStr)
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if choices, ok := result["choices"].([]interface{}); !ok || len(choices) == 0 {
		t.Fatal("No choices in NVIDIA response")
	}
	fmt.Println("✓ NVIDIA DeepSeek OK")
}
