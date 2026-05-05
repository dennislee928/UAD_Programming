package common

import (
	"fmt"
	"strings"
)

// ErrorKind represents the category of an error
type ErrorKind int

const (
	// Lexical errors
	ErrLexical ErrorKind = iota
	
	// Syntax errors
	ErrSyntax
	
	// Type errors
	ErrType
	
	// Semantic errors
	ErrSemantic
	
	// Runtime errors
	ErrRuntime
	
	// Internal errors (compiler bugs)
	ErrInternal
)

// String returns the string representation of ErrorKind
func (k ErrorKind) String() string {
	switch k {
	case ErrLexical:
		return "Lexical Error"
	case ErrSyntax:
		return "Syntax Error"
	case ErrType:
		return "Type Error"
	case ErrSemantic:
		return "Semantic Error"
	case ErrRuntime:
		return "Runtime Error"
	case ErrInternal:
		return "Internal Error"
	default:
		return "Unknown Error"
	}
}

// Error represents a compiler or runtime error with position information
type Error struct {
	Kind    ErrorKind
	Message string
	Span    Span
	Note    string // Optional note or suggestion
}

// NewError creates a new error
func NewError(kind ErrorKind, message string, span Span) *Error {
	return &Error{
		Kind:    kind,
		Message: message,
		Span:    span,
	}
}

// NewErrorWithNote creates a new error with a note
func NewErrorWithNote(kind ErrorKind, message string, span Span, note string) *Error {
	return &Error{
		Kind:    kind,
		Message: message,
		Span:    span,
		Note:    note,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	var sb strings.Builder
	
	sb.WriteString(e.Kind.String())
	sb.WriteString(": ")
	sb.WriteString(e.Message)
	
	if e.Span.IsValid() {
		sb.WriteString("\n  --> ")
		sb.WriteString(e.Span.String())
	}
	
	if e.Note != "" {
		sb.WriteString("\n  note: ")
		sb.WriteString(e.Note)
	}
	
	return sb.String()
}

// ErrorList represents a collection of errors
type ErrorList struct {
	Errors []*Error
}

// NewErrorList creates a new error list
func NewErrorList() *ErrorList {
	return &ErrorList{
		Errors: make([]*Error, 0),
	}
}

// Add adds an error to the list
func (l *ErrorList) Add(err *Error) {
	l.Errors = append(l.Errors, err)
}

// AddError creates and adds an error to the list
func (l *ErrorList) AddError(kind ErrorKind, message string, span Span) {
	l.Add(NewError(kind, message, span))
}

// HasErrors returns true if the list contains any errors
func (l *ErrorList) HasErrors() bool {
	return len(l.Errors) > 0
}

// Count returns the number of errors
func (l *ErrorList) Count() int {
	return len(l.Errors)
}

// Error implements the error interface
func (l *ErrorList) Error() string {
	if len(l.Errors) == 0 {
		return "no errors"
	}
	
	var sb strings.Builder
	for i, err := range l.Errors {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Sort sorts errors by position
func (l *ErrorList) Sort() {
	// Simple bubble sort (fine for small error lists)
	n := len(l.Errors)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			e1 := l.Errors[j]
			e2 := l.Errors[j+1]
			
			if e1.Span.Start.Line > e2.Span.Start.Line ||
				(e1.Span.Start.Line == e2.Span.Start.Line &&
					e1.Span.Start.Column > e2.Span.Start.Column) {
				l.Errors[j], l.Errors[j+1] = e2, e1
			}
		}
	}
}

// Common error constructors

// LexicalError creates a lexical error
func LexicalError(message string, span Span) *Error {
	return NewError(ErrLexical, message, span)
}

// SyntaxError creates a syntax error
func SyntaxError(message string, span Span) *Error {
	return NewError(ErrSyntax, message, span)
}

// TypeError creates a type error
func TypeError(message string, span Span) *Error {
	return NewError(ErrType, message, span)
}

// SemanticError creates a semantic error
func SemanticError(message string, span Span) *Error {
	return NewError(ErrSemantic, message, span)
}

// RuntimeError creates a runtime error
func RuntimeError(message string, span Span) *Error {
	return NewError(ErrRuntime, message, span)
}

// InternalError creates an internal error
func InternalError(message string, span Span) *Error {
	return NewError(ErrInternal, message, span)
}

// Errorf creates a formatted error message
func Errorf(kind ErrorKind, span Span, format string, args ...interface{}) *Error {
	return NewError(kind, fmt.Sprintf(format, args...), span)
}


