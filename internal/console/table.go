package console

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Color codes for terminal output
const (
	ColorReset   = "\033[0m"
	ColorBold    = "\033[1m"
	ColorDim     = "\033[2m"
	ColorCyan    = "\033[96m"  // Bright cyan for borders
	ColorGreen   = "\033[92m"  // Bright green
	ColorYellow  = "\033[93m"  // Bright yellow
	ColorRed     = "\033[91m"  // Bright red
	ColorBlue    = "\033[94m"  // Bright blue for agent output
	ColorMagenta = "\033[95m"  // Bright magenta
	ColorWhite   = "\033[97m"  // Bright white
	ColorGray    = "\033[90m"  // Dark gray
)

// BoxTable provides professional table formatting for CLI output
type BoxTable struct {
	output      io.Writer
	title       string
	headers     []string
	rows        [][]string
	maxWidths   []int
	footer      string
	useBoxChars bool
	useColors   bool
}

// NewBoxTable creates a new table with beautiful box-drawing characters and colors
func NewBoxTable(output io.Writer) *BoxTable {
	if output == nil {
		output = os.Stdout
	}
	return &BoxTable{
		output:      output,
		useBoxChars: true,
		useColors:   true,
	}
}

// NewSimpleTable creates a table with simple ASCII characters (no colors)
func NewSimpleTable(output io.Writer) *BoxTable {
	if output == nil {
		output = os.Stdout
	}
	return &BoxTable{
		output:      output,
		useBoxChars: false,
		useColors:   false,
	}
}

// SetTitle sets the table title
func (t *BoxTable) SetTitle(title string) {
	t.title = title
}

// SetFooter sets the table footer
func (t *BoxTable) SetFooter(footer string) {
	t.footer = footer
}

// SetHeaders sets the table headers
func (t *BoxTable) SetHeaders(headers []string) {
	t.headers = headers
	t.calculateWidths()
}

// AddRow adds a new row to the table
func (t *BoxTable) AddRow(row []string) {
	t.rows = append(t.rows, row)
	t.calculateWidths()
}

// DisableColors disables color output
func (t *BoxTable) DisableColors() {
	t.useColors = false
}

// calculateWidths calculates the maximum width for each column
func (t *BoxTable) calculateWidths() {
	colCount := len(t.headers)
	if len(t.rows) > 0 {
		for _, row := range t.rows {
			if len(row) > colCount {
				colCount = len(row)
			}
		}
	}

	t.maxWidths = make([]int, colCount)

	// Header widths
	for i, header := range t.headers {
		if i < len(t.maxWidths) {
			t.maxWidths[i] = len(header)
		}
	}

	// Row widths
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(t.maxWidths) && len(cell) > t.maxWidths[i] {
				t.maxWidths[i] = len(cell)
			}
		}
	}
}

// Render renders the table with simple formatting (backward compatibility)
func (t *BoxTable) Render() {
	if t.title != "" {
		fmt.Fprintf(t.output, "\n%s\n", t.title)
		fmt.Fprintf(t.output, "%s\n", strings.Repeat("=", len(t.title)))
	}

	if len(t.headers) > 0 {
		t.printRow(t.headers)
		t.printSeparator()
	}

	for _, row := range t.rows {
		t.printRow(row)
	}

	if t.footer != "" {
		fmt.Fprintf(t.output, "\n%s\n", t.footer)
	}

	if len(t.headers) > 0 || len(t.rows) > 0 {
		fmt.Fprintln(t.output)
	}
}

// printRow prints a single row (simple format)
func (t *BoxTable) printRow(row []string) {
	for i, cell := range row {
		if i > 0 {
			fmt.Fprint(t.output, "  ")
		}
		width := 0
		if i < len(t.maxWidths) {
			width = t.maxWidths[i]
		}
		fmt.Fprintf(t.output, "%-*s", width, cell)
	}
	fmt.Fprintln(t.output)
}

// printSeparator prints a separator line (simple format)
func (t *BoxTable) printSeparator() {
	for i, width := range t.maxWidths {
		if i > 0 {
			fmt.Fprint(t.output, "  ")
		}
		fmt.Fprintf(t.output, "%s", strings.Repeat("-", width))
	}
	fmt.Fprintln(t.output)
}

