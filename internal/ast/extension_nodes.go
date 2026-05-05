package ast

import (
	"github.com/dennislee928/uad-lang/internal/common"
)

// extension_nodes.go defines AST nodes for future language extensions.
// These interfaces and placeholders reserve space for:
// - Musical DSL (score, track, bar, motif, variation)
// - String Theory semantics (string, brane, coupling, resonance)
// - Quantum Entanglement (entangle, synchronization)
//
// Implementation status: Interfaces defined, concrete types to be implemented in M2.3-M2.5

// ==================== Extension Interfaces ====================

// MusicalNode represents nodes in the musical DSL layer.
// Musical nodes introduce temporal structure and harmony-based semantics.
type MusicalNode interface {
	Node
	musicalNode() // Marker method
}

// StringTheoryNode represents nodes for string theory semantics.
// String theory nodes model field coupling and resonance.
type StringTheoryNode interface {
	Node
	stringTheoryNode() // Marker method
}

// EntanglementNode represents nodes for quantum entanglement semantics.
// Entanglement nodes describe variable synchronization and shared state.
type EntanglementNode interface {
	Node
	entanglementNode() // Marker method
}

// ==================== Musical DSL Nodes (M2.3) ====================

// ScoreNode represents a musical score (top-level temporal structure).
// A score contains multiple tracks and defines the overall temporal grid.
// Planned syntax: `score <name> { ... }`
// Status: Interface placeholder (implementation in M2.3.1)
type ScoreNode struct {
	baseNode
	Name   *Ident  // Score name
	Tracks []*TrackNode // Tracks within this score
	// TODO(M2.3.1): Add tempo, time signature, etc.
}

func (s *ScoreNode) declNode() {}
func (s *ScoreNode) musicalNode() {}

// TrackNode represents a single track (voice/agent timeline).
// A track contains bars and motifs arranged in time.
// Planned syntax: `track <name> { ... }`
// Status: Interface placeholder (implementation in M2.3.1)
type TrackNode struct {
	baseNode
	Name *Ident // Track name
	Bars []*BarRangeNode // Bar ranges in this track
	// TODO(M2.3.1): Add instrument, agent binding, etc.
}

func (t *TrackNode) musicalNode() {}

// BarRangeNode represents a time range in bars.
// Defines what happens within a specific bar range.
// Planned syntax: `bars <start>..<end> { ... }`
// Status: Interface placeholder (implementation in M2.3.1)
type BarRangeNode struct {
	baseNode
	StartBar int  // Start bar number
	EndBar   int  // End bar number
	Body     *BlockExpr // Actions within this bar range
	// TODO(M2.3.1): Add beat subdivision, duration, etc.
}

func (b *BarRangeNode) musicalNode() {}
func (b *BarRangeNode) stmtNode() {} // Allow bars as statements for nesting

// MotifDeclNode represents a motif (reusable musical pattern) declaration.
// Motifs are like functions but for temporal patterns.
// Planned syntax: `motif <name> { ... }`
// Status: Interface placeholder (implementation in M2.3.1)
type MotifDeclNode struct {
	baseNode
	Name       *Ident     // Motif name
	Params     []*Param   // Motif parameters (for variations)
	Body       *BlockExpr // Motif body
	// TODO(M2.3.1): Add duration, key signature, etc.
}

func (m *MotifDeclNode) declNode() {}
func (m *MotifDeclNode) musicalNode() {}

// EmitStmt represents an emit statement for events.
// Emits an event of a specific type with given fields.
// Syntax: `emit <TypeName> { <fields>... };`
type EmitStmt struct {
	baseNode
	TypeName *Ident          // Event type name (e.g., "Event")
	Fields   *StructLiteral  // Event fields as struct literal
}

func (e *EmitStmt) stmtNode() {}
func (e *EmitStmt) musicalNode() {}

// NewEmitStmt creates a new emit statement.
func NewEmitStmt(typeName *Ident, fields *StructLiteral, span common.Span) *EmitStmt {
	return &EmitStmt{
		baseNode: baseNode{span},
		TypeName: typeName,
		Fields:   fields,
	}
}

// UseStmt represents a use statement for calling motifs.
// Syntax: `use <motif_name>;` or `use <motif_name>(<args>);`
type UseStmt struct {
	baseNode
	MotifName *Ident    // Motif name to use
	Args      []Expr    // Arguments for parameterized motifs (optional)
}

func (u *UseStmt) stmtNode() {}
func (u *UseStmt) musicalNode() {}

