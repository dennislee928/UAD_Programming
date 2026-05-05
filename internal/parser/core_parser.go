package parser

import (
	"fmt"
	"strconv"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/common"
	"github.com/dennislee928/uad-lang/internal/lexer"
)

// Parser parses tokens into an AST
type Parser struct {
	tokens  []lexer.Token
	pos     int
	errors  *common.ErrorList
	file    string
}

// New creates a new parser
func New(tokens []lexer.Token, file string) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		errors: common.NewErrorList(),
		file:   file,
	}
}

// ParseModule parses a complete module
func (p *Parser) ParseModule() (*ast.Module, error) {
	decls := []ast.Decl{}
	
	for !p.isAtEnd() {
		// Skip comments
		if p.match(lexer.TokenComment) {
			continue
		}
		
		decl, err := p.parseDecl()
		if err != nil {
			p.errors.Add(err.(*common.Error))
			// Try to recover by skipping to next declaration
			p.synchronize()
			continue
		}
		
		if decl != nil {
			decls = append(decls, decl)
		}
	}
	
	if p.errors.HasErrors() {
		return nil, p.errors
	}
	
	span := common.Span{File: p.file}
	if len(decls) > 0 {
		span.Start = decls[0].Span().Start
		span.End = decls[len(decls)-1].Span().End
	}
	
	return ast.NewModule(decls, span), nil
}

// Errors returns the error list
func (p *Parser) Errors() *common.ErrorList {
	return p.errors
}

// ==================== Declarations ====================

func (p *Parser) parseDecl() (ast.Decl, error) {
	switch {
	case p.check(lexer.TokenFn):
		return p.parseFnDecl()
	case p.check(lexer.TokenStruct):
		return p.parseStructDecl()
	case p.check(lexer.TokenEnum):
		return p.parseEnumDecl()
	case p.check(lexer.TokenTypeDecl):
		return p.parseTypeAlias()
	case p.check(lexer.TokenImport):
		return p.parseImportDecl()
	default:
		// Try extension declarations (M2.3-M2.5)
		if extDecl, err := p.ParseDeclExtension(); extDecl != nil || err != nil {
			return extDecl, err
		}
		return nil, p.error("expected declaration")
	}
}

func (p *Parser) parseFnDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenFn, "expected 'fn'")
	
	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected function name")
	}
	
	p.consume(lexer.TokenLParen, "expected '('")
	
	params, err := p.parseParamList()
	if err != nil {
		return nil, err
	}
	
	p.consume(lexer.TokenRParen, "expected ')'")
	
	var returnType ast.TypeExpr
	if p.match(lexer.TokenArrow) {
		returnType, err = p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
	}
	
	body, err := p.parseBlockExpr()
	if err != nil {
		return nil, err
	}
	
	end := body.Span().End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewFnDecl(name, params, returnType, body, span), nil
}

func (p *Parser) parseParamList() ([]*ast.Param, error) {
	params := []*ast.Param{}
	
	if p.check(lexer.TokenRParen) {
		return params, nil
	}
	
	for {
		start := p.current().Span.Start
		name := p.parseIdent()
		if name == nil {
			return nil, p.error("expected parameter name")
		}
		
		p.consume(lexer.TokenColon, "expected ':'")
		
		typeExpr, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		params = append(params, ast.NewParam(name, typeExpr, span))
		
		if !p.match(lexer.TokenComma) {
			break
		}
	}
	
	return params, nil
}

func (p *Parser) parseStructDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenStruct, "expected 'struct'")
	
	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected struct name")
	}
	
	p.consume(lexer.TokenLBrace, "expected '{'")
	
	fields, err := p.parseFieldList()
	if err != nil {
		return nil, err
	}
	
	p.consume(lexer.TokenRBrace, "expected '}'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewStructDecl(name, fields, span), nil
}

func (p *Parser) parseFieldList() ([]*ast.Field, error) {
	fields := []*ast.Field{}
	
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		start := p.current().Span.Start
		// Allow keywords as field names (e.g., "type" as a field name)
		var name *ast.Ident
		if p.check(lexer.TokenIdent) {
			name = p.parseIdent()
		} else if lexer.IsKeyword(p.current().Type) {
			tok := p.current()
			p.advance()
			name = ast.NewIdent(tok.Lexeme, tok.Span)
		}
		if name == nil {
			return nil, p.error("expected field name")
		}
		
		p.consume(lexer.TokenColon, "expected ':'")
		
		typeExpr, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		fields = append(fields, ast.NewField(name, typeExpr, span))
		
		// Optional comma or newline
		p.match(lexer.TokenComma)
	}
	
	return fields, nil
}

