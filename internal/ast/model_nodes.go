package ast

import (
	"github.com/dennislee928/uad-lang/internal/common"
)

// ==================== Model Module ====================

// ModelModule represents a complete .uadmodel file
type ModelModule struct {
	baseNode
	Decls []ModelDecl
}

func (m *ModelModule) modelNode() {}

// NewModelModule creates a new model module
func NewModelModule(decls []ModelDecl, span common.Span) *ModelModule {
	return &ModelModule{
		baseNode: baseNode{span},
		Decls:    decls,
	}
}

// ==================== Model Declarations ====================

// ModelDecl represents a model-level declaration
type ModelDecl interface {
	Node
	modelDeclNode()
}

// ==================== ActionClass ====================

// ActionClassDecl represents an action_class declaration
type ActionClassDecl struct {
	baseNode
	Name       *Ident
	Fields     []*ActionField
	Methods    []*ActionMethod
}

func (a *ActionClassDecl) modelDeclNode() {}

// NewActionClassDecl creates a new action_class declaration
func NewActionClassDecl(name *Ident, fields []*ActionField, methods []*ActionMethod, span common.Span) *ActionClassDecl {
	return &ActionClassDecl{
		baseNode: baseNode{span},
		Name:     name,
		Fields:   fields,
		Methods:  methods,
	}
}

// ActionField represents a field in an action_class
type ActionField struct {
	baseNode
	Name     *Ident
	TypeExpr TypeExpr
	Default  Expr // Optional default value
}

// NewActionField creates a new action field
func NewActionField(name *Ident, typeExpr TypeExpr, defaultVal Expr, span common.Span) *ActionField {
	return &ActionField{
		baseNode: baseNode{span},
		Name:     name,
		TypeExpr: typeExpr,
		Default:  defaultVal,
	}
}

// ActionMethod represents a method in an action_class
type ActionMethod struct {
	baseNode
	Name   *Ident
	Params []*Param
	Body   *BlockExpr
}

// NewActionMethod creates a new action method
func NewActionMethod(name *Ident, params []*Param, body *BlockExpr, span common.Span) *ActionMethod {
	return &ActionMethod{
		baseNode: baseNode{span},
		Name:     name,
		Params:   params,
		Body:     body,
	}
}

// ==================== Judge ====================

// JudgeDecl represents a judge declaration
type JudgeDecl struct {
	baseNode
	Name       *Ident
	InputType  TypeExpr
	OutputType TypeExpr
	Rules      []*JudgeRule
}

func (j *JudgeDecl) modelDeclNode() {}

// NewJudgeDecl creates a new judge declaration
func NewJudgeDecl(name *Ident, inputType, outputType TypeExpr, rules []*JudgeRule, span common.Span) *JudgeDecl {
	return &JudgeDecl{
		baseNode:   baseNode{span},
		Name:       name,
		InputType:  inputType,
		OutputType: outputType,
		Rules:      rules,
	}
}

// JudgeRule represents a rule in a judge
type JudgeRule struct {
	baseNode
	Condition Expr
	Result    Expr
}

// NewJudgeRule creates a new judge rule
func NewJudgeRule(condition, result Expr, span common.Span) *JudgeRule {
	return &JudgeRule{
		baseNode:  baseNode{span},
		Condition: condition,
		Result:    result,
	}
}

// ==================== ERH Profile ====================

// ErhProfileDecl represents an erh_profile declaration
type ErhProfileDecl struct {
	baseNode
	Name           *Ident
	DatasetBinding *DatasetBinding
	ActionClass    *Ident
	Judge          *Ident
	Config         *ErhConfig
	Output         []*OutputSpec
}

func (e *ErhProfileDecl) modelDeclNode() {}

// NewErhProfileDecl creates a new erh_profile declaration
func NewErhProfileDecl(
	name *Ident,
	dataset *DatasetBinding,
	actionClass, judge *Ident,
	config *ErhConfig,
	output []*OutputSpec,
	span common.Span,
) *ErhProfileDecl {
	return &ErhProfileDecl{
		baseNode:       baseNode{span},
		Name:           name,
		DatasetBinding: dataset,
		ActionClass:    actionClass,
		Judge:          judge,
		Config:         config,
		Output:         output,
	}
}

// DatasetBinding represents a dataset binding
type DatasetBinding struct {
	baseNode
	Source Expr // String literal or expression
}

// NewDatasetBinding creates a new dataset binding
func NewDatasetBinding(source Expr, span common.Span) *DatasetBinding {
	return &DatasetBinding{
		baseNode: baseNode{span},
		Source:   source,
	}
}

// ErhConfig represents ERH configuration
type ErhConfig struct {
	baseNode
	PrimeThreshold   Expr // Float expression
	FitAlpha         Expr // String or expression for fitting method
	ImportanceMin    Expr // Optional
	ComplexityMin    Expr // Optional
}

// NewErhConfig creates a new ERH config
func NewErhConfig(primeThreshold, fitAlpha, importanceMin, complexityMin Expr, span common.Span) *ErhConfig {
	return &ErhConfig{
		baseNode:       baseNode{span},
		PrimeThreshold: primeThreshold,
		FitAlpha:       fitAlpha,
		ImportanceMin:  importanceMin,
		ComplexityMin:  complexityMin,
	}
}

// OutputSpec represents an output specification
type OutputSpec struct {
	baseNode
	Field *Ident
	Expr  Expr
}

