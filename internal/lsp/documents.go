package lsp

import (
	"sync"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/lsp/protocol"
)

// Document represents an open document in the editor
type Document struct {
	URI     string
	Content string
	Version int
	
	// Parsed AST (cached)
	AST *ast.Module
	
	// Analysis results
	Diagnostics []protocol.Diagnostic
}

// DocumentManager manages all open documents
type DocumentManager struct {
	documents map[string]*Document
	mu        sync.RWMutex
}

// NewDocumentManager creates a new document manager
func NewDocumentManager() *DocumentManager {
	return &DocumentManager{
		documents: make(map[string]*Document),
	}
}

// DidOpen handles document open event
func (dm *DocumentManager) DidOpen(uri, content string, version int) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	
	dm.documents[uri] = &Document{
		URI:     uri,
		Content: content,
		Version: version,
	}
}

// DidChange handles document change event
func (dm *DocumentManager) DidChange(uri, content string, version int) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	
	if doc, exists := dm.documents[uri]; exists {
		doc.Content = content
		doc.Version = version
		// Invalidate cached AST
		doc.AST = nil
	}
}

// DidClose handles document close event
func (dm *DocumentManager) DidClose(uri string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	
	delete(dm.documents, uri)
}

// GetDocument retrieves a document by URI
func (dm *DocumentManager) GetDocument(uri string) (*Document, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	
	doc, exists := dm.documents[uri]
	return doc, exists
}

// AllDocuments returns all open documents
func (dm *DocumentManager) AllDocuments() []*Document {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	
	docs := make([]*Document, 0, len(dm.documents))
	for _, doc := range dm.documents {
		docs = append(docs, doc)
	}
	
	return docs
}

