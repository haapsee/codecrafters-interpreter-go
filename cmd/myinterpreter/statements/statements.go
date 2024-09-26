package statements

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type VarStatement struct {
	Name       token.Token
	Expression interfaces.Expr
}

func (vs VarStatement) GetExpression() (interfaces.Expr, error) {
	return vs.Expression, nil
}

func (vs VarStatement) Accept(visitor interfaces.StatementVisitor) (interface{}, error) {
	return visitor.VisitVarStatement(vs)
}

func NewVarStatement(name token.Token, expression interfaces.Expr) VarStatement {
	return VarStatement{
		Name:       name,
		Expression: expression,
	}
}

type ExpressionStatement struct {
	Expression interfaces.Expr
}

func (es ExpressionStatement) GetExpression() (interfaces.Expr, error) {
	return es.Expression, nil
}

func (es ExpressionStatement) Accept(visitor interfaces.StatementVisitor) (interface{}, error) {
	return visitor.VisitExpressionStatement(es)
}

func NewExpressionStatement(expression interfaces.Expr) ExpressionStatement {
	return ExpressionStatement{
		Expression: expression,
	}
}

type PrintStatement struct {
	Expression interfaces.Expr
}

func (ps PrintStatement) GetExpression() (interfaces.Expr, error) {
	return ps.Expression, nil
}

func (ps PrintStatement) Accept(visitor interfaces.StatementVisitor) (interface{}, error) {
	return visitor.VisitPrintStatement(ps)
}

func NewPrintStatement(expression interfaces.Expr) PrintStatement {
	return PrintStatement{
		Expression: expression,
	}
}
