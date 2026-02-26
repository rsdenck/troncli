package console

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOutputJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantErr  bool
		validate func(string) bool
	}{
		{
			name: "simple map",
			data: map[string]string{
				"name":    "test",
				"version": "1.0",
			},
			wantErr: false,
			validate: func(output string) bool {
				var result map[string]string
				err := json.Unmarshal([]byte(output), &result)
				return err == nil && result["name"] == "test" && result["version"] == "1.0"
			},
		},
		{
			name: "struct with nested data",
			data: struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			}{
				Name:  "example",
				Count: 42,
			},
			wantErr: false,
			validate: func(output string) bool {
				return strings.Contains(output, `"name": "example"`) &&
					strings.Contains(output, `"count": 42`)
			},
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: false,
			validate: func(output string) bool {
				return strings.TrimSpace(output) == "null"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := OutputJSON(&buf, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("OutputJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()
			if !tt.validate(output) {
				t.Errorf("OutputJSON() output validation failed, got: %s", output)
			}

			// Verify proper indentation (2 spaces)
			if tt.data != nil && strings.Contains(output, "{") {
				lines := strings.Split(output, "\n")
				hasIndentation := false
				for _, line := range lines {
					if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") {
						hasIndentation = true
						break
					}
				}
				if !hasIndentation && len(lines) > 2 {
					t.Errorf("OutputJSON() expected 2-space indentation, got: %s", output)
				}
			}
		})
	}
}

func TestOutputYAML(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantErr  bool
		validate func(string) bool
	}{
		{
			name: "simple map",
			data: map[string]string{
				"name":    "test",
				"version": "1.0",
			},
			wantErr: false,
			validate: func(output string) bool {
				var result map[string]string
				err := yaml.Unmarshal([]byte(output), &result)
				return err == nil && result["name"] == "test" && result["version"] == "1.0"
			},
		},
		{
			name: "struct with nested data",
			data: struct {
				Name  string `yaml:"name"`
				Count int    `yaml:"count"`
			}{
				Name:  "example",
				Count: 42,
			},
			wantErr: false,
			validate: func(output string) bool {
				return strings.Contains(output, "name: example") &&
					strings.Contains(output, "count: 42")
			},
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: false,
			validate: func(output string) bool {
				trimmed := strings.TrimSpace(output)
				return trimmed == "null" || trimmed == ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := OutputYAML(&buf, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("OutputYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()
			if !tt.validate(output) {
				t.Errorf("OutputYAML() output validation failed, got: %s", output)
			}
		})
	}
}

func TestOutputTable(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		headers []string
		rows    [][]string
		footer  string
		want    []string // strings that should be present in output
	}{
		{
			name:    "simple table",
			title:   "Test Table",
			headers: []string{"Name", "Value"},
			rows: [][]string{
				{"foo", "bar"},
				{"baz", "qux"},
			},
			footer: "2 rows",
			want:   []string{"Test Table", "Name", "Value", "foo", "bar", "baz", "qux", "2 rows"},
		},
		{
			name:    "table without title",
			headers: []string{"Col1", "Col2"},
			rows: [][]string{
				{"a", "b"},
			},
			want: []string{"Col1", "Col2", "a", "b"},
		},
		{
			name:    "table without footer",
			title:   "No Footer",
			headers: []string{"X"},
			rows: [][]string{
				{"y"},
			},
			want: []string{"No Footer", "X", "y"},
		},
		{
			name:  "empty table",
			title: "Empty",
			want:  []string{"Empty"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			OutputTable(&buf, tt.title, tt.headers, tt.rows, tt.footer)

			output := buf.String()
			for _, expected := range tt.want {
				if !strings.Contains(output, expected) {
					t.Errorf("OutputTable() output missing expected string %q, got: %s", expected, output)
				}
			}
		})
	}
}

func TestOutputJSON_NilWriter(t *testing.T) {
	// This test verifies that OutputJSON handles nil writer by defaulting to os.Stdout
	// We can't easily capture os.Stdout in a test, so we just verify it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OutputJSON() with nil writer panicked: %v", r)
		}
	}()

	// This will write to stdout, which is acceptable behavior
	data := map[string]string{"test": "value"}
	_ = OutputJSON(nil, data)
}

func TestOutputYAML_NilWriter(t *testing.T) {
	// This test verifies that OutputYAML handles nil writer by defaulting to os.Stdout
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OutputYAML() with nil writer panicked: %v", r)
		}
	}()

	data := map[string]string{"test": "value"}
	_ = OutputYAML(nil, data)
}

func TestOutputTable_NilWriter(t *testing.T) {
	// This test verifies that OutputTable handles nil writer by defaulting to os.Stdout
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OutputTable() with nil writer panicked: %v", r)
		}
	}()

	OutputTable(nil, "Test", []string{"Col"}, [][]string{{"val"}}, "")
}
