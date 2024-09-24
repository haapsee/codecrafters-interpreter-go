package parser

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expr"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Parser struct {
	Tokens  []token.Token
	Current int
}

func (p *Parser) error(t token.Token, message string) error {
	if t.TokenType == token.EOF {
		return errors.NewLexicalError(t.Line, " at end", message)
	} else {
		return errors.NewLexicalError(t.Line, " at '"+t.Lexeme+"'", message)
	}
}

func (p *Parser) isAtEnd() bool {
	return p.Current >= len(p.Tokens)
}

func (p *Parser) previous() token.Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) peek() token.Token {
	return p.Tokens[p.Current]
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) check(tokentype token.TokenType) bool {
	return !p.isAtEnd() && p.peek().TokenType == tokentype
}

func (p *Parser) match(tokentypes ...token.TokenType) bool {
	for _, tokentype := range tokentypes {
		if p.check(tokentype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokentype token.TokenType, message string) (token.Token, error) {
	if p.check(tokentype) {
		return p.advance(), nil
	}
	return token.NewTokenNil(), errors.New(message)
}

func (p *Parser) primary() (interfaces.Expr, error) {
	if p.match(token.FALSE) {
		return expr.NewLiteral(false), nil
	} else if p.match(token.TRUE) {
		return expr.NewLiteral(true), nil
	} else if p.match(token.NIL) {
		return expr.NewLiteral(nil), nil
	} else if p.match(token.NUMBER, token.STRING) {
		return expr.NewLiteral(p.previous().Literal), nil
	} else if p.match(token.LEFT_PAREN) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return expr.NewGrouping(expression), nil
	}
	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) unary() (interfaces.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expr.NewUnary(operator, right), nil
	}

	primary, err := p.primary()
	if err != nil {
		return nil, err
	}
	return primary, nil
}

func (p *Parser) factor() (interfaces.Expr, error) {
	expression, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expression = expr.NewBinary(expression, operator, right)
	}
	return expression, nil
}

func (p *Parser) term() (interfaces.Expr, error) {
	expression, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = expr.NewBinary(expression, operator, right)
	}
	return expression, nil
}

func (p *Parser) comparsion() (interfaces.Expr, error) {
	expression, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expression = expr.NewBinary(expression, operator, right)
	}
	return expression, nil
}

func (p *Parser) expression() (interfaces.Expr, error) {
	expression, err := p.comparsion()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparsion()
		if err != nil {
			return nil, err
		}
		expression = expr.NewBinary(expression, operator, right)
	}

	return expression, nil
}

func (p *Parser) Parse() (interfaces.Expr, error) {
	return p.expression()
}

func New(tokens []token.Token) Parser {
	return Parser{
		Tokens:  tokens,
		Current: 0,
	}
}