func (p *Parser) parseEnumDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenEnum, "expected 'enum'")
	
	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected enum name")
	}
	
	p.consume(lexer.TokenLBrace, "expected '{'")
	
	variants, err := p.parseVariantList()
	if err != nil {
		return nil, err
	}
	
	p.consume(lexer.TokenRBrace, "expected '}'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewEnumDecl(name, variants, span), nil
}

func (p *Parser) parseVariantList() ([]*ast.Variant, error) {
	variants := []*ast.Variant{}
	
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		start := p.current().Span.Start
		name := p.parseIdent()
		if name == nil {
			return nil, p.error("expected variant name")
		}
		
		var types []ast.TypeExpr
		if p.match(lexer.TokenLParen) {
			for !p.check(lexer.TokenRParen) && !p.isAtEnd() {
				typeExpr, err := p.parseTypeExpr()
				if err != nil {
					return nil, err
				}
				types = append(types, typeExpr)
				
				if !p.match(lexer.TokenComma) {
					break
				}
			}
			p.consume(lexer.TokenRParen, "expected ')'")
		}
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		variants = append(variants, ast.NewVariant(name, types, span))
		
		p.match(lexer.TokenComma)
	}
	
	return variants, nil
}

func (p *Parser) parseTypeAlias() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenTypeDecl, "expected 'type'")
	
	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected type name")
	}
	
	p.consume(lexer.TokenAssign, "expected '='")
	
	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}
	
	p.consume(lexer.TokenSemicolon, "expected ';'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewTypeAlias(name, typeExpr, span), nil
}

func (p *Parser) parseImportDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenImport, "expected 'import'")
	
	if !p.check(lexer.TokenString) {
		return nil, p.error("expected string literal")
	}
	
	pathToken := p.advance()
	path := pathToken.Lexeme
	// Remove quotes
	if len(path) >= 2 {
		path = path[1 : len(path)-1]
	}
	
	p.consume(lexer.TokenSemicolon, "expected ';'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewImportDecl(path, span), nil
}

// ==================== Type Expressions ====================

func (p *Parser) parseTypeExpr() (ast.TypeExpr, error) {
	start := p.current().Span.Start
	
	// Array type
	if p.match(lexer.TokenLBracket) {
		elemType, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		p.consume(lexer.TokenRBracket, "expected ']'")
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		return ast.NewArrayType(elemType, span), nil
	}
	
	// Map type
	if p.check(lexer.TokenIdent) && p.current().Lexeme == "Map" {
		p.advance()
		p.consume(lexer.TokenLBracket, "expected '['")
		
		keyType, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		
		p.consume(lexer.TokenComma, "expected ','")
		
		valueType, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		
		p.consume(lexer.TokenRBracket, "expected ']'")
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		return ast.NewMapType(keyType, valueType, span), nil
	}
	
	// Function type
	if p.match(lexer.TokenFn) {
		p.consume(lexer.TokenLParen, "expected '('")
		
		paramTypes := []ast.TypeExpr{}
		if !p.check(lexer.TokenRParen) {
			for {
				paramType, err := p.parseTypeExpr()
				if err != nil {
					return nil, err
				}
				paramTypes = append(paramTypes, paramType)
				
				if !p.match(lexer.TokenComma) {
					break
				}
			}
		}
		
		p.consume(lexer.TokenRParen, "expected ')'")
		p.consume(lexer.TokenArrow, "expected '->'")
		
		returnType, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		return ast.NewFunctionType(paramTypes, returnType, span), nil
	}
	
	// Named type
	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected type name")
	}
	
	return ast.NewNamedType(name, name.Span()), nil
}

// ==================== Statements ====================

func (p *Parser) parseStmt() (ast.Stmt, error) {
	switch {
	case p.check(lexer.TokenLet):
		return p.parseLetStmt()
	case p.check(lexer.TokenReturn):
		return p.parseReturnStmt()
	case p.check(lexer.TokenWhile):
		return p.parseWhileStmt()
	case p.check(lexer.TokenFor):
		return p.parseForStmt()
	case p.check(lexer.TokenBreak):
		return p.parseBreakStmt()
	case p.check(lexer.TokenContinue):
		return p.parseContinueStmt()
	default:
		// Try extension statements (M2.3-M2.5)
		if extStmt, err := p.ParseStmtExtension(); extStmt != nil || err != nil {
			return extStmt, err
		}
		// Try to parse as expression statement or assignment
		return p.parseExprOrAssignStmt()
	}
}

