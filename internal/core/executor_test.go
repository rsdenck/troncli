package core

import (
	"fmt"
	"testing"
)

func TestRealExecutor(t *testing.T) {
	executor := &RealExecutor{}
	
	// Test simple command
	output, err := executor.Run("echo", "hello")
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}
	
	if output != "hello" {
		t.Errorf("Expected 'hello', got '%s'", output)
	}
}

func TestMockExecutor(t *testing.T) {
	mock := &MockExecutor{
		Output: "mock output",
		Err:    nil,
	}
	
	output, err := mock.Run("any", "command")
	if err != nil {
		t.Errorf("Mock Run failed: %v", err)
	}
	
	if output != "mock output" {
		t.Errorf("Expected 'mock output', got '%s'", output)
	}
}

func TestMockExecutorError(t *testing.T) {
	mock := &MockExecutor{
		Output: "",
		Err:    fmt.Errorf("mock error"),
	}
	
	_, err := mock.Run("fail", "command")
	if err == nil {
		t.Error("Expected error from mock")
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello; rm -rf /", "hello rm -rf /"},
		{"test && cat /etc/passwd", "test  cat /etc/passwd"},
		{"normal text", "normal text"},
		{"$(whoami)", "whoami)"},  // $( removed for security
		{"`id`", "id"},           // backticks removed
		{"ls -la /etc/passwd", "ls -la /etc/passwd"}, // safe command
	}
	
	for _, tt := range tests {
		got := SanitizeInput(tt.input)
		if got != tt.want {
			t.Errorf("SanitizeInput(%s) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/etc/passwd", true},
		{"/usr/bin", true},
		{"relative/path", false},
		{"/etc/../passwd", false},
		{"/safe/path", true},
	}
	
	for _, tt := range tests {
		got := ValidatePath(tt.path)
		if got != tt.want {
			t.Errorf("ValidatePath(%s) = %v, want %v", tt.path, got, tt.want)
		}
	}
}