// NewUseStmt creates a new use statement.
func NewUseStmt(motifName *Ident, args []Expr, span common.Span) *UseStmt {
	return &UseStmt{
		baseNode: baseNode{span},
		MotifName: motifName,
		Args:      args,
	}
}

// MotifUseNode represents the use/instantiation of a motif.
// Applies a motif at a specific bar range.
// Planned syntax: `use motif <name> at bars <range>`
// Status: Interface placeholder (implementation in M2.3.1)
type MotifUseNode struct {
	baseNode
	MotifName *Ident       // Name of motif to use
	BarRange  *BarRangeNode // Where to apply the motif
	Args      []Expr       // Arguments for parameterized motifs
	// TODO(M2.3.1): Add transposition, variation parameters, etc.
}

func (m *MotifUseNode) musicalNode() {}

// MotifVariationNode represents a variation on a motif.
// Creates a new motif by transforming an existing one.
// Planned syntax: `variation motif <name> as <new_name> { transpose ... }`
// Status: Interface placeholder (implementation in M2.3.1)
type MotifVariationNode struct {
	baseNode
	OriginalName *Ident      // Original motif name
	VariantName  *Ident      // New variant name
	Transforms   []Expr      // List of transformations
	// TODO(M2.3.1): Add transpose, augmentation, diminution, etc.
}

func (m *MotifVariationNode) declNode() {}
func (m *MotifVariationNode) musicalNode() {}

// ==================== String Theory Nodes (M2.4) ====================

// StringDeclNode represents a string field declaration.
// Strings have multiple modes (vibrational degrees of freedom).
// Planned syntax: `string <name> { modes { ... } }`
// Status: Interface placeholder (implementation in M2.4.1)
type StringDeclNode struct {
	baseNode
	Name  *Ident   // String name
	Modes []*Field // Mode fields (name: type pairs)
	// TODO(M2.4.1): Add frequency, amplitude, phase, etc.
}

func (s *StringDeclNode) declNode() {}
func (s *StringDeclNode) stringTheoryNode() {}

// BraneDeclNode represents a brane (background space) declaration.
// Branes provide the dimensional context for strings.
// Planned syntax: `brane <name> { dimensions [...] }`
// Status: Interface placeholder (implementation in M2.4.1)
type BraneDeclNode struct {
	baseNode
	Name       *Ident   // Brane name
	Dimensions []string // Dimension names
	// TODO(M2.4.1): Add geometry, metric, boundary conditions, etc.
}

func (b *BraneDeclNode) declNode() {}
func (b *BraneDeclNode) stringTheoryNode() {}

// StringOnBraneNode represents a string attached to branes.
// Describes which branes a string is coupled to.
// Planned syntax: `string <name> on <brane1>, <brane2>, ...`
// Status: Interface placeholder (implementation in M2.4.1)
type StringOnBraneNode struct {
	baseNode
	StringName *Ident   // String name
	BraneNames []*Ident // List of brane names
	// TODO(M2.4.1): Add attachment points, boundary conditions, etc.
}

func (s *StringOnBraneNode) stringTheoryNode() {}

// CouplingNode represents mode coupling between strings.
// Defines how vibrations in one string affect another.
// Planned syntax: `coupling <string>.<mode> { mode_pair (a, b) with strength <value> }`
// Status: Interface placeholder (implementation in M2.4.1)
type CouplingNode struct {
	baseNode
	StringA  *Ident // First string name
	ModeA    *Ident // First mode name
	StringB  *Ident // Second string name
	ModeB    *Ident // Second mode name
	Strength Expr   // Coupling strength (expression)
	// TODO(M2.4.1): Add coupling type, nonlinearity, etc.
}

func (c *CouplingNode) declNode() {}
func (c *CouplingNode) stringTheoryNode() {}

// ResonanceRuleNode represents a resonance condition and action.
// When certain frequency/amplitude conditions are met, trigger actions.
// Planned syntax: `resonance when <condition> { <action> }`
// Status: Interface placeholder (implementation in M2.4.1)
type ResonanceRuleNode struct {
	baseNode
	Condition Expr       // Boolean condition (frequency matching, amplitude threshold, etc.)
	Action    *BlockExpr // Action to execute when resonance occurs
	// TODO(M2.4.1): Add resonance strength, damping, feedback, etc.
}

func (r *ResonanceRuleNode) declNode() {}
func (r *ResonanceRuleNode) stringTheoryNode() {}

// ==================== Entanglement Nodes (M2.5) ====================

