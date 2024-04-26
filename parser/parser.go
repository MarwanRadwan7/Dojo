package parser

import (
	"dojo/ast"
	"dojo/lexer"
	"dojo/token"
	"fmt"
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read the two tokens, so curr and peek are both set
	p.nextToken()
	p.nextToken()

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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}

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