func (p *Parser) parseLetStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenLet, "expected 'let'")
	
	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected variable name")
	}
	
	var typeExpr ast.TypeExpr
	if p.match(lexer.TokenColon) {
		var err error
		typeExpr, err = p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
	}
	
	p.consume(lexer.TokenAssign, "expected '='")
	
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	
	p.consume(lexer.TokenSemicolon, "expected ';'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewLetStmt(name, typeExpr, value, span), nil
}

func (p *Parser) parseReturnStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenReturn, "expected 'return'")
	
	var value ast.Expr
	if !p.check(lexer.TokenSemicolon) {
		var err error
		value, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	
	p.consume(lexer.TokenSemicolon, "expected ';'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewReturnStmt(value, span), nil
}

func (p *Parser) parseWhileStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenWhile, "expected 'while'")
	
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	
	body, err := p.parseBlockExpr()
	if err != nil {
		return nil, err
	}
	
	end := body.Span().End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewWhileStmt(cond, body, span), nil
}

func (p *Parser) parseForStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenFor, "expected 'for'")
	
	varName := p.parseIdent()
	if varName == nil {
		return nil, p.error("expected loop variable")
	}
	
	p.consume(lexer.TokenIn, "expected 'in'")
	
	iter, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	
	body, err := p.parseBlockExpr()
	if err != nil {
		return nil, err
	}
	
	end := body.Span().End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewForStmt(varName, iter, body, span), nil
}

func (p *Parser) parseBreakStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenBreak, "expected 'break'")
	p.consume(lexer.TokenSemicolon, "expected ';'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewBreakStmt(span), nil
}

func (p *Parser) parseContinueStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenContinue, "expected 'continue'")
	p.consume(lexer.TokenSemicolon, "expected ';'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewContinueStmt(span), nil
}

func (p *Parser) parseExprOrAssignStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	
	// Check for assignment
	if p.match(lexer.TokenAssign) {
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		
		p.consume(lexer.TokenSemicolon, "expected ';'")
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		
		return ast.NewAssignStmt(expr, value, span), nil
	}
	
	p.consume(lexer.TokenSemicolon, "expected ';'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewExprStmt(expr, span), nil
}

// ==================== Expressions (Pratt Parser) ====================

// Operator precedence levels
const (
	PrecLowest      = iota
	PrecOr          // ||
	PrecAnd         // &&
	PrecEquality    // == !=
	PrecComparison  // < > <= >=
	PrecSum         // + -
	PrecProduct     // * / %
	PrecUnary       // -x !x
	PrecCall        // f() a.b a[i]
)

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseExprWithPrecedence(PrecLowest)
}

func (p *Parser) parseExprWithPrecedence(precedence int) (ast.Expr, error) {
	// Parse prefix expression
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}
	
	// Parse infix expressions
	for !p.isAtEnd() && precedence < p.currentPrecedence() {
		left, err = p.parseInfixExpr(left)
		if err != nil {
			return nil, err
		}
	}
	
	return left, nil
}

func (p *Parser) parsePrefixExpr() (ast.Expr, error) {
	start := p.current().Span.Start
	
	// Unary operators
	if p.match(lexer.TokenMinus) {
		expr, err := p.parseExprWithPrecedence(PrecUnary)
		if err != nil {
			return nil, err
		}
		end := expr.Span().End
		span := common.NewSpan(start, end, p.file)
		return ast.NewUnaryExpr(ast.OpNeg, expr, span), nil
	}
	
	if p.match(lexer.TokenNot) {
		expr, err := p.parseExprWithPrecedence(PrecUnary)
		if err != nil {
			return nil, err
		}
		end := expr.Span().End
		span := common.NewSpan(start, end, p.file)
		return ast.NewUnaryExpr(ast.OpNot, expr, span), nil
	}
	
	// Primary expressions
	return p.parsePrimaryExpr()
}

