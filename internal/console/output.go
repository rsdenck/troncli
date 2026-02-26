package console

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// OutputJSON writes data as JSON to the given writer with 2-space indentation
func OutputJSON(w io.Writer, data interface{}) error {
	if w == nil {
		w = os.Stdout
	}
	
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// OutputYAML writes data as YAML to the given writer
func OutputYAML(w io.Writer, data interface{}) error {
	if w == nil {
		w = os.Stdout
	}
	
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	defer encoder.Close()
	
	return encoder.Encode(data)
}

// OutputTable is a wrapper for BoxTable that provides a consistent interface
// for table output alongside JSON and YAML output functions
func OutputTable(w io.Writer, title string, headers []string, rows [][]string, footer string) {
	if w == nil {
		w = os.Stdout
	}
	
	table := NewBoxTable(w)
	
	if title != "" {
		table.SetTitle(title)
	}
	
	if len(headers) > 0 {
		table.SetHeaders(headers)
	}
	
	for _, row := range rows {
		table.AddRow(row)
	}
	
	if footer != "" {
		table.SetFooter(footer)
	}
	
	table.Render()
}
