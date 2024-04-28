package parser

import (
	"fmt"
	"strconv"

	"dojo/ast"
	"dojo/lexer"
	"dojo/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > OR <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X OR !X
	CALL        // fooFunc(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParserFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token

	// With these maps in place, we can just check
	// if the appropriate map (infix or prefix) has a parsing
	// function associated with curToken.Type.

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParserFns map[token.TokenType]infixParserFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read the two tokens, so curr and peek are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.", t, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // parent node
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

// registerPrefix helper func to add entries for the prefixParseFnx map of the parser
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParserFn) {
	p.infixParserFns[tokenType] = fn
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseExpression checks whether we have a parsing function associated
// with p.curToken.Type in the prefix position.
// If we do, it calls this parsing function, if not, it returns nil.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// check for optional semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//	TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	//	TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek  checks the type of the peekToken and only if the type is correct does it advance the tokens
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = val
	return lit
}

// function parseProgram() {
// 	program = newProgramASTNode()
// 	advanceTokens() --> get more tokens by moving the peek and curr
// 	for (currentToken() != EOF_TOKEN) {
// 		statement = null
// 		if (currentToken() == LET_TOKEN) {
// 		statement = parseLetStatement()
// 		} else if (currentToken() == RETURN_TOKEN) {
// 		statement = parseReturnStatement()
// 		} else if (currentToken() == IF_TOKEN) {
// 		statement = parseIfStatement()
// 		}
// 		if (statement != null) {
// 		program.Statements.push(statement)
// 		}
// 		advanceTokens()
// 	}
// 	return program
// }

// function parseLetStatement() {
// 	advanceTokens()
// 	identifier = parseIdentifier()
// 	advanceTokens()
// 	if currentToken() != EQUAL_TOKEN {
// 		parseError("no equal sign!")
// 		return null
// 	}
// 	advanceTokens()
// 	value = parseExpression()
// 	variableStatement = newVariableStatementASTNode()
// 	variableStatement.identifier = identifier
// 	variableStatement.value = value
// 	return variableStatement
// }

// function parseIdentifier() {
// 	identifier = newIdentifierASTNode()
// 	identifier.token = currentToken()
// 	return identifier
// }

// function parseExpression() {
// 	if (currentToken() == INTEGER_TOKEN) {
// 		if (nextToken() == PLUS_TOKEN)
// 		{
// 			return parseOperatorExpression()
// 		}else if (nextToken() == SEMICOLON_TOKEN) {
// 			return parseIntegerLiteral()
// 		}
// 	} else if (currentToken() == LEFT_PAREN) {
// 		return parseGroupedExpression()
// 	}
// 	// [...]
// }
// function parseOperatorExpression() {
// 	operatorExpression = newOperatorExpression()
// 	operatorExpression.left = parseIntegerLiteral()
// 	operatorExpression.operator = currentToken()
// 	operatorExpression.right = parseExpression()
// 	return operatorExpression()
// }
// 	// [...]
