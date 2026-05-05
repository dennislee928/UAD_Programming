package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/dennislee928/uad-lang/internal/lsp/completion"
	"github.com/dennislee928/uad-lang/internal/lsp/hover"
	"github.com/dennislee928/uad-lang/internal/lsp/protocol"
)

// Server is the main LSP server instance
type Server struct {
	// State
	initialized bool
	shutdown    bool
	
	// Document management
	documents *DocumentManager
	
	// Analysis
	analyzer *Analyzer
	
	// Completion
	completionEngine *completion.Engine
	
	// Hover
	hoverProvider *hover.Provider
	
	// Mutex for thread safety
	mu sync.RWMutex
	
	// Request/response tracking
	nextID int
	
	// Context
	ctx    context.Context
	cancel context.CancelFunc
}

// NewServer creates a new LSP server instance
func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Server{
		documents:        NewDocumentManager(),
		analyzer:         NewAnalyzer(),
		completionEngine: completion.NewEngine(),
		hoverProvider:    hover.NewProvider(),
		ctx:              ctx,
		cancel:           cancel,
	}
}

// RunStdio runs the server using stdio for communication
func (s *Server) RunStdio() error {
	log.Println("LSP Server starting with stdio transport")
	
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout
	
	for {
		// Check if shutdown
		if s.shutdown {
			log.Println("Server shutting down")
			break
		}
		
		// Read message
		msg, err := s.readMessage(reader)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				break
			}
			log.Printf("Error reading message: %v", err)
			continue
		}
		
		// Handle message
		response, err := s.handleMessage(msg)
		if err != nil {
			log.Printf("Error handling message: %v", err)
			// Send error response
			s.sendError(writer, msg, err)
			continue
		}
		
		// Send response if not nil (notifications don't need responses)
		if response != nil {
			if err := s.sendMessage(writer, response); err != nil {
				log.Printf("Error sending response: %v", err)
			}
		}
	}
	
	return nil
}

// readMessage reads a JSON-RPC message from the reader
func (s *Server) readMessage(reader *bufio.Reader) (map[string]interface{}, error) {
	// Read headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		
		// Empty line marks end of headers
		if line == "\r\n" || line == "\n" {
			break
		}
		
		// Parse header
		var key, value string
		fmt.Sscanf(line, "%s %s", &key, &value)
		headers[key] = value
	}
	
	// Get content length
	contentLengthStr, ok := headers["Content-Length:"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}
	
	var contentLength int
	fmt.Sscanf(contentLengthStr, "%d", &contentLength)
	
	// Read content
	content := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, content); err != nil {
		return nil, err
	}
	
	// Parse JSON
	var msg map[string]interface{}
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, err
	}
	
	return msg, nil
}

// sendMessage sends a JSON-RPC message to the writer
func (s *Server) sendMessage(writer io.Writer, msg map[string]interface{}) error {
	// Marshal JSON
	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	// Write headers
	headers := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	if _, err := writer.Write([]byte(headers)); err != nil {
		return err
	}
	
	// Write content
	if _, err := writer.Write(content); err != nil {
		return err
	}
	
	return nil
}

// sendError sends an error response
func (s *Server) sendError(writer io.Writer, request map[string]interface{}, err error) error {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      request["id"],
		"error": map[string]interface{}{
			"code":    -32603, // Internal error
			"message": err.Error(),
		},
	}
	
	return s.sendMessage(writer, response)
}

