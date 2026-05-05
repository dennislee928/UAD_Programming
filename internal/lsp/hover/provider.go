package hover

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/lsp/protocol"
)

// Provider provides hover information
type Provider struct{}

// NewProvider creates a new hover provider
func NewProvider() *Provider {
	return &Provider{}
}

// HoverInfo represents hover information
type HoverInfo struct {
	Contents string
	Range    *protocol.Range
}

// GetHover provides hover information at a position
func (p *Provider) GetHover(module *ast.Module, line, character int) *HoverInfo {
	// TODO: Implement proper hover by:
	// 1. Finding the symbol at the position
	// 2. Looking up its type information
	// 3. Generating documentation
	
	// For now, return nil (no hover)
	return nil
}

// GetSymbolHover returns hover info for a specific symbol
func (p *Provider) GetSymbolHover(symbol string, symbolType string) *HoverInfo {
	// Build hover content
	content := fmt.Sprintf("**%s**\n\n```uad\n%s\n```", symbol, symbolType)
	
	return &HoverInfo{
		Contents: content,
	}
}

// GetKeywordHover returns hover info for keywords
func (p *Provider) GetKeywordHover(keyword string) *HoverInfo {
	docs := getKeywordDocumentation(keyword)
	if docs == "" {
		return nil
	}
	
	return &HoverInfo{
		Contents: docs,
	}
}

// getKeywordDocumentation returns documentation for keywords
func getKeywordDocumentation(keyword string) string {
	keywordDocs := map[string]string{
		"fn":       "**fn** - Function declaration\n\nDefine a new function.",
		"let":      "**let** - Variable binding\n\nBind a value to an immutable variable.",
		"mut":      "**mut** - Mutable variable\n\nDeclare a mutable variable.",
		"struct":   "**struct** - Structure type\n\nDefine a composite data type.",
		"enum":     "**enum** - Enumeration type\n\nDefine a type with multiple variants.",
		"if":       "**if** - Conditional expression\n\nConditionally execute code.",
		"match":    "**match** - Pattern matching\n\nMatch a value against patterns.",
		"for":      "**for** - For loop\n\nIterate over a collection.",
		"while":    "**while** - While loop\n\nLoop while a condition is true.",
		"return":   "**return** - Return from function\n\nReturn a value from a function.",
		"score":    "**score** - Musical score (Musical DSL)\n\nDefine a time-structured multi-agent score.",
		"track":    "**track** - Track declaration (Musical DSL)\n\nDefine an agent's timeline within a score.",
		"motif":    "**motif** - Motif declaration (Musical DSL)\n\nDefine a reusable behavior pattern.",
		"string":   "**string** - String field (String Theory)\n\nDefine a string field with modes.",
		"brane":    "**brane** - Brane context (String Theory)\n\nDefine a dimensional context for string fields.",
		"coupling": "**coupling** - Coupling declaration (String Theory)\n\nDefine a coupling between string modes.",
		"entangle": "**entangle** - Entanglement (Quantum Semantics)\n\nCreate quantum entanglement between variables.",
		"Int":      "**Int** - Integer type\n\nSigned 64-bit integer.",
		"Float":    "**Float** - Floating-point type\n\n64-bit floating-point number.",
		"Bool":     "**Bool** - Boolean type\n\nTrue or false value.",
		"String":   "**String** - String type\n\nUTF-8 encoded text.",
	}
	
	return keywordDocs[keyword]
}

// GetBuiltinFunctionHover returns hover info for builtin functions
func (p *Provider) GetBuiltinFunctionHover(funcName string) *HoverInfo {
	builtinDocs := map[string]string{
		"println": "**println**(value: Any) -> Unit\n\nPrint a value followed by a newline.",
		"print":   "**print**(value: Any) -> Unit\n\nPrint a value without a newline.",
		"len":     "**len**(collection: Array|Map|String) -> Int\n\nReturn the length of a collection.",
		"push":    "**push**(array: Array<T>, value: T) -> Unit\n\nAppend a value to an array.",
		"pop":     "**pop**(array: Array<T>) -> T?\n\nRemove and return the last element.",
	}
	
	docs, exists := builtinDocs[funcName]
	if !exists {
		return nil
	}
	
	return &HoverInfo{
		Contents: docs,
	}
}