func (p *Parser) parseInfixExpr(left ast.Expr) (ast.Expr, error) {
	start := left.Span().Start
	
	// Binary operators
	op := p.current()
	var binaryOp ast.BinaryOp
	
	switch op.Type {
	case lexer.TokenPlus:
		binaryOp = ast.OpAdd
	case lexer.TokenMinus:
		binaryOp = ast.OpSub
	case lexer.TokenStar:
		binaryOp = ast.OpMul
	case lexer.TokenSlash:
		binaryOp = ast.OpDiv
	case lexer.TokenMod:
		binaryOp = ast.OpMod
	case lexer.TokenEq:
		binaryOp = ast.OpEq
	case lexer.TokenNeq:
		binaryOp = ast.OpNeq
	case lexer.TokenLt:
		binaryOp = ast.OpLt
	case lexer.TokenGt:
		binaryOp = ast.OpGt
	case lexer.TokenLe:
		binaryOp = ast.OpLe
	case lexer.TokenGe:
		binaryOp = ast.OpGe
	case lexer.TokenAnd:
		binaryOp = ast.OpAnd
	case lexer.TokenOr:
		binaryOp = ast.OpOr
	default:
		// Postfix operators
		return p.parsePostfixExpr(left)
	}
	
	p.advance()
	precedence := p.tokenPrecedence(op.Type)
	
	right, err := p.parseExprWithPrecedence(precedence)
	if err != nil {
		return nil, err
	}
	
	end := right.Span().End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewBinaryExpr(binaryOp, left, right, span), nil
}

func (p *Parser) parsePostfixExpr(left ast.Expr) (ast.Expr, error) {
	start := left.Span().Start
	
	// Call expression
	if p.match(lexer.TokenLParen) {
		args := []ast.Expr{}
		
		if !p.check(lexer.TokenRParen) {
			for {
				arg, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				
				if !p.match(lexer.TokenComma) {
					break
				}
			}
		}
		
		p.consume(lexer.TokenRParen, "expected ')'")
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		
		return ast.NewCallExpr(left, args, span), nil
	}
	
	// Field access
	if p.match(lexer.TokenDot) {
		field := p.parseIdent()
		if field == nil {
			return nil, p.error("expected field name")
		}
		
		end := field.Span().End
		span := common.NewSpan(start, end, p.file)
		
		return ast.NewFieldAccess(left, field, span), nil
	}
	
	// Index expression
	if p.match(lexer.TokenLBracket) {
		index, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		
		p.consume(lexer.TokenRBracket, "expected ']'")
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		
		return ast.NewIndexExpr(left, index, span), nil
	}
	
	return left, nil
}

func (p *Parser) parsePrimaryExpr() (ast.Expr, error) {
	start := p.current().Span.Start
	
	// Literals
	if p.check(lexer.TokenInt) || p.check(lexer.TokenFloat) ||
		p.check(lexer.TokenString) || p.check(lexer.TokenTrue) ||
		p.check(lexer.TokenFalse) || p.check(lexer.TokenNil) ||
		p.check(lexer.TokenDuration) {
		return p.parseLiteral()
	}
	
	// Parenthesized expression
	if p.match(lexer.TokenLParen) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		p.consume(lexer.TokenRParen, "expected ')'")
		
		end := p.previous().Span.End
		span := common.NewSpan(start, end, p.file)
		return ast.NewParenExpr(expr, span), nil
	}
	
	// If expression
	if p.check(lexer.TokenIf) {
		return p.parseIfExpr()
	}
	
	// Match expression
	if p.check(lexer.TokenMatch) {
		return p.parseMatchExpr()
	}
	
	// Block expression
	if p.check(lexer.TokenLBrace) {
		return p.parseBlockExpr()
	}
	
	// Array literal
	if p.check(lexer.TokenLBracket) {
		return p.parseArrayLiteral()
	}
	
	// Identifier (could be variable, function call, or struct literal)
	if p.check(lexer.TokenIdent) {
		return p.parseIdentOrStructLiteral()
	}
	
	return nil, p.error(fmt.Sprintf("unexpected token: %s", p.current().Type))
}

func (p *Parser) parseLiteral() (ast.Expr, error) {
	tok := p.advance()
	
	var kind ast.LiteralKind
	switch tok.Type {
	case lexer.TokenInt:
		kind = ast.LitInt
	case lexer.TokenFloat:
		kind = ast.LitFloat
	case lexer.TokenString:
		kind = ast.LitString
	case lexer.TokenTrue, lexer.TokenFalse:
		kind = ast.LitBool
	case lexer.TokenDuration:
		kind = ast.LitDuration
	case lexer.TokenNil:
		kind = ast.LitNil
	}
	
	return ast.NewLiteral(kind, tok.Lexeme, tok.Span), nil
}

