package table

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// TableOptions provides configuration for table rendering
type TableOptions struct {
	// AutoWrapText controls text wrapping in cells
	AutoWrapText bool
	// AutoFormatHeaders controls header formatting
	AutoFormatHeaders bool
	// HeaderAlignment sets the alignment for header cells
	HeaderAlignment int
	// CellAlignment sets the alignment for non-header cells
	CellAlignment int
	// UseBorders enables table borders
	UseBorders bool
	// Padding sets the padding between columns
	Padding string
	// FooterData adds a footer row to the table
	FooterData []string
	// Caption sets a caption for the table
	Caption string
}

// DefaultOptions returns the default table formatting options
func DefaultOptions() TableOptions {
	return TableOptions{
		AutoWrapText:      false,
		AutoFormatHeaders: true,
		HeaderAlignment:   tablewriter.ALIGN_LEFT,
		CellAlignment:     tablewriter.ALIGN_LEFT,
		UseBorders:        false,
		Padding:           "\t",
	}
}

// Print outputs a table to stdout with default styling
func Print(header []string, body [][]string) {
	PrintWithOptions(os.Stdout, header, body, DefaultOptions())
}

// PrintWithWriter outputs a table to the provided writer with default styling
func PrintWithWriter(writer io.Writer, header []string, body [][]string) {
	PrintWithOptions(writer, header, body, DefaultOptions())
}

// PrintWithOptions outputs a table with the provided configuration options
func PrintWithOptions(writer io.Writer, header []string, body [][]string, options TableOptions) {
	table := tablewriter.NewWriter(writer)

	if len(header) > 0 {
		table.SetHeader(header)
	}

	// Configure the table based on options
	table.SetAutoWrapText(options.AutoWrapText)
	table.SetAutoFormatHeaders(options.AutoFormatHeaders)
	table.SetHeaderAlignment(options.HeaderAlignment)
	table.SetAlignment(options.CellAlignment)

	if !options.UseBorders {
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
	}

	table.SetTablePadding(options.Padding)
	table.SetNoWhiteSpace(true)

	if options.Caption != "" {
		table.SetCaption(true, options.Caption)
	}

	if len(options.FooterData) > 0 {
		table.SetFooter(options.FooterData)
	}

	if len(body) > 0 {
		table.AppendBulk(body)
	}

	table.Render()
}

// RenderToString renders a table and returns it as a string
func RenderToString(header []string, body [][]string) string {
	var builder strings.Builder
	PrintWithWriter(&builder, header, body)
	return builder.String()
}

// Colors for styling table content
const (
	// Color constants for ANSI colors
	ColorBlack   = 0
	ColorRed     = 1
	ColorGreen   = 2
	ColorYellow  = 3
	ColorBlue    = 4
	ColorMagenta = 5
	ColorCyan    = 6
	ColorWhite   = 7
)

// CreateColoredCell returns a cell with the specified color
func CreateColoredCell(content string, color int) string {
	return fmt.Sprintf("\033[3%dm%s\033[0m", color, content)
}

// Red returns text colored in red
func Red(content string) string {
	return CreateColoredCell(content, ColorRed)
}

// Green returns text colored in green
func Green(content string) string {
	return CreateColoredCell(content, ColorGreen)
}

// Yellow returns text colored in yellow
func Yellow(content string) string {
	return CreateColoredCell(content, ColorYellow)
}

// Blue returns text colored in blue
func Blue(content string) string {
	return CreateColoredCell(content, ColorBlue)
}
