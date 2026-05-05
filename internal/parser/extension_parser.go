package parser

import (
	"strconv"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/common"
	"github.com/dennislee928/uad-lang/internal/lexer"
)

// extension_parser.go handles parsing for M2.3 (Musical DSL), M2.4 (String Theory), and M2.5 (Entanglement).

// ==================== Musical DSL Parsing (M2.3) ====================

// parseScoreDecl parses a score declaration.
// Syntax: score <name> { <tracks>... }
func (p *Parser) parseScoreDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenScore, "expected 'score'")

	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected score name")
	}

	p.consume(lexer.TokenLBrace, "expected '{'")

	tracks := []*ast.TrackNode{}
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		if p.match(lexer.TokenComment) {
			continue
		}
		
		if p.check(lexer.TokenTrack) {
			track, err := p.parseTrackNode()
			if err != nil {
				return nil, err
			}
			tracks = append(tracks, track)
		} else if p.check(lexer.TokenIdent) {
			// Skip metadata fields like "tempo: 120,"
			p.advance() // field name
			if p.match(lexer.TokenColon) {
				// Skip the value expression
				_, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				p.match(lexer.TokenComma) // optional comma
			}
		} else {
			return nil, p.error("expected 'track' or metadata field in score body")
		}
	}

	end := p.current().Span.End
	p.consume(lexer.TokenRBrace, "expected '}'")

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewScoreNode(name, tracks, span), nil
}

// parseTrackNode parses a track within a score.
// Syntax: track <name> { <bars>... }
func (p *Parser) parseTrackNode() (*ast.TrackNode, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenTrack, "expected 'track'")

	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected track name")
	}

	p.consume(lexer.TokenLBrace, "expected '{'")

	bars := []*ast.BarRangeNode{}
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		if p.match(lexer.TokenComment) {
			continue
		}
		
		if p.check(lexer.TokenBars) {
			barRange, err := p.parseBarRangeNode()
			if err != nil {
				return nil, err
			}
			bars = append(bars, barRange)
		} else {
			return nil, p.error("expected 'bars' in track body")
		}
	}

	end := p.current().Span.End
	p.consume(lexer.TokenRBrace, "expected '}'")

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewTrackNode(name, bars, span), nil
}

// parseBarRangeNode parses a bar range.
// Syntax: bars <start>..<end> { <statements>... }
func (p *Parser) parseBarRangeNode() (*ast.BarRangeNode, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenBars, "expected 'bars'")

	if !p.check(lexer.TokenInt) {
		return nil, p.error("expected integer for start bar")
	}
	startBarToken := p.advance()
	startBar, err := strconv.Atoi(startBarToken.Lexeme)
	if err != nil {
		return nil, p.error("invalid start bar number")
	}

	p.consume(lexer.TokenDotDot, "expected '..' for bar range")

	if !p.check(lexer.TokenInt) {
		return nil, p.error("expected integer for end bar")
	}
	endBarToken := p.advance()
	endBar, err := strconv.Atoi(endBarToken.Lexeme)
	if err != nil {
		return nil, p.error("invalid end bar number")
	}

	// parseBlockExpr will consume the '{'
	body, err := p.parseBlockExpr()
	if err != nil {
		return nil, err
	}

	end := p.current().Span.End

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewBarRangeNode(startBar, endBar, body, span), nil
}

// parseMotifDecl parses a motif declaration.
// Syntax: motif <name>(<params>) { <body> }
func (p *Parser) parseMotifDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenMotif, "expected 'motif'")

	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected motif name")
	}

	// Parse optional parameters
	params := []*ast.Param{}
	if p.match(lexer.TokenLParen) {
		if !p.check(lexer.TokenRParen) {
			for {
				param, err := p.parseParam()
				if err != nil {
					return nil, err
				}
				params = append(params, param)

				if !p.match(lexer.TokenComma) {
					break
				}
			}
		}
		p.consume(lexer.TokenRParen, "expected ')' after parameters")
	}

	// parseBlockExpr will consume the '{'
	body, err := p.parseBlockExpr()
	if err != nil {
		return nil, err
	}

	end := p.current().Span.End

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewMotifDeclNode(name, params, body, span), nil
}

// ==================== String Theory DSL Parsing (M2.4) ====================

