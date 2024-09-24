package visitor

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expr"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/functions"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
)

type AstPrinter struct {
}

func (printer *AstPrinter) parenthesize(name string, expressions ...interfaces.Expr) (string, error) {
	str := fmt.Sprintf("(%s", name)

	for _, expression := range expressions {
		result, err := expression.Accept(printer)
		if err != nil {
			return "", err
		}
		str = fmt.Sprintf("%s %v", str, result)
	}

	return str + ")", nil
}

func (printer *AstPrinter) VisitBinaryExpr(b interfaces.Expr) (interface{}, error) {
	binary := b.(expr.BinaryExpr)
	result, err := printer.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (printer *AstPrinter) VisitGroupingExpr(g interfaces.Expr) (interface{}, error) {
	grouping := g.(expr.GroupingExpr)
	result, err := printer.parenthesize("group", grouping.Expression)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (printer *AstPrinter) VisitLiteralExpr(l interfaces.Expr) (interface{}, error) {
	literal := l.(expr.LiteralExpr)
	if literal.Literal == nil {
		return "nil", nil
	}
	switch literal.Literal.(type) {
	case bool:
		value := literal.Literal.(bool)
		return strconv.FormatBool(value), nil
	case float64:
		value := literal.Literal.(float64)
		formattedValue := functions.FormatWithFixedPrecision(value)
		return formattedValue, nil
	case string:
		value := literal.Literal.(string)
		return value, nil
	default:
		return "", nil
	}
}

func (printer *AstPrinter) VisitUnaryExpr(u interfaces.Expr) (interface{}, error) {
	unary := u.(expr.UnaryExpr)
	result, err := printer.parenthesize(unary.Operator.Lexeme, unary.Right)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (printer *AstPrinter) Print(expression interfaces.Expr) (string, error) {
	result, err := expression.Accept(printer)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}
