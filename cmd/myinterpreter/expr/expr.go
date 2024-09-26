package expr

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type VarExpr struct {
	Token token.Token
}

func (ve VarExpr) Accept(v interfaces.Visitor) (interface{}, error) {
	return v.VisitVarExpr(ve)
}

func NewVarExpr(t token.Token) VarExpr {
	return VarExpr{
		Token: t,
	}
}

type GroupingExpr struct {
	Expression interfaces.Expr
}

func (g GroupingExpr) Accept(v interfaces.Visitor) (interface{}, error) {
	result, err := v.VisitGroupingExpr(g)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewGrouping(expression interfaces.Expr) GroupingExpr {
	return GroupingExpr{
		Expression: expression,
	}
}

type LiteralExpr struct {
	Literal interface{}
}

func (l LiteralExpr) Accept(v interfaces.Visitor) (interface{}, error) {
	result, err := v.VisitLiteralExpr(l)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewLiteral(literal interface{}) LiteralExpr {
	return LiteralExpr{
		Literal: literal,
	}
}

type BinaryExpr struct {
	Left     interfaces.Expr
	Right    interfaces.Expr
	Operator token.Token
}

func (b BinaryExpr) Accept(v interfaces.Visitor) (interface{}, error) {
	result, err := v.VisitBinaryExpr(b)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewBinary(left interfaces.Expr, operator token.Token, right interfaces.Expr) BinaryExpr {
	return BinaryExpr{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

type UnaryExpr struct {
	Operator token.Token
	Right    interfaces.Expr
}

func (u UnaryExpr) Accept(v interfaces.Visitor) (interface{}, error) {
	result, err := v.VisitUnaryExpr(u)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewUnary(operator token.Token, right interfaces.Expr) UnaryExpr {
	return UnaryExpr{
		Operator: operator,
		Right:    right,
	}
}