// parseStringDecl parses a string declaration.
// Syntax: string <name> { modes { <field>: <type>, ... } }
func (p *Parser) parseStringDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenStringKw, "expected 'string'")

	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected string name")
	}

	p.consume(lexer.TokenLBrace, "expected '{'")

	// Expect "modes { ... }"
	if !p.match(lexer.TokenModes) {
		return nil, p.error("expected 'modes' in string declaration")
	}

	p.consume(lexer.TokenLBrace, "expected '{' after 'modes'")

	modes := []*ast.Field{}
	for !p.check(lexer.TokenRBrace) && !p.isAtEnd() {
		if p.match(lexer.TokenComment) {
			continue
		}

		field, err := p.parseField()
		if err != nil {
			return nil, err
		}
		modes = append(modes, field)

		// Optional comma
		p.match(lexer.TokenComma)
	}

	p.consume(lexer.TokenRBrace, "expected '}' after modes")
	p.consume(lexer.TokenRBrace, "expected '}' after string body")

	end := p.previous().Span.End

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewStringDeclNode(name, modes, span), nil
}

// parseBraneDecl parses a brane declaration.
// Syntax: brane <name> { dimensions [<dim1>, <dim2>, ...] }
func (p *Parser) parseBraneDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenBrane, "expected 'brane'")

	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected brane name")
	}

	p.consume(lexer.TokenLBrace, "expected '{'")

	// Expect "dimensions [ ... ]"
	if !p.match(lexer.TokenDimensions) {
		return nil, p.error("expected 'dimensions' in brane declaration")
	}

	p.consume(lexer.TokenLBracket, "expected '[' after 'dimensions'")

	dimensions := []string{}
	for !p.check(lexer.TokenRBracket) && !p.isAtEnd() {
		if p.match(lexer.TokenComment) {
			continue
		}

		if !p.check(lexer.TokenIdent) {
			return nil, p.error("expected dimension name")
		}
		dimName := p.advance().Lexeme
		dimensions = append(dimensions, dimName)

		// Optional comma
		p.match(lexer.TokenComma)
	}

	p.consume(lexer.TokenRBracket, "expected ']' after dimensions")
	p.consume(lexer.TokenRBrace, "expected '}' after brane body")

	end := p.previous().Span.End

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewBraneDeclNode(name, dimensions, span), nil
}

// parseCouplingDecl parses a coupling declaration.
// Syntax: coupling <string1>.<mode1> <string2>.<mode2> with strength <expr>
func (p *Parser) parseCouplingDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenCoupling, "expected 'coupling'")

	// Parse first string.mode
	if !p.check(lexer.TokenIdent) {
		return nil, p.error("expected first string name")
	}
	stringA := p.parseIdent()

	p.consume(lexer.TokenDot, "expected '.' after string name")

	if !p.check(lexer.TokenIdent) {
		return nil, p.error("expected first mode name")
	}
	modeA := p.parseIdent()

	// Parse second string.mode
	if !p.check(lexer.TokenIdent) {
		return nil, p.error("expected second string name")
	}
	stringB := p.parseIdent()

	p.consume(lexer.TokenDot, "expected '.' after string name")

	if !p.check(lexer.TokenIdent) {
		return nil, p.error("expected second mode name")
	}
	modeB := p.parseIdent()

	// Parse "with strength <expr>"
	if !p.match(lexer.TokenWith) {
		return nil, p.error("expected 'with' after mode pair")
	}

	if !p.match(lexer.TokenStrength) {
		return nil, p.error("expected 'strength' after 'with'")
	}

	strength, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	end := p.previous().Span.End

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewCouplingNode(stringA, modeA, stringB, modeB, strength, span), nil
}

// parseResonanceDecl parses a resonance rule.
// Syntax: resonance when <condition> { <action> }
func (p *Parser) parseResonanceDecl() (ast.Decl, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenResonance, "expected 'resonance'")

	if !p.match(lexer.TokenWhen) {
		return nil, p.error("expected 'when' after 'resonance'")
	}

	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// parseBlockExpr will consume the '{'
	action, err := p.parseBlockExpr()
	if err != nil {
		return nil, err
	}

	end := p.current().Span.End

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewResonanceRuleNode(condition, action, span), nil
}

// ==================== Emit Statement Parsing (M2.3) ====================

// parseEmitStmt parses an emit statement.
// Syntax: emit <TypeName> { <field>: <value>, ... };
func (p *Parser) parseEmitStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenEmit, "expected 'emit'")

	// Parse type name
	typeName := p.parseIdent()
	if typeName == nil {
		return nil, p.error("expected event type name after 'emit'")
	}

	// Parse struct literal
	if !p.check(lexer.TokenLBrace) {
		return nil, p.error("expected '{' after type name in emit statement")
	}

	structLiteral, err := p.parseStructLiteralWithName(typeName)
	if err != nil {
		return nil, err
	}

	// Consume semicolon
	p.consume(lexer.TokenSemicolon, "expected ';' after emit statement")

	end := p.previous().Span.End
	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	structLit, ok := structLiteral.(*ast.StructLiteral)
	if !ok {
		return nil, p.error("expected struct literal in emit statement")
	}

	return ast.NewEmitStmt(typeName, structLit, span), nil
}

