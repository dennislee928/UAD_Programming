package common

import "fmt"

// Position represents a location in source code
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column number
	Offset int // 0-based byte offset from start of file
}

// Span represents a range in source code
type Span struct {
	Start Position
	End   Position
	File  string // Source file path
}

// NewPosition creates a new Position
func NewPosition(line, column, offset int) Position {
	return Position{
		Line:   line,
		Column: column,
		Offset: offset,
	}
}

// NewSpan creates a new Span
func NewSpan(start, end Position, file string) Span {
	return Span{
		Start: start,
		End:   end,
		File:  file,
	}
}

// String returns a human-readable string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// String returns a human-readable string representation of the span
func (s Span) String() string {
	if s.File == "" {
		return fmt.Sprintf("%s-%s", s.Start.String(), s.End.String())
	}
	return fmt.Sprintf("%s:%s-%s", s.File, s.Start.String(), s.End.String())
}

// IsValid checks if the position is valid (non-zero)
func (p Position) IsValid() bool {
	return p.Line > 0 && p.Column > 0
}

// IsValid checks if the span is valid
func (s Span) IsValid() bool {
	return s.Start.IsValid() && s.End.IsValid()
}

// Contains checks if this span contains a given position
func (s Span) Contains(p Position) bool {
	if !s.IsValid() || !p.IsValid() {
		return false
	}
	
	// Check if position is between start and end
	if p.Line < s.Start.Line || p.Line > s.End.Line {
		return false
	}
	
	if p.Line == s.Start.Line && p.Column < s.Start.Column {
		return false
	}
	
	if p.Line == s.End.Line && p.Column > s.End.Column {
		return false
	}
	
	return true
}

// Overlaps checks if this span overlaps with another span
func (s Span) Overlaps(other Span) bool {
	if !s.IsValid() || !other.IsValid() {
		return false
	}
	
	return s.Contains(other.Start) || s.Contains(other.End) ||
		other.Contains(s.Start) || other.Contains(s.End)
}

// Merge creates a new span covering both this and another span
func (s Span) Merge(other Span) Span {
	if !s.IsValid() {
		return other
	}
	if !other.IsValid() {
		return s
	}
	
	start := s.Start
	if other.Start.Line < start.Line ||
		(other.Start.Line == start.Line && other.Start.Column < start.Column) {
		start = other.Start
	}
	
	end := s.End
	if other.End.Line > end.Line ||
		(other.End.Line == end.Line && other.End.Column > end.Column) {
		end = other.End
	}
	
	file := s.File
	if file == "" {
		file = other.File
	}
	
	return NewSpan(start, end, file)
}