func (p *Parser) parseIdent() *ast.Ident {
	if !p.check(lexer.TokenIdent) {
		return nil
	}
	
	tok := p.advance()
	return ast.NewIdent(tok.Lexeme, tok.Span)
}

func (p *Parser) parseIdentOrStructLiteral() (ast.Expr, error) {
	ident := p.parseIdent()
	if ident == nil {
		return nil, p.error("expected identifier")
	}
	
	// Check for struct literal
	if p.check(lexer.TokenLBrace) {
		return p.parseStructLiteralWithName(ident)
	}
	
	return ident, nil
}

func (p *Parser) parseStructLiteralWithName(name *ast.Ident) (ast.Expr, error) {
	start := name.Span().Start
	p.consume(lexer.TokenLBrace, "expected '{'")
	
	fields := []*ast.FieldInit{}
	
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		fieldStart := p.current().Span.Start
		// In struct literals, keywords can be used as field names
		var fieldName *ast.Ident
		if p.check(lexer.TokenIdent) {
			fieldName = p.parseIdent()
		} else {
			// Try to parse keyword as identifier (e.g., "type" as field name)
			tok := p.current()
			if lexer.IsKeyword(tok.Type) {
				p.advance()
				fieldName = ast.NewIdent(tok.Lexeme, tok.Span)
			}
		}
		if fieldName == nil {
			return nil, p.error("expected field name")
		}
		
		p.consume(lexer.TokenColon, "expected ':'")
		
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		
		fieldEnd := value.Span().End
		fieldSpan := common.NewSpan(fieldStart, fieldEnd, p.file)
		fields = append(fields, ast.NewFieldInit(fieldName, value, fieldSpan))
		
		if !p.match(lexer.TokenComma) {
			break
		}
	}
	
	p.consume(lexer.TokenRBrace, "expected '}'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewStructLiteral(name, fields, span), nil
}

func (p *Parser) parseArrayLiteral() (ast.Expr, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenLBracket, "expected '['")
	
	elements := []ast.Expr{}
	
	if !p.check(lexer.TokenRBracket) {
		for {
			elem, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
			
			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}
	
	p.consume(lexer.TokenRBracket, "expected ']'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewArrayLiteral(elements, span), nil
}

func (p *Parser) parseIfExpr() (ast.Expr, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenIf, "expected 'if'")
	
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	
	thenBlock, err := p.parseBlockExpr()
	if err != nil {
		return nil, err
	}
	
	var elseExpr ast.Expr
	if p.match(lexer.TokenElse) {
		if p.check(lexer.TokenIf) {
			// else if
			elseExpr, err = p.parseIfExpr()
		} else {
			// else block
			elseExpr, err = p.parseBlockExpr()
		}
		if err != nil {
			return nil, err
		}
	}
	
	end := thenBlock.Span().End
	if elseExpr != nil {
		end = elseExpr.Span().End
	}
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewIfExpr(cond, thenBlock, elseExpr, span), nil
}

