package completion

import (
	"strings"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/lsp/protocol"
)

// Engine provides code completion
type Engine struct {
	// Completion context
	keywords []CompletionItem
	snippets []CompletionItem
}

// CompletionItem represents a completion item
type CompletionItem struct {
	Label         string
	Kind          protocol.CompletionItemKind
	Detail        string
	Documentation string
	InsertText    string
}

// NewEngine creates a new completion engine
func NewEngine() *Engine {
	return &Engine{
		keywords: buildKeywordCompletions(),
		snippets: buildSnippetCompletions(),
	}
}

// Complete provides completions for a given position
func (e *Engine) Complete(doc interface{}, line, character int) []CompletionItem {
	items := []CompletionItem{}
	
	// Add keywords
	items = append(items, e.keywords...)
	
	// Add snippets
	items = append(items, e.snippets...)
	
	// TODO: Add context-aware completions
	// - Variables in scope
	// - Function names
	// - Struct fields
	// - Enum variants
	
	return items
}

// CompleteWithPrefix filters completions by prefix
func (e *Engine) CompleteWithPrefix(doc interface{}, line, character int, prefix string) []CompletionItem {
	allItems := e.Complete(doc, line, character)
	
	if prefix == "" {
		return allItems
	}
	
	// Filter by prefix
	filtered := []CompletionItem{}
	lowerPrefix := strings.ToLower(prefix)
	
	for _, item := range allItems {
		if strings.HasPrefix(strings.ToLower(item.Label), lowerPrefix) {
			filtered = append(filtered, item)
		}
	}
	
	return filtered
}

// buildKeywordCompletions creates completion items for keywords
func buildKeywordCompletions() []CompletionItem {
	keywords := []string{
		"fn", "let", "mut", "if", "else", "while", "for", "in",
		"match", "return", "break", "continue",
		"struct", "enum", "type", "import",
		"pub", "priv", "as",
		// Musical DSL
		"score", "track", "bars", "motif", "use", "at",
		"variation", "transpose", "tempo", "time_signature",
		"beats", "duration",
		// String Theory
		"string", "brane", "on", "coupling", "resonance",
		"modes", "with", "strength", "dimensions",
		"when", "and", "freq", "amplitude", "threshold", "mark",
		// Entanglement
		"entangle",
		// Types
		"Int", "Float", "Bool", "String", "Unit", "Duration", "Time",
		// Constants
		"true", "false", "nil",
	}
	
	items := make([]CompletionItem, len(keywords))
	for i, kw := range keywords {
		items[i] = CompletionItem{
			Label:  kw,
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "keyword",
		}
	}
	
	return items
}

// buildSnippetCompletions creates snippet completions
func buildSnippetCompletions() []CompletionItem {
	return []CompletionItem{
		{
			Label:         "fn",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Function declaration",
			Documentation: "Create a new function",
			InsertText:    "fn ${1:name}(${2:params}) -> ${3:Type} {\n\t$0\n}",
		},
		{
			Label:         "struct",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Struct declaration",
			Documentation: "Create a new struct",
			InsertText:    "struct ${1:Name} {\n\t${2:field}: ${3:Type},\n}",
		},
		{
			Label:         "enum",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Enum declaration",
			Documentation: "Create a new enum",
			InsertText:    "enum ${1:Name} {\n\t${2:Variant},\n}",
		},
		{
			Label:         "if",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "If statement",
			InsertText:    "if ${1:condition} {\n\t$0\n}",
		},
		{
			Label:         "for",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "For loop",
			InsertText:    "for ${1:item} in ${2:items} {\n\t$0\n}",
		},
		{
			Label:         "match",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Match expression",
			InsertText:    "match ${1:value} {\n\t${2:pattern} => $0,\n}",
		},
		{
			Label:         "score",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Score declaration (Musical DSL)",
			Documentation: "Create a musical score",
			InsertText:    "score ${1:Name} {\n\ttempo: ${2:120},\n\ttime_signature: ${3:4/4},\n\t$0\n}",
		},
		{
			Label:         "track",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Track declaration (Musical DSL)",
			InsertText:    "track ${1:Name} {\n\t$0\n}",
		},
		{
			Label:         "motif",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Motif declaration (Musical DSL)",
			InsertText:    "motif ${1:name} {\n\tduration: ${2:4} beats,\n\t$0\n}",
		},
		{
			Label:         "string",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "String field declaration (String Theory)",
			InsertText:    "string ${1:Name} {\n\tmodes {\n\t\t${2:mode}: ${3:Float},\n\t}\n}",
		},
		{
			Label:         "brane",
			Kind:          protocol.CompletionItemKindSnippet,
			Detail:        "Brane declaration (String Theory)",
			InsertText:    "brane ${1:Name} {\n\tdimensions [${2:time, space}]\n}",
		},
	}
}

// GetSymbolsInScope retrieves symbols visible at a position
// TODO: Implement proper scope analysis
func (e *Engine) GetSymbolsInScope(module *ast.Module, line, character int) []CompletionItem {
	items := []CompletionItem{}
	
	// TODO: Walk AST to find:
	// - Function declarations
	// - Variable declarations in scope
	// - Struct/Enum types
	// - Import names
	
	return items
}

