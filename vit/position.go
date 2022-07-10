package vit

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// position describes a specific position in a file
type Position struct {
	FilePath string
	Line     int // line inside the file starting at 1
	Column   int // column inside the line starting at 1 (this is pointing to the rune, not the byte)
}

// String returns a human readable description of the position
func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.FilePath, p.Line, p.Column)
}

// IsEqual returns true if with positions point to the same location in the same file
func (p Position) IsEqual(o Position) bool {
	return p.FilePath == o.FilePath && p.Line == o.Line && p.Column == o.Column
}

// positionRange describes a range of runes in a file
type PositionRange struct {
	FilePath string
	// start points to the first rune of the range
	StartLine   int // line inside the file starting at 1
	StartColumn int // column inside the line starting at 1 (this is pointing to the rune, not the byte)
	// end points to the last rune of the range
	EndLine   int // line inside the file starting at 1
	EndColumn int // column inside the line starting at 1 (this is pointing to the rune, not the byte)
}

// newRangeFromStartToEnd returns a range that starts and ends at the same position.
func NewRangeFromPosition(pos Position) PositionRange {
	return PositionRange{
		FilePath:    pos.FilePath,
		StartLine:   pos.Line,
		StartColumn: pos.Column,
		EndLine:     pos.Line,
		EndColumn:   pos.Column,
	}
}

// newRangeFromStartToEnd returns a range from the start position to the end position.
// The filePath is taken from the start position.
func NewRangeFromStartToEnd(start Position, end Position) PositionRange {
	return PositionRange{
		FilePath:    start.FilePath,
		StartLine:   start.Line,
		StartColumn: start.Column,
		EndLine:     end.Line,
		EndColumn:   end.Column,
	}
}

func CombineRanges(a, b PositionRange) PositionRange {
	if a.FilePath != b.FilePath {
		fmt.Printf("RangeUnion has been called with two ranges from different files")
	}
	out := PositionRange{FilePath: a.FilePath}

	if a.StartLine < b.StartLine {
		out.StartLine = a.StartLine
		out.StartColumn = a.StartColumn
	} else if a.StartLine > b.StartLine {
		out.StartLine = b.StartLine
		out.StartColumn = b.StartColumn
	} else {
		out.StartLine = a.StartLine
		if a.StartColumn < b.StartColumn {
			out.StartColumn = a.StartColumn
		} else {
			out.StartColumn = b.StartColumn
		}
	}

	if a.EndLine < b.EndLine {
		out.EndLine = b.EndLine
		out.EndColumn = b.EndColumn
	} else if a.EndLine > b.EndLine {
		out.EndLine = a.EndLine
		out.EndColumn = a.EndColumn
	} else {
		out.EndLine = a.EndLine
		if a.EndColumn < b.EndColumn {
			out.EndColumn = b.EndColumn
		} else {
			out.EndColumn = a.EndColumn
		}
	}

	return out
}

// String returns a human readable description of the range between two position
func (p PositionRange) String() string {
	return fmt.Sprintf("%s:%d:%d", p.FilePath, p.StartLine, p.StartColumn)
}

func (p PositionRange) Report() string {
	f, err := os.ReadFile(p.FilePath)
	if err != nil {
		fmt.Printf("(unable to generate detailed position string: %v)\r\n", err)
		return p.String()
	}
	lines := strings.Split(string(f), "\n")
	var out strings.Builder

	if p.StartLine == p.EndLine {
		out.WriteString(fmt.Sprintf("%s:%d:%d\r\n", p.FilePath, p.StartLine, p.StartColumn))

		// The line prefix. We need to know it's length to know where to start printing the markers in the line below.
		lineNumberPrefix := fmt.Sprintf(" %s | ", formatLineNumber(p.StartLine, p.EndLine))
		lineContent := lines[p.StartLine-1]
		// Because we will remove any leading spaces and tabs we need to know how many there are to adjust the markers.
		trimmedCharacters := len(lineContent) - len(strings.TrimLeft(lineContent, " \t"))
		// Here we actually remove all leading and trailing spaces and tabs (and potential \r while we are at it). We also replace all remaining tabs that might be in the line with spaces
		// to make sure the markers line up properly because we can't tell how wide tabs in this line would be and only a single marker would be printed.
		trimmedContent := strings.ReplaceAll(strings.Trim(lineContent, " \t\r"), "\t", " ")

		out.WriteString(lineNumberPrefix)
		out.WriteString(trimmedContent)
		out.WriteString("\r\n")
		out.WriteString(strings.Repeat(" ", len(lineNumberPrefix)-trimmedCharacters+p.StartColumn-1)) // leading spaces to get the right offset
		out.WriteString(strings.Repeat("^", p.EndColumn-p.StartColumn+1))                             // now print the markers
	} else {
		out.WriteString("report for multiline range is not implemnted yet\r\n")
	}

	return out.String()
}

// Start returns the position of the first rune
func (p PositionRange) Start() Position {
	return Position{
		FilePath: p.FilePath,
		Line:     p.StartLine,
		Column:   p.StartColumn,
	}
}

// End returns the position of the last rune
func (p PositionRange) End() Position {
	return Position{
		FilePath: p.FilePath,
		Line:     p.EndLine,
		Column:   p.EndColumn,
	}
}

func (p *PositionRange) SetEnd(pos Position) {
	p.EndLine = pos.Line
	p.EndColumn = pos.Column
}

// StartColumnShifted returns a new position range whose start has been shifted by the given amount.
// If the end is on the same line it's column is shifted as well.
func (p PositionRange) StartColumnShifted(amount int) PositionRange {
	pos := PositionRange{
		FilePath:    p.FilePath,
		StartLine:   p.StartLine,
		StartColumn: p.StartColumn - 10,
		EndLine:     p.EndLine,
		EndColumn:   p.EndColumn + 10,
	}
	if pos.EndLine == pos.StartLine {
		pos.EndColumn -= 10
	}
	return pos
}

// formatLineNumber returns a stringified version of the line number with enough padded spaces to accommodate the highest line number
func formatLineNumber(line int, highest int) string {
	return fmt.Sprintf(
		fmt.Sprintf("%%%dd", digits(highest)), // generate a format string with the correct padding number
		line,
	)
}

// digits returns the number of digits in the given number
func digits(n int) int {
	if n == 0 {
		return 1
	}
	return int(math.Log10(float64(n))) + 1
}