// RenderBox renders a table with headers and multiple columns (PERFECT alignment)
func (t *BoxTable) RenderBox() {
	if !t.useBoxChars {
		t.Render()
		return
	}

	t.calculateWidths()

	// Calculate content width: sum of column widths + separators between columns
	contentWidth := 0
	for i, width := range t.maxWidths {
		contentWidth += width
		if i > 0 {
			contentWidth += 3 // " │ " separator
		}
	}

	// Box width = left padding (2) + content + right padding (2)
	boxWidth := contentWidth + 4

	// Ensure minimum width for title
	if t.title != "" {
		minWidth := len(t.title) + 6 // "┌─ " + title + " ─┐"
		if minWidth > boxWidth {
			boxWidth = minWidth
		}
	}

	// Ensure minimum width for footer
	if t.footer != "" {
		minWidth := len(t.footer) + 4
		if minWidth > boxWidth {
			boxWidth = minWidth
		}
	}

	// Top border with title
	if t.title != "" {
		remaining := boxWidth - len(t.title) - 4
		if remaining < 0 {
			remaining = 0
		}
		if t.useColors {
			fmt.Fprintf(t.output, "\n%s┌─ %s%s%s %s┐%s\n",
				ColorCyan, ColorBold+ColorWhite, t.title, ColorReset+ColorCyan,
				strings.Repeat("─", remaining), ColorReset)
		} else {
			fmt.Fprintf(t.output, "\n┌─ %s %s┐\n", t.title, strings.Repeat("─", remaining))
		}
	} else {
		if t.useColors {
			fmt.Fprintf(t.output, "\n%s┌%s┐%s\n", ColorCyan, strings.Repeat("─", boxWidth), ColorReset)
		} else {
			fmt.Fprintf(t.output, "\n┌%s┐\n", strings.Repeat("─", boxWidth))
		}
	}

	// Empty line after title
	t.printEmptyLine(boxWidth)

	// Headers
	if len(t.headers) > 0 {
		t.printBoxRow(t.headers, boxWidth, true)
		// Separator after headers
		if t.useColors {
			fmt.Fprintf(t.output, "%s├%s┤%s\n", ColorCyan, strings.Repeat("─", boxWidth), ColorReset)
		} else {
			fmt.Fprintf(t.output, "├%s┤\n", strings.Repeat("─", boxWidth))
		}
	}

	// Data rows
	for _, row := range t.rows {
		t.printBoxRow(row, boxWidth, false)
	}

	// Empty line before footer
	t.printEmptyLine(boxWidth)

	// Footer
	if t.footer != "" {
		padding := boxWidth - len(t.footer) - 4
		if padding < 0 {
			padding = 0
		}
		if t.useColors {
			fmt.Fprintf(t.output, "%s│%s  %s%s%s%s  %s│%s\n",
				ColorCyan, ColorReset, ColorDim, t.footer, ColorReset,
				strings.Repeat(" ", padding), ColorCyan, ColorReset)
		} else {
			fmt.Fprintf(t.output, "│  %s%s  │\n", t.footer, strings.Repeat(" ", padding))
		}
	}

	// Bottom border
	if t.useColors {
		fmt.Fprintf(t.output, "%s└%s┘%s\n\n", ColorCyan, strings.Repeat("─", boxWidth), ColorReset)
	} else {
		fmt.Fprintf(t.output, "└%s┘\n\n", strings.Repeat("─", boxWidth))
	}
}

// printBoxRow prints a single row in box format with PERFECT alignment
func (t *BoxTable) printBoxRow(row []string, boxWidth int, isHeader bool) {
	if t.useColors {
		fmt.Fprintf(t.output, "%s│%s  ", ColorCyan, ColorReset)
	} else {
		fmt.Fprint(t.output, "│  ")
	}

	// Print cells with separators
	for i, cell := range row {
		if i > 0 {
			if t.useColors {
				fmt.Fprintf(t.output, " %s│%s ", ColorCyan, ColorReset)
			} else {
				fmt.Fprint(t.output, " │ ")
			}
		}

		width := 0
		if i < len(t.maxWidths) {
			width = t.maxWidths[i]
		}

		if isHeader && t.useColors {
			fmt.Fprintf(t.output, "%s%-*s%s", ColorBold+ColorWhite, width, cell, ColorReset)
		} else {
			fmt.Fprintf(t.output, "%-*s", width, cell)
		}
	}

	// Calculate used width and add padding to reach box width
	usedWidth := 4 // left (2) + right (2) padding
	for i, width := range t.maxWidths {
		usedWidth += width
		if i > 0 {
			usedWidth += 3 // " │ "
		}
	}

	padding := boxWidth - usedWidth
	if padding < 0 {
		padding = 0
	}

	if t.useColors {
		fmt.Fprintf(t.output, "%s  %s│%s\n", strings.Repeat(" ", padding), ColorCyan, ColorReset)
	} else {
		fmt.Fprintf(t.output, "%s  │\n", strings.Repeat(" ", padding))
	}
}

