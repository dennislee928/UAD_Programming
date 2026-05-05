package protocol

// TextDocumentSyncKind defines how text documents are synced
type TextDocumentSyncKind int

const (
	// TextDocumentSyncKindNone - Documents should not be synced
	TextDocumentSyncKindNone TextDocumentSyncKind = 0
	
	// TextDocumentSyncKindFull - Documents are synced by always sending the full content
	TextDocumentSyncKindFull TextDocumentSyncKind = 1
	
	// TextDocumentSyncKindIncremental - Documents are synced by sending the full content on open
	// After that only incremental updates are sent
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

// CompletionItemKind defines completion item types
type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

// Diagnostic represents a diagnostic message
type Diagnostic struct {
	Range    Range
	Severity DiagnosticSeverity
	Code     string
	Source   string
	Message  string
}

// DiagnosticSeverity levels
type DiagnosticSeverity int

const (
	DiagnosticSeverityError       DiagnosticSeverity = 1
	DiagnosticSeverityWarning     DiagnosticSeverity = 2
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	DiagnosticSeverityHint        DiagnosticSeverity = 4
)

// Range in a document
type Range struct {
	Start Position
	End   Position
}

// Position in a document
type Position struct {
	Line      int
	Character int
}

