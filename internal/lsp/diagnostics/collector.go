package diagnostics

import (
	"github.com/dennislee928/uad-lang/internal/common"
	"github.com/dennislee928/uad-lang/internal/lsp/protocol"
)

// Collector collects diagnostics from various sources
type Collector struct{}

// NewCollector creates a new diagnostics collector
func NewCollector() *Collector {
	return &Collector{}
}

// FromErrorList converts common.ErrorList to LSP diagnostics
func (c *Collector) FromErrorList(errors *common.ErrorList) []protocol.Diagnostic {
	if errors == nil || len(errors.Errors) == 0 {
		return []protocol.Diagnostic{}
	}
	
	diagnostics := make([]protocol.Diagnostic, 0, len(errors.Errors))
	
	for _, err := range errors.Errors {
		severity := protocol.DiagnosticSeverityError
		
		// Map error kind to severity
		switch err.Kind {
		case common.ErrSyntax:
			severity = protocol.DiagnosticSeverityError
		case common.ErrType:
			severity = protocol.DiagnosticSeverityError
		case common.ErrSemantic:
			severity = protocol.DiagnosticSeverityError
		default:
			severity = protocol.DiagnosticSeverityError
		}
		
		diag := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      err.Span.Start.Line - 1, // Convert to 0-based
					Character: err.Span.Start.Column - 1,
				},
				End: protocol.Position{
					Line:      err.Span.End.Line - 1,
					Character: err.Span.End.Column - 1,
				},
			},
			Severity: severity,
			Source:   errorKindToSource(err.Kind),
			Message:  err.Message,
		}
		
		// Add note as related information if present
		if err.Note != "" {
			diag.Message += "\nNote: " + err.Note
		}
		
		diagnostics = append(diagnostics, diag)
	}
	
	return diagnostics
}

func errorKindToSource(kind common.ErrorKind) string {
	switch kind {
	case common.ErrSyntax:
		return "parser"
	case common.ErrType:
		return "type-checker"
	case common.ErrSemantic:
		return "semantic-analyzer"
	default:
		return "uad"
	}
}

