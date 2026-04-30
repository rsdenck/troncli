package console

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

type OutputFormat int

const (
	FormatText OutputFormat = iota
	FormatJSON
	FormatYAML
)

var currentFormat = FormatText

func SetFormat(json bool, yaml bool) {
	if json {
		currentFormat = FormatJSON
	} else if yaml {
		currentFormat = FormatYAML
	} else {
		currentFormat = FormatText
	}
}

func Printf(format string, args ...interface{}) {
	if currentFormat == FormatText {
		fmt.Printf(format, args...)
	}
}

func Println(args ...interface{}) {
	if currentFormat == FormatText {
		fmt.Println(args...)
	}
}

func PrintJSON(data interface{}) {
	if currentFormat == FormatJSON {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(jsonData))
	}
}

func PrintYAML(data interface{}) {
	if currentFormat == FormatYAML {
		Println("---")
		Println(data)
	}
}

type BoxTable struct {
	writer  io.Writer
	headers []string
	rows    [][]string
	title   string
	footer  string
}

func NewBoxTable(w io.Writer) *BoxTable {
	return &BoxTable{writer: w}
}

func (t *BoxTable) SetTitle(title string) {
	t.title = title
}

func (t *BoxTable) SetHeaders(headers []string) {
	t.headers = headers
}

func (t *BoxTable) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

func (t *BoxTable) SetFooter(footer string) {
	t.footer = footer
}

func (t *BoxTable) Render() {
	if t.title != "" {
		fmt.Fprintln(t.writer)
		fmt.Fprintln(t.writer, t.title)
		fmt.Fprintln(t.writer, strings.Repeat("=", len(t.title)))
	}

	if len(t.headers) > 0 {
		w := tabwriter.NewWriter(t.writer, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, strings.Join(t.headers, "\t"))
		fmt.Fprintln(w, strings.Repeat("-\t", len(t.headers)))
		for _, row := range t.rows {
			fmt.Fprintln(w, strings.Join(row, "\t"))
		}
		w.Flush()
	}

	if t.footer != "" {
		fmt.Fprintln(t.writer, t.footer)
	}

	fmt.Fprintln(t.writer)
}

func PrintSuccess(msg string) {
	fmt.Printf("✓ %s\n", msg)
}

func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
}

func PrintInfo(msg string) {
	fmt.Printf("ℹ %s\n", msg)
}
