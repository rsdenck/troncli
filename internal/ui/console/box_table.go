package console

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// BoxTable renders a table with borders
type BoxTable struct {
	Title   string
	Headers []string
	Rows    [][]string
	Footer  string
	Writer  io.Writer
}

// NewBoxTable creates a new BoxTable
func NewBoxTable(w io.Writer) *BoxTable {
	return &BoxTable{Writer: w}
}

// SetTitle sets the table title
func (t *BoxTable) SetTitle(title string) {
	t.Title = title
}

// SetHeaders sets the table headers
func (t *BoxTable) SetHeaders(headers []string) {
	t.Headers = headers
}

// AddRow adds a row of data
func (t *BoxTable) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
}

// SetFooter sets the table footer
func (t *BoxTable) SetFooter(footer string) {
	t.Footer = footer
}

// Render prints the table to the writer
func (t *BoxTable) Render() {
	if len(t.Headers) == 0 && len(t.Rows) == 0 {
		return
	}

	// 1. Calculate column widths
	numCols := len(t.Headers)
	if numCols == 0 && len(t.Rows) > 0 {
		numCols = len(t.Rows[0])
	}
	
	colWidths := make([]int, numCols)
	
	// Headers width
	for i, h := range t.Headers {
		if i < numCols {
			w := utf8.RuneCountInString(h)
			if w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	// Rows width
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < numCols {
				w := utf8.RuneCountInString(cell)
				if w > colWidths[i] {
					colWidths[i] = w
				}
			}
		}
	}

	// Add padding (1 space each side)
	for i := range colWidths {
		colWidths[i] += 2 
	}

	// Calculate total inner width (sum of col widths + internal borders)
	totalInnerWidth := 0
	for _, w := range colWidths {
		totalInnerWidth += w
	}
	if numCols > 1 {
		totalInnerWidth += numCols - 1
	}

	// Ensure title fits
	if t.Title != "" {
		titleLen := utf8.RuneCountInString(t.Title)
		if titleLen+2 > totalInnerWidth { 
			diff := (titleLen + 2) - totalInnerWidth
			if numCols > 0 {
				colWidths[numCols-1] += diff
				totalInnerWidth += diff
			}
		}
	}
    
    // Ensure footer fits
    if t.Footer != "" {
		footerLen := utf8.RuneCountInString(t.Footer)
		if footerLen+2 > totalInnerWidth {
			diff := (footerLen + 2) - totalInnerWidth
			if numCols > 0 {
				colWidths[numCols-1] += diff
				totalInnerWidth += diff
			}
		}
	}

	// Helper to print separator line
	printSep := func(left, mid, right, fill string) {
		fmt.Fprint(t.Writer, left)
		for i, w := range colWidths {
			fmt.Fprint(t.Writer, strings.Repeat(fill, w))
			if i < numCols-1 {
				fmt.Fprint(t.Writer, mid)
			}
		}
		fmt.Fprintln(t.Writer, right)
	}
    
    // Helper to print full line (for Title/Footer box)
    printFullLine := func(left, right, fill string) {
        fmt.Fprint(t.Writer, left)
        fmt.Fprint(t.Writer, strings.Repeat(fill, totalInnerWidth))
        fmt.Fprintln(t.Writer, right)
    }

    // Helper to print row
    printRow := func(row []string) {
        fmt.Fprint(t.Writer, "│")
        for i, w := range colWidths {
            cell := ""
            if i < len(row) {
                cell = row[i]
            }
            
            fmt.Fprint(t.Writer, " ")
            fmt.Fprint(t.Writer, cell)
            padding := w - 1 - utf8.RuneCountInString(cell) // -1 for left space
            if padding > 0 {
                fmt.Fprint(t.Writer, strings.Repeat(" ", padding))
            }
            
            if i < numCols-1 {
                fmt.Fprint(t.Writer, "│")
            }
        }
        fmt.Fprintln(t.Writer, "│")
    }

	// RENDER START
    
    // 1. Top Border
    if t.Title != "" {
        printFullLine("┌", "┐", "─")
        // Title Row
        fmt.Fprint(t.Writer, "│ ")
        fmt.Fprint(t.Writer, t.Title)
        pad := totalInnerWidth - 2 - utf8.RuneCountInString(t.Title)
        if pad > 0 {
            fmt.Fprint(t.Writer, strings.Repeat(" ", pad))
        }
        fmt.Fprintln(t.Writer, " │")
        
        // Separator below title
        printSep("├", "┬", "┤", "─")
    } else {
        printSep("┌", "┬", "┐", "─")
    }

	// 2. Headers
	if len(t.Headers) > 0 {
		printRow(t.Headers)
		printSep("├", "┼", "┤", "─")
	}

	// 3. Rows
	for _, row := range t.Rows {
		printRow(row)
	}

	// 4. Footer
	if t.Footer != "" {
        printSep("├", "┴", "┤", "─")
        // Footer Row
        fmt.Fprint(t.Writer, "│ ")
        fmt.Fprint(t.Writer, t.Footer)
        pad := totalInnerWidth - 2 - utf8.RuneCountInString(t.Footer)
        if pad > 0 {
            fmt.Fprint(t.Writer, strings.Repeat(" ", pad))
        }
        fmt.Fprintln(t.Writer, " │")
        
        // Bottom
        printFullLine("└", "┘", "─")
	} else {
		printSep("└", "┴", "┘", "─")
	}
}