// handleMessage routes messages to appropriate handlers
func (s *Server) handleMessage(msg map[string]interface{}) (map[string]interface{}, error) {
	method, ok := msg["method"].(string)
	if !ok {
		return nil, fmt.Errorf("missing method field")
	}
	
	log.Printf("Handling method: %s", method)
	
	// Get params
	params, _ := msg["params"].(map[string]interface{})
	
	// Get ID (nil for notifications)
	id := msg["id"]
	
	// Route to handler
	switch method {
	case "initialize":
		return s.handleInitialize(id, params)
	
	case "initialized":
		// Notification, no response needed
		s.handleInitialized()
		return nil, nil
	
	case "shutdown":
		return s.handleShutdown(id)
	
	case "exit":
		s.handleExit()
		return nil, nil
	
	case "textDocument/didOpen":
		s.handleDidOpen(params)
		return nil, nil
	
	case "textDocument/didChange":
		s.handleDidChange(params)
		return nil, nil
	
	case "textDocument/didClose":
		s.handleDidClose(params)
		return nil, nil
	
	case "textDocument/completion":
		return s.handleCompletion(id, params)
	
	case "textDocument/hover":
		return s.handleHover(id, params)
	
	default:
		log.Printf("Unhandled method: %s", method)
		return s.successResponse(id, nil), nil
	}
}

// successResponse creates a success response
func (s *Server) successResponse(id interface{}, result interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	log.Println("Initialize request received")
	
	capabilities := map[string]interface{}{
		"textDocumentSync": map[string]interface{}{
			"openClose": true,
			"change":    protocol.TextDocumentSyncKindIncremental,
		},
		"completionProvider": map[string]interface{}{
			"triggerCharacters": []string{".", ":", "::"},
		},
		"hoverProvider": true,
		// More capabilities will be added in future
	}
	
	result := map[string]interface{}{
		"capabilities": capabilities,
		"serverInfo": map[string]interface{}{
			"name":    "UAD Language Server",
			"version": "0.1.0",
		},
	}
	
	s.mu.Lock()
	s.initialized = true
	s.mu.Unlock()
	
	return s.successResponse(id, result), nil
}

// handleInitialized handles the initialized notification
func (s *Server) handleInitialized() {
	log.Println("Client initialized")
}

// handleShutdown handles the shutdown request
func (s *Server) handleShutdown(id interface{}) (map[string]interface{}, error) {
	log.Println("Shutdown request received")
	
	s.mu.Lock()
	s.shutdown = true
	s.mu.Unlock()
	
	return s.successResponse(id, nil), nil
}

// handleExit handles the exit notification
func (s *Server) handleExit() {
	log.Println("Exit notification received")
	s.cancel()
	os.Exit(0)
}

// handleDidOpen handles textDocument/didOpen notification
func (s *Server) handleDidOpen(params map[string]interface{}) {
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		log.Println("Invalid textDocument in didOpen")
		return
	}
	
	uri, _ := textDocument["uri"].(string)
	text, _ := textDocument["text"].(string)
	version, _ := textDocument["version"].(float64)
	
	log.Printf("Document opened: %s", uri)
	
	// Add document
	s.documents.DidOpen(uri, text, int(version))
	
	// Analyze document and publish diagnostics
	go s.analyzeAndPublishDiagnostics(uri)
}

// handleDidChange handles textDocument/didChange notification
func (s *Server) handleDidChange(params map[string]interface{}) {
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		return
	}
	
	uri, _ := textDocument["uri"].(string)
	version, _ := textDocument["version"].(float64)
	
	contentChanges, ok := params["contentChanges"].([]interface{})
	if !ok || len(contentChanges) == 0 {
		return
	}
	
	// Get full text from first change (incremental sync not yet implemented)
	change := contentChanges[0].(map[string]interface{})
	text, _ := change["text"].(string)
	
	log.Printf("Document changed: %s", uri)
	
	// Update document
	s.documents.DidChange(uri, text, int(version))
	
	// Re-analyze and send diagnostics
	go s.analyzeAndPublishDiagnostics(uri)
}

// handleDidClose handles textDocument/didClose notification
func (s *Server) handleDidClose(params map[string]interface{}) {
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		return
	}
	
	uri, _ := textDocument["uri"].(string)
	
	log.Printf("Document closed: %s", uri)
	
	s.documents.DidClose(uri)
}