// printEmptyLine prints an empty line in the box
func (t *BoxTable) printEmptyLine(boxWidth int) {
	if t.useColors {
		fmt.Fprintf(t.output, "%s│%s%s│%s\n", ColorCyan, ColorReset, strings.Repeat(" ", boxWidth), ColorCyan+ColorReset)
	} else {
		fmt.Fprintf(t.output, "│%s│\n", strings.Repeat(" ", boxWidth))
	}
}

// RenderKeyValue renders key-value pairs with PERFECT alignment
func (t *BoxTable) RenderKeyValue() {
	if !t.useBoxChars {
		t.Render()
		return
	}

	// Calculate max key width
	maxKeyWidth := 0
	for _, row := range t.rows {
		if len(row) > 0 && len(row[0]) > maxKeyWidth {
			maxKeyWidth = len(row[0])
		}
	}

	// Calculate max value width
	maxValueWidth := 0
	for _, row := range t.rows {
		if len(row) > 1 && len(row[1]) > maxValueWidth {
			maxValueWidth = len(row[1])
		}
	}

	// Box width = left pad (2) + key + " › " (3) + value + right pad (2)
	boxWidth := 2 + maxKeyWidth + 3 + maxValueWidth + 2

	// Ensure minimum width for title
	if t.title != "" {
		minWidth := len(t.title) + 7 // "┌── " + title + " ─┐"
		if minWidth > boxWidth {
			boxWidth = minWidth
		}
	}

	// Ensure minimum width for footer
	if t.footer != "" {
		minWidth := len(t.footer) + 4
		if minWidth > boxWidth {
			boxWidth = minWidth
		}
	}

	// Top border with title
	if t.title != "" {
		remaining := boxWidth - len(t.title) - 5
		if remaining < 0 {
			remaining = 0
		}
		if t.useColors {
			fmt.Fprintf(t.output, "\n%s┌── %s%s%s %s┐%s\n",
				ColorCyan, ColorBold+ColorWhite, t.title, ColorReset+ColorCyan,
				strings.Repeat("─", remaining), ColorReset)
		} else {
			fmt.Fprintf(t.output, "\n┌── %s %s┐\n", t.title, strings.Repeat("─", remaining))
		}
	} else {
		if t.useColors {
			fmt.Fprintf(t.output, "\n%s┌%s┐%s\n", ColorCyan, strings.Repeat("─", boxWidth), ColorReset)
		} else {
			fmt.Fprintf(t.output, "\n┌%s┐\n", strings.Repeat("─", boxWidth))
		}
	}

	// Empty line after title
	t.printEmptyLine(boxWidth)

	// Key-value rows
	for _, row := range t.rows {
		if len(row) >= 2 {
			key := row[0]
			value := row[1]

			// Calculate padding to reach box width
			usedWidth := 2 + maxKeyWidth + 3 + len(value) + 2
			padding := boxWidth - usedWidth
			if padding < 0 {
				padding = 0
			}

			if t.useColors {
				fmt.Fprintf(t.output, "%s│%s  %-*s %s›%s %s%s  %s│%s\n",
					ColorCyan, ColorReset, maxKeyWidth, key,
					ColorCyan, ColorReset, value,
					strings.Repeat(" ", padding), ColorCyan, ColorReset)
			} else {
				fmt.Fprintf(t.output, "│  %-*s › %s%s  │\n",
					maxKeyWidth, key, value, strings.Repeat(" ", padding))
			}
		}
	}

	// Empty line before footer
	t.printEmptyLine(boxWidth)

	// Footer
	if t.footer != "" {
		padding := boxWidth - len(t.footer) - 4
		if padding < 0 {
			padding = 0
		}
		if t.useColors {
			fmt.Fprintf(t.output, "%s│%s  %s%s%s%s  %s│%s\n",
				ColorCyan, ColorReset, ColorDim, t.footer, ColorReset,
				strings.Repeat(" ", padding), ColorCyan, ColorReset)
		} else {
			fmt.Fprintf(t.output, "│  %s%s  │\n", t.footer, strings.Repeat(" ", padding))
		}
	}

	// Bottom border
	if t.useColors {
		fmt.Fprintf(t.output, "%s└%s┘%s\n\n", ColorCyan, strings.Repeat("─", boxWidth), ColorReset)
	} else {
		fmt.Fprintf(t.output, "└%s┘\n\n", strings.Repeat("─", boxWidth))
	}
}
