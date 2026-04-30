package output

import (
	"testing"
	"encoding/json"
)

func TestNewSuccess(t *testing.T) {
	result := NewSuccess("test data")
	
	if result == nil {
		t.Error("NewSuccess should return non-nil Output")
	}
	
	if result.Data != "test data" {
		t.Errorf("Expected Data to be 'test data', got %v", result.Data)
	}
	
	if result.Status != "success" {
		t.Errorf("Expected Status to be 'success', got %s", result.Status)
	}
}

func TestNewError(t *testing.T) {
	err := NewError("test error", "TEST_CODE")
	
	if err == nil {
		t.Error("NewError should return non-nil Output")
	}
	
	// Note: NewError sets Error field, not Message
	if err.Error != "test error" {
		t.Errorf("Expected Error to be 'test error', got %s", err.Error)
	}
	
	if err.Code != "TEST_CODE" {
		t.Errorf("Expected Code to be 'TEST_CODE', got %s", err.Code)
	}
}

func TestWithMessage(t *testing.T) {
	result := NewSuccess("data").WithMessage("test message")
	
	if result.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %s", result.Message)
	}
}

func TestNewList(t *testing.T) {
	items := []map[string]interface{}{
		{"name": "item1"},
		{"name": "item2"},
	}
	
	list := NewList(items, len(items))
	
	if list == nil {
		t.Error("NewList should return non-nil Output")
	}
	
	if list.Items == nil {
		t.Error("Items should not be nil")
	}
	
	if list.Total != 2 {
		t.Errorf("Expected Total to be 2, got %d", list.Total)
	}
}

func TestPrintJSON(t *testing.T) {
	// Set JSON format
	SetFormat(true, false)
	
	result := NewSuccess(map[string]interface{}{"key": "value"})
	
	// This should not panic
	result.Print()
	
	// Reset format
	SetFormat(false, false)
	
	// Verify JSON output
	data, err := json.Marshal(result.Data)
	if err != nil {
		t.Errorf("Failed to marshal result: %v", err)
	}
	
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal result: %v", err)
	}
}

func TestNewInfo(t *testing.T) {
	// Note: NewInfo sets Data field, not Message
	info := NewInfo("info message")
	
	if info == nil {
		t.Error("NewInfo should return non-nil Output")
	}
	
	if info.Data != "info message" {
		t.Errorf("Expected Data to be 'info message', got %v", info.Data)
	}
	
	if info.Status != "info" {
		t.Errorf("Expected Status to be 'info', got %s", info.Status)
	}
}

func TestOutputStatus(t *testing.T) {
	success := NewSuccess("data")
	if success.Status != "success" {
		t.Errorf("Expected Status to be 'success', got %s", success.Status)
	}
	
	err := NewError("error", "CODE")
	if err.Status != "error" {
		t.Errorf("Expected Status to be 'error', got %s", err.Status)
	}
	
	info := NewInfo("info")
	if info.Status != "info" {
		t.Errorf("Expected Status to be 'info', got %s", info.Status)
	}
}