func (p *Parser) parseMatchExpr() (ast.Expr, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenMatch, "expected 'match'")
	
	// Parse the expression being matched
	// Special handling: parse identifier but don't allow struct literals
	// because the { is for the match arms
	var expr ast.Expr
	var err error
	
	if p.check(lexer.TokenIdent) {
		// Just parse the identifier, don't check for struct literal
		ident := p.parseIdent()
		expr = ident
	} else {
		// For other expressions (like field access, etc.), use normal parsing
		// but this should stop before consuming {
		expr, err = p.parsePrefixExpr()
		if err != nil {
			return nil, err
		}
		
		// Parse postfix operations (field access, indexing, etc.)
		// but not struct literals
		for p.check(lexer.TokenDot) || p.check(lexer.TokenLBracket) {
			if p.match(lexer.TokenDot) {
				field := p.parseIdent()
				if field == nil {
					return nil, p.error("expected field name")
				}
				fieldEnd := field.Span().End
				span := common.NewSpan(expr.Span().Start, fieldEnd, p.file)
				expr = ast.NewFieldAccess(expr, field, span)
			} else if p.match(lexer.TokenLBracket) {
				index, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				p.consume(lexer.TokenRBracket, "expected ']'")
				end := p.previous().Span.End
				span := common.NewSpan(expr.Span().Start, end, p.file)
				expr = ast.NewIndexExpr(expr, index, span)
			}
		}
	}
	
	p.consume(lexer.TokenLBrace, "expected '{'")
	
	arms := []*ast.MatchArm{}
	
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		armStart := p.current().Span.Start
		
		pattern, err := p.parsePattern()
		if err != nil {
			return nil, err
		}
		
		if err := p.consume(lexer.TokenFatArrow, "expected '=>'"); err != nil {
			return nil, err
		}
		
		armExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		
		armEnd := armExpr.Span().End
		armSpan := common.NewSpan(armStart, armEnd, p.file)
		arms = append(arms, ast.NewMatchArm(pattern, armExpr, armSpan))
		
		// Optional trailing comma
		if !p.match(lexer.TokenComma) {
			break
		}
		
		// Allow trailing comma before closing brace
		if p.check(lexer.TokenRBrace) {
			break
		}
	}
	
	p.consume(lexer.TokenRBrace, "expected '}'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewMatchExpr(expr, arms, span), nil
}

func (p *Parser) parseBlockExpr() (*ast.BlockExpr, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenLBrace, "expected '{'")
	
	stmts := []ast.Stmt{}
	var finalExpr ast.Expr
	
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		// Skip comments
		if p.match(lexer.TokenComment) {
			continue
		}
		
		// Check if this is a statement keyword
		if p.isStmtStart() {
			stmt, err := p.parseStmt()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, stmt)
			continue
		}
		
		// Otherwise, try to parse as expression/assignment
		exprStart := p.current().Span.Start
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		
		// Check for assignment
		if p.match(lexer.TokenAssign) {
			value, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			
			p.consume(lexer.TokenSemicolon, "expected ';'")
			
			end := p.previous().Span.End
			span := common.NewSpan(exprStart, end, p.file)
			stmts = append(stmts, ast.NewAssignStmt(expr, value, span))
			continue
		}
		
		// Check if followed by semicolon (statement) or closing brace (final expression)
		if p.match(lexer.TokenSemicolon) {
			span := common.NewSpan(expr.Span().Start, p.previous().Span.End, p.file)
			stmts = append(stmts, ast.NewExprStmt(expr, span))
		} else if p.check(lexer.TokenRBrace) {
			finalExpr = expr
			break
		} else {
			return nil, p.error("expected ';' or '}'")
		}
	}
	
	p.consume(lexer.TokenRBrace, "expected '}'")
	
	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	
	return ast.NewBlockExpr(stmts, finalExpr, span), nil
}

// ==================== Patterns ====================

func (p *Parser) parsePattern() (ast.Pattern, error) {
	start := p.current().Span.Start
	
	// Wildcard pattern
	if p.check(lexer.TokenIdent) && p.current().Lexeme == "_" {
		p.advance()
		span := common.NewSpan(start, p.previous().Span.End, p.file)
		return ast.NewWildcardPattern(span), nil
	}
	
	// Literal pattern
	if p.check(lexer.TokenInt) || p.check(lexer.TokenFloat) ||
		p.check(lexer.TokenString) || p.check(lexer.TokenTrue) ||
		p.check(lexer.TokenFalse) {
		lit, err := p.parseLiteral()
		if err != nil {
			return nil, err
		}
		return ast.NewLiteralPattern(lit.(*ast.Literal), lit.Span()), nil
	}
	
	// Identifier or enum/struct pattern
	if p.check(lexer.TokenIdent) {
		ident := p.parseIdent()
		
		// Struct pattern
		if p.match(lexer.TokenLBrace) {
			fields := []*ast.FieldPattern{}
			
			for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
				fieldStart := p.current().Span.Start
				fieldName := p.parseIdent()
				if fieldName == nil {
					return nil, p.error("expected field name")
				}
				
				var fieldPattern ast.Pattern
				if p.match(lexer.TokenColon) {
					var err error
					fieldPattern, err = p.parsePattern()
					if err != nil {
						return nil, err
					}
				}
				
				fieldEnd := p.previous().Span.End
				fieldSpan := common.NewSpan(fieldStart, fieldEnd, p.file)
				fields = append(fields, ast.NewFieldPattern(fieldName, fieldPattern, fieldSpan))
				
				if !p.match(lexer.TokenComma) {
					break
				}
			}
			
			p.consume(lexer.TokenRBrace, "expected '}'")
			
			end := p.previous().Span.End
			span := common.NewSpan(start, end, p.file)
			return ast.NewStructPattern(ident, fields, span), nil
		}
		
		// Enum pattern
		if p.match(lexer.TokenLParen) {
			patterns := []ast.Pattern{}
			
			// Parse patterns inside parentheses
			if !p.check(lexer.TokenRParen) {
				for {
					pattern, err := p.parsePattern()
					if err != nil {
						return nil, err
					}
					patterns = append(patterns, pattern)
					
					if p.check(lexer.TokenRParen) {
						break
					}
					
					if !p.match(lexer.TokenComma) {
						return nil, p.error("expected ',' or ')'")
					}
				}
			}
			
			p.consume(lexer.TokenRParen, "expected ')'")
			
			end := p.previous().Span.End
			span := common.NewSpan(start, end, p.file)
			return ast.NewEnumPattern(ident, patterns, span), nil
		}
		
		// Simple identifier pattern
		return ast.NewIdentPattern(ident, ident.Span()), nil
	}
	
	return nil, p.error("expected pattern")
}

