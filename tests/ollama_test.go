package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const ollamaHost = "http://192.168.130.25:11434"

func TestOllamaGenerate(t *testing.T) {
	payload := map[string]interface{}{
		"model":  "gemma:2b",
		"prompt": "test",
		"stream": false,
	}
	data, _ := json.Marshal(payload)
	
	client := &http.Client{Timeout: 360 * time.Second}
	resp, err := client.Post(ollamaHost+"/api/generate", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to connect to Ollama: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if _, ok := result["response"].(string); !ok {
		t.Fatal("No response field in Ollama result")
	}
	fmt.Println("✓ Ollama Generate OK")
}

func TestOllamaTags(t *testing.T) {
	resp, err := http.Get(ollamaHost + "/api/tags")
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse tags: %v", err)
	}
	
	if _, ok := result["models"]; !ok {
		t.Fatal("No models field in tags response")
	}
	fmt.Println("✓ Ollama Tags OK")
}

func TestOllamaVersion(t *testing.T) {
	resp, err := http.Get(ollamaHost + "/api/version")
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("version")) {
		t.Fatal("Invalid version response")
	}
	fmt.Println("✓ Ollama Version OK")
}
