package ui

import (
	"fmt"
	"strings"
)

// Table represents a simple text table
type Table struct {
	headers []string
	rows    [][]string
	widths  []int
}

// NewTable creates a new table with headers
func NewTable(headers ...string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
		widths:  widths,
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(values ...string) {
	// Pad values if needed
	for len(values) < len(t.headers) {
		values = append(values, "")
	}

	// Update widths
	for i, v := range values {
		if i < len(t.widths) && len(v) > t.widths[i] {
			t.widths[i] = len(v)
		}
	}

	t.rows = append(t.rows, values)
}

// Render renders the table to stdout
func (t *Table) Render() {
	// Print headers
	headerLine := t.formatRow(t.headers, true)
	fmt.Println(headerLine)

	// Print separator
	sep := make([]string, len(t.headers))
	for i, w := range t.widths {
		sep[i] = strings.Repeat(SymbolHorizontal, w+2)
	}
	fmt.Println(Dim.Sprint(strings.Join(sep, "")))

	// Print rows
	for _, row := range t.rows {
		fmt.Println(t.formatRow(row, false))
	}
}

// formatRow formats a single row
func (t *Table) formatRow(values []string, isHeader bool) string {
	parts := make([]string, len(values))
	for i, v := range values {
		width := t.widths[i]
		padded := fmt.Sprintf("%-*s", width, v)
		if isHeader {
			parts[i] = Bold.Sprint(padded)
		} else {
			parts[i] = padded
		}
	}
	return "  " + strings.Join(parts, "  ")
}

// RenderCompact renders a more compact list format
func (t *Table) RenderCompact() {
	for _, row := range t.rows {
		if len(row) > 0 {
			fmt.Printf("  %s", row[0])
			if len(row) > 1 && row[1] != "" {
				fmt.Printf("  %s", Dim.Sprint(row[1]))
			}
			fmt.Println()
		}
	}
}
