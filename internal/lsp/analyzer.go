package lsp

import (
	"github.com/dennislee928/uad-lang/internal/common"
	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/lsp/diagnostics"
	"github.com/dennislee928/uad-lang/internal/lsp/protocol"
	"github.com/dennislee928/uad-lang/internal/parser"
	"github.com/dennislee928/uad-lang/internal/typer"
)

// Analyzer performs code analysis
type Analyzer struct {
	// Configuration
	config AnalyzerConfig
}

// AnalyzerConfig holds analyzer configuration
type AnalyzerConfig struct {
	EnableTypeChecking bool
	EnableLinting      bool
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		config: AnalyzerConfig{
			EnableTypeChecking: true,
			EnableLinting:      true,
		},
	}
}

// Analyze performs full analysis on a document
func (a *Analyzer) Analyze(doc *Document) []protocol.Diagnostic {
	allDiagnostics := []protocol.Diagnostic{}
	collector := diagnostics.NewCollector()
	
	// Tokenize
	l := lexer.New(doc.Content, doc.URI)
	tokens := []lexer.Token{}
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lexer.TokenEOF {
			break
		}
	}
	
	// Parse
	p := parser.New(tokens, doc.URI)
	module, err := p.ParseModule()
	if err != nil {
		// Convert parser errors to diagnostics
		if errorList, ok := err.(*common.ErrorList); ok {
			parseErrors := collector.FromErrorList(errorList)
			allDiagnostics = append(allDiagnostics, parseErrors...)
		} else {
			// Fallback for unexpected error types
			allDiagnostics = append(allDiagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 1},
				},
				Severity: protocol.DiagnosticSeverityError,
				Source:   "parser",
				Message:  err.Error(),
			})
		}
		
		return allDiagnostics
	}
	
	// Cache AST
	doc.AST = module
	
	// Type check
	if a.config.EnableTypeChecking {
		tc := typer.NewTypeChecker()
		if err := tc.Check(module); err != nil {
			// Convert type errors to diagnostics
			if errorList, ok := err.(*common.ErrorList); ok {
				typeErrors := collector.FromErrorList(errorList)
				allDiagnostics = append(allDiagnostics, typeErrors...)
			} else {
				// Fallback
				allDiagnostics = append(allDiagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 0},
						End:   protocol.Position{Line: 0, Character: 1},
					},
					Severity: protocol.DiagnosticSeverityError,
					Source:   "type-checker",
					Message:  err.Error(),
				})
			}
		}
	}
	
	// TODO: Additional linting checks
	// - Unused variables
	// - Unreachable code
	// - Style violations
	
	return allDiagnostics
}