// EntangleStmt represents an entanglement statement.
// Entangles multiple variables so they share the same underlying value.
// When one changes, all others change synchronously.
// Planned syntax: `entangle x, y, z;`
// Status: Interface placeholder (implementation in M2.5.1)
type EntangleStmt struct {
	baseNode
	Variables []*Ident // List of variables to entangle
	// TODO(M2.5.1): Add entanglement strength, scope, lifetime, etc.
}

func (e *EntangleStmt) stmtNode() {}
func (e *EntangleStmt) entanglementNode() {}

// NewEntangleStmt creates a new entanglement statement.
// This is a placeholder constructor for future implementation.
func NewEntangleStmt(variables []*Ident, span common.Span) *EntangleStmt {
	return &EntangleStmt{
		baseNode:  baseNode{span},
		Variables: variables,
	}
}

// ==================== Future Extension Notes ====================

// TODO(M2.3): Musical DSL Implementation
// - Implement TemporalGrid in internal/runtime/temporal.go
// - Implement MotifRegistry for motif instantiation
// - Add parser support for score/track/bars/motif syntax
// - Add semantic validation (bar ranges, motif parameters)

// TODO(M2.4): String Theory Implementation
// - Implement StringState in internal/runtime/string_state.go
// - Implement BraneContext in internal/runtime/brane_context.go
// - Implement ResonanceGraph in internal/runtime/resonance.go
// - Add parser support for string/brane/coupling syntax
// - Add semantic validation (mode compatibility, brane dimensions)

// TODO(M2.5): Entanglement Implementation
// - Implement EntanglementGroup in internal/runtime/entanglement.go
// - Add semantic pass in internal/semantic/entangle_pass.go
// - Implement synchronization mechanism for variable updates
// - Add type checking for entangled variables (must have compatible types)

// ==================== Placeholder Constructors ====================

// The following constructors are placeholders for future implementation.
// They provide a consistent API for creating extension nodes.

// NewScoreNode creates a placeholder for a score node.
func NewScoreNode(name *Ident, tracks []*TrackNode, span common.Span) *ScoreNode {
	return &ScoreNode{
		baseNode: baseNode{span},
		Name:     name,
		Tracks:   tracks,
	}
}

// NewTrackNode creates a placeholder for a track node.
func NewTrackNode(name *Ident, bars []*BarRangeNode, span common.Span) *TrackNode {
	return &TrackNode{
		baseNode: baseNode{span},
		Name:     name,
		Bars:     bars,
	}
}

// NewBarRangeNode creates a placeholder for a bar range node.
func NewBarRangeNode(startBar, endBar int, body *BlockExpr, span common.Span) *BarRangeNode {
	return &BarRangeNode{
		baseNode: baseNode{span},
		StartBar: startBar,
		EndBar:   endBar,
		Body:     body,
	}
}

// NewMotifDeclNode creates a placeholder for a motif declaration.
func NewMotifDeclNode(name *Ident, params []*Param, body *BlockExpr, span common.Span) *MotifDeclNode {
	return &MotifDeclNode{
		baseNode: baseNode{span},
		Name:     name,
		Params:   params,
		Body:     body,
	}
}

// NewStringDeclNode creates a placeholder for a string declaration.
func NewStringDeclNode(name *Ident, modes []*Field, span common.Span) *StringDeclNode {
	return &StringDeclNode{
		baseNode: baseNode{span},
		Name:     name,
		Modes:    modes,
	}
}

// NewBraneDeclNode creates a placeholder for a brane declaration.
func NewBraneDeclNode(name *Ident, dimensions []string, span common.Span) *BraneDeclNode {
	return &BraneDeclNode{
		baseNode:   baseNode{span},
		Name:       name,
		Dimensions: dimensions,
	}
}

// NewCouplingNode creates a placeholder for a coupling node.
func NewCouplingNode(stringA, modeA, stringB, modeB *Ident, strength Expr, span common.Span) *CouplingNode {
	return &CouplingNode{
		baseNode: baseNode{span},
		StringA:  stringA,
		ModeA:    modeA,
		StringB:  stringB,
		ModeB:    modeB,
		Strength: strength,
	}
}

// NewResonanceRuleNode creates a placeholder for a resonance rule.
func NewResonanceRuleNode(condition Expr, action *BlockExpr, span common.Span) *ResonanceRuleNode {
	return &ResonanceRuleNode{
		baseNode:  baseNode{span},
		Condition: condition,
		Action:    action,
	}
}