// ==================== Helper Functions ====================

func (p *Parser) currentPrecedence() int {
	if p.isAtEnd() {
		return PrecLowest
	}
	
	return p.tokenPrecedence(p.current().Type)
}

func (p *Parser) tokenPrecedence(typ lexer.TokenType) int {
	switch typ {
	case lexer.TokenOr:
		return PrecOr
	case lexer.TokenAnd:
		return PrecAnd
	case lexer.TokenEq, lexer.TokenNeq:
		return PrecEquality
	case lexer.TokenLt, lexer.TokenGt, lexer.TokenLe, lexer.TokenGe:
		return PrecComparison
	case lexer.TokenPlus, lexer.TokenMinus:
		return PrecSum
	case lexer.TokenStar, lexer.TokenSlash, lexer.TokenMod:
		return PrecProduct
	case lexer.TokenLParen, lexer.TokenDot, lexer.TokenLBracket:
		return PrecCall
	default:
		return PrecLowest
	}
}

func (p *Parser) isStmtStart() bool {
	return p.check(lexer.TokenLet) ||
		p.check(lexer.TokenReturn) ||
		p.check(lexer.TokenWhile) ||
		p.check(lexer.TokenFor) ||
		p.check(lexer.TokenBreak) ||
		p.check(lexer.TokenContinue) ||
		p.check(lexer.TokenEmit) ||
		p.check(lexer.TokenEntangle) ||
		p.check(lexer.TokenUse) ||
		p.check(lexer.TokenBars)
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.pos]
}

func (p *Parser) previous() lexer.Token {
	if p.pos == 0 {
		return p.tokens[0]
	}
	return p.tokens[p.pos-1]
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.pos++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	if p.pos >= len(p.tokens) {
		return true
	}
	return p.tokens[p.pos].Type == lexer.TokenEOF
}

func (p *Parser) check(typ lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.current().Type == typ
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(typ lexer.TokenType, message string) error {
	if p.check(typ) {
		p.advance()
		return nil
	}
	return p.error(message)
}

func (p *Parser) error(message string) error {
	tok := p.current()
	return common.SyntaxError(message, tok.Span)
}

func (p *Parser) synchronize() {
	p.advance()
	
	for !p.isAtEnd() {
		if p.previous().Type == lexer.TokenSemicolon {
			return
		}
		
		switch p.current().Type {
		case lexer.TokenFn, lexer.TokenStruct, lexer.TokenEnum,
			lexer.TokenLet, lexer.TokenReturn, lexer.TokenIf,
			lexer.TokenWhile, lexer.TokenFor:
			return
		}
		
		p.advance()
	}
}

// ParseInt parses an integer literal
func ParseInt(s string) (int64, error) {
	if len(s) >= 2 {
		if s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			return strconv.ParseInt(s[2:], 16, 64)
		}
		if s[0] == '0' && (s[1] == 'b' || s[1] == 'B') {
			return strconv.ParseInt(s[2:], 2, 64)
		}
	}
	return strconv.ParseInt(s, 10, 64)
}

// ParseFloat parses a float literal
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