// parseUseStmt parses a use statement for calling motifs.
// Syntax: use <motif_name>; or use <motif_name>(<args>);
func (p *Parser) parseUseStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenUse, "expected 'use'")

	// Parse motif name
	motifName := p.parseIdent()
	if motifName == nil {
		return nil, p.error("expected motif name after 'use'")
	}

	// Parse optional arguments
	var args []ast.Expr
	if p.match(lexer.TokenLParen) {
		// Parse argument list
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
		p.consume(lexer.TokenRParen, "expected ')' after arguments")
	}

	// Consume semicolon
	p.consume(lexer.TokenSemicolon, "expected ';' after use statement")

	end := p.previous().Span.End
	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewUseStmt(motifName, args, span), nil
}

// ==================== Entanglement Parsing (M2.5) ====================

// parseEntangleStmt parses an entangle statement.
// Syntax: entangle <var1>, <var2>, ...;
func (p *Parser) parseEntangleStmt() (ast.Stmt, error) {
	start := p.current().Span.Start
	p.consume(lexer.TokenEntangle, "expected 'entangle'")

	variables := []*ast.Ident{}
	for {
		if !p.check(lexer.TokenIdent) {
			return nil, p.error("expected variable name in entangle statement")
		}
		variable := p.parseIdent()
		variables = append(variables, variable)

		if !p.match(lexer.TokenComma) {
			break
		}
	}

	if len(variables) < 2 {
		return nil, p.error("entangle requires at least 2 variables")
	}

	p.consume(lexer.TokenSemicolon, "expected ';' after entangle statement")

	end := p.previous().Span.End

	span := common.Span{
		File:  p.file,
		Start: start,
		End:   end,
	}

	return ast.NewEntangleStmt(variables, span), nil
}

// ==================== Integration with Core Parser ====================

// These methods should be called from the main parseDecl() and parseStmt() methods.

// ParseDeclExtension extends parseDecl to handle new declaration types.
// Returns nil if the token doesn't match any extension declaration.
func (p *Parser) ParseDeclExtension() (ast.Decl, error) {
	switch {
	case p.check(lexer.TokenScore):
		return p.parseScoreDecl()
	case p.check(lexer.TokenMotif):
		return p.parseMotifDecl()
	case p.check(lexer.TokenStringKw):
		return p.parseStringDecl()
	case p.check(lexer.TokenBrane):
		return p.parseBraneDecl()
	case p.check(lexer.TokenCoupling):
		return p.parseCouplingDecl()
	case p.check(lexer.TokenResonance):
		return p.parseResonanceDecl()
	default:
		return nil, nil // Not an extension declaration
	}
}

// ParseStmtExtension extends parseStmt to handle new statement types.
// Returns nil if the token doesn't match any extension statement.
func (p *Parser) ParseStmtExtension() (ast.Stmt, error) {
	switch {
	case p.check(lexer.TokenBars):
		// Allow bars as statements for nesting
		barNode, err := p.parseBarRangeNode()
		if err != nil {
			return nil, err
		}
		// BarRangeNode implements stmtNode(), so we can return it as a statement
		return barNode, nil
	case p.check(lexer.TokenEmit):
		return p.parseEmitStmt()
	case p.check(lexer.TokenUse):
		return p.parseUseStmt()
	case p.check(lexer.TokenEntangle):
		return p.parseEntangleStmt()
	default:
		return nil, nil // Not an extension statement
	}
}

// ==================== Helper Functions ====================

// parseParam parses a function/motif parameter.
// Syntax: <name>: <type>
// Note: This is a helper that wraps the core parser's parameter parsing.
func (p *Parser) parseParam() (*ast.Param, error) {
	start := p.current().Span.Start
	
	name := p.parseIdent()
	if name == nil {
		return nil, p.error("expected parameter name")
	}

	p.consume(lexer.TokenColon, "expected ':' after parameter name")

	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}

	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	param := ast.NewParam(name, typeExpr, span)

	return param, nil
}

// parseField parses a struct/string field.
// Syntax: <name>: <type>
// Note: This is a helper that wraps the core parser's field parsing logic.
func (p *Parser) parseField() (*ast.Field, error) {
	start := p.current().Span.Start
	
	// Allow keywords as field names (similar to struct literals)
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

	p.consume(lexer.TokenColon, "expected ':' after field name")

	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}

	end := p.previous().Span.End
	span := common.NewSpan(start, end, p.file)
	field := ast.NewField(name, typeExpr, span)

	return field, nil
}

// parseBlockExpr is assumed to exist in core_parser.go.
// If not, we need to add it or adapt this code.