// handleCompletion handles textDocument/completion request
func (s *Server) handleCompletion(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	log.Println("Completion request received")
	
	// Extract position
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		return s.successResponse(id, []interface{}{}), nil
	}
	
	uri, _ := textDocument["uri"].(string)
	position, _ := params["position"].(map[string]interface{})
	line := int(position["line"].(float64))
	character := int(position["character"].(float64))
	
	// Get document
	doc, exists := s.documents.GetDocument(uri)
	if !exists {
		return s.successResponse(id, []interface{}{}), nil
	}
	
	// Get completions
	completionItems := s.completionEngine.Complete(doc, line, character)
	
	// Convert to LSP format
	lspItems := make([]interface{}, len(completionItems))
	for i, item := range completionItems {
		lspItems[i] = map[string]interface{}{
			"label":         item.Label,
			"kind":          item.Kind,
			"detail":        item.Detail,
			"documentation": item.Documentation,
			"insertText":    item.InsertText,
		}
	}
	
	log.Printf("Returning %d completion items", len(lspItems))
	
	return s.successResponse(id, lspItems), nil
}

// handleHover handles textDocument/hover request
func (s *Server) handleHover(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	log.Println("Hover request received")
	
	// Extract position
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		return s.successResponse(id, nil), nil
	}
	
	uri, _ := textDocument["uri"].(string)
	position, _ := params["position"].(map[string]interface{})
	line := int(position["line"].(float64))
	character := int(position["character"].(float64))
	
	// Get document
	doc, exists := s.documents.GetDocument(uri)
	if !exists || doc.AST == nil {
		return s.successResponse(id, nil), nil
	}
	
	// Get hover info
	hoverInfo := s.hoverProvider.GetHover(doc.AST, line, character)
	if hoverInfo == nil {
		return s.successResponse(id, nil), nil
	}
	
	// Convert to LSP format
	result := map[string]interface{}{
		"contents": map[string]interface{}{
			"kind":  "markdown",
			"value": hoverInfo.Contents,
		},
	}
	
	if hoverInfo.Range != nil {
		result["range"] = map[string]interface{}{
			"start": map[string]interface{}{
				"line":      hoverInfo.Range.Start.Line,
				"character": hoverInfo.Range.Start.Character,
			},
			"end": map[string]interface{}{
				"line":      hoverInfo.Range.End.Line,
				"character": hoverInfo.Range.End.Character,
			},
		}
	}
	
	return s.successResponse(id, result), nil
}

// analyzeAndPublishDiagnostics analyzes a document and publishes diagnostics
func (s *Server) analyzeAndPublishDiagnostics(uri string) {
	doc, exists := s.documents.GetDocument(uri)
	if !exists {
		return
	}
	
	// Analyze document
	diagnostics := s.analyzer.Analyze(doc)
	
	// Store diagnostics in document
	doc.Diagnostics = diagnostics
	
	// Publish diagnostics to client
	s.publishDiagnostics(uri, diagnostics)
}

// publishDiagnostics sends diagnostics to the client
func (s *Server) publishDiagnostics(uri string, diagnostics []protocol.Diagnostic) {
	// Convert diagnostics to LSP format
	lspDiagnostics := make([]interface{}, len(diagnostics))
	for i, diag := range diagnostics {
		lspDiagnostics[i] = map[string]interface{}{
			"range": map[string]interface{}{
				"start": map[string]interface{}{
					"line":      diag.Range.Start.Line,
					"character": diag.Range.Start.Character,
				},
				"end": map[string]interface{}{
					"line":      diag.Range.End.Line,
					"character": diag.Range.End.Character,
				},
			},
			"severity": diag.Severity,
			"source":   diag.Source,
			"message":  diag.Message,
		}
	}
	
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]interface{}{
			"uri":         uri,
			"diagnostics": lspDiagnostics,
		},
	}
	
	// Send notification to stdout
	if err := s.sendMessage(os.Stdout, notification); err != nil {
		log.Printf("Error publishing diagnostics: %v", err)
	}
}