// NewOutputSpec creates a new output spec
func NewOutputSpec(field *Ident, expr Expr, span common.Span) *OutputSpec {
	return &OutputSpec{
		baseNode: baseNode{span},
		Field:    field,
		Expr:     expr,
	}
}

// ==================== Scenario ====================

// ScenarioDecl represents a scenario declaration
type ScenarioDecl struct {
	baseNode
	Name        *Ident
	Description Expr // String literal
	Stages      []*ScenarioStage
	Assertions  []*Assertion
}

func (s *ScenarioDecl) modelDeclNode() {}

// NewScenarioDecl creates a new scenario declaration
func NewScenarioDecl(name *Ident, description Expr, stages []*ScenarioStage, assertions []*Assertion, span common.Span) *ScenarioDecl {
	return &ScenarioDecl{
		baseNode:    baseNode{span},
		Name:        name,
		Description: description,
		Stages:      stages,
		Assertions:  assertions,
	}
}

// ScenarioStage represents a stage in a scenario
type ScenarioStage struct {
	baseNode
	Name    *Ident
	Actions []Expr // Action expressions
}

// NewScenarioStage creates a new scenario stage
func NewScenarioStage(name *Ident, actions []Expr, span common.Span) *ScenarioStage {
	return &ScenarioStage{
		baseNode: baseNode{span},
		Name:     name,
		Actions:  actions,
	}
}

// Assertion represents an assertion in a scenario
type Assertion struct {
	baseNode
	Condition Expr
	Message   Expr // Optional error message
}

// NewAssertion creates a new assertion
func NewAssertion(condition, message Expr, span common.Span) *Assertion {
	return &Assertion{
		baseNode:  baseNode{span},
		Condition: condition,
		Message:   message,
	}
}

// ==================== Cognitive SIEM ====================

// CognitiveSiemDecl represents a cognitive_siem declaration
type CognitiveSiemDecl struct {
	baseNode
	Name       *Ident
	Agent      *AgentConfig
	Rules      []*SiemRule
	Response   []*ResponseAction
}

func (c *CognitiveSiemDecl) modelDeclNode() {}

// NewCognitiveSiemDecl creates a new cognitive_siem declaration
func NewCognitiveSiemDecl(name *Ident, agent *AgentConfig, rules []*SiemRule, response []*ResponseAction, span common.Span) *CognitiveSiemDecl {
	return &CognitiveSiemDecl{
		baseNode: baseNode{span},
		Name:     name,
		Agent:    agent,
		Rules:    rules,
		Response: response,
	}
}

// AgentConfig represents agent configuration
type AgentConfig struct {
	baseNode
	Type       *Ident
	Parameters map[string]Expr
}

// NewAgentConfig creates a new agent config
func NewAgentConfig(typ *Ident, params map[string]Expr, span common.Span) *AgentConfig {
	return &AgentConfig{
		baseNode:   baseNode{span},
		Type:       typ,
		Parameters: params,
	}
}

// SiemRule represents a SIEM rule
type SiemRule struct {
	baseNode
	Name      *Ident
	Condition Expr
	Severity  Expr // Int or String expression
	Action    Expr
}

// NewSiemRule creates a new SIEM rule
func NewSiemRule(name *Ident, condition, severity, action Expr, span common.Span) *SiemRule {
	return &SiemRule{
		baseNode:  baseNode{span},
		Name:      name,
		Condition: condition,
		Severity:  severity,
		Action:    action,
	}
}

// ResponseAction represents a response action
type ResponseAction struct {
	baseNode
	Type   *Ident
	Params []Expr
}

// NewResponseAction creates a new response action
func NewResponseAction(typ *Ident, params []Expr, span common.Span) *ResponseAction {
	return &ResponseAction{
		baseNode: baseNode{span},
		Type:     typ,
		Params:   params,
	}
}

// ==================== Dataset Reference ====================

// DatasetRef represents a reference to a dataset field
type DatasetRef struct {
	baseNode
	Field *Ident
}

func (d *DatasetRef) exprNode() {}

// NewDatasetRef creates a new dataset reference
func NewDatasetRef(field *Ident, span common.Span) *DatasetRef {
	return &DatasetRef{
		baseNode: baseNode{span},
		Field:    field,
	}
}

// ==================== Model-Specific Expressions ====================

// CaseExpr represents a case expression (for iterating over dataset)
type CaseExpr struct {
	baseNode
	Variable *Ident
	Source   Expr // Dataset reference
	Body     *BlockExpr
}

func (c *CaseExpr) exprNode() {}

// NewCaseExpr creates a new case expression
func NewCaseExpr(variable *Ident, source Expr, body *BlockExpr, span common.Span) *CaseExpr {
	return &CaseExpr{
		baseNode: baseNode{span},
		Variable: variable,
		Source:   source,
		Body:     body,
	}
}

// AggregateExpr represents an aggregation expression
type AggregateExpr struct {
	baseNode
	Func   *Ident // count, sum, avg, etc.
	Source Expr
	Filter Expr // Optional filter expression
}

func (a *AggregateExpr) exprNode() {}

// NewAggregateExpr creates a new aggregate expression
func NewAggregateExpr(fn *Ident, source, filter Expr, span common.Span) *AggregateExpr {
	return &AggregateExpr{
		baseNode: baseNode{span},
		Func:     fn,
		Source:   source,
		Filter:   filter,
	}
}


