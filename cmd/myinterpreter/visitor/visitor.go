package visitor

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expr"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/functions"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

/****************
 ****************
 *			  	*
 *	Interpreter	*
 *             	*
 ****************
 ****************/
type Interpreter struct{}

func checkNumberOperand(t token.Token, operands ...interface{}) error {
	for _, operand := range operands {
		if functions.TypeOf(operand) != "float64" {
			return errors.NewRuntimeError(t, "Operand must be a number.")
		}
	}
	return nil
}

func (interpreter *Interpreter) evaluate(expression interfaces.Expr) (interface{}, error) {
	return expression.Accept(interpreter)
}

func (interpreter *Interpreter) VisitBinaryExpr(b interfaces.Expr) (interface{}, error) {
	binary := b.(expr.BinaryExpr)
	left, err := interpreter.evaluate(binary.Left)
	if err != nil {
		return nil, err
	}

	right, err := interpreter.evaluate(binary.Right)
	if err != nil {
		return nil, err
	}

	switch binary.Operator.TokenType {
	case token.BANG:
		return functions.IsTruthy(right), nil
	case token.MINUS:
		err := checkNumberOperand(binary.Operator, right, left)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.PLUS:
		isSame := functions.TypeOf(left) == functions.TypeOf(right)
		isString, isFloat := functions.TypeOf(left) == "string", functions.TypeOf(left) == "float64"

		if isSame && isString {
			return left.(string) + right.(string), nil
		}

		if isSame && isFloat {
			return left.(float64) + right.(float64), nil
		}
		return nil, errors.NewRuntimeError(binary.Operator, "Operands must be two numbers or two strings.")
	case token.SLASH:
		err := checkNumberOperand(binary.Operator, right, left)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.STAR:
		err := checkNumberOperand(binary.Operator, right, left)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.GREATER:
		err := checkNumberOperand(binary.Operator, right, left)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		err := checkNumberOperand(binary.Operator, right, left)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		err := checkNumberOperand(binary.Operator, right, left)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		err := checkNumberOperand(binary.Operator, right, left)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.BANG_EQUAL:
		return left != right, nil
	case token.EQUAL_EQUAL:
		return left == right, nil
	}

	return nil, nil
}

// VisitGroupingExpr implements interfaces.Visitor.
func (interpreter *Interpreter) VisitGroupingExpr(g interfaces.Expr) (interface{}, error) {
	grouping := g.(expr.GroupingExpr)
	return interpreter.evaluate(grouping.Expression)
}

// VisitLiteralExpr implements interfaces.Visitor.
func (interpreter *Interpreter) VisitLiteralExpr(l interfaces.Expr) (interface{}, error) {
	literal := l.(expr.LiteralExpr)
	return literal.Literal, nil
}

// VisitUnaryExpr implements interfaces.Visitor.
func (interpreter *Interpreter) VisitUnaryExpr(u interfaces.Expr) (interface{}, error) {
	unary := u.(expr.UnaryExpr)
	right, err := interpreter.evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	if unary.Operator.TokenType == token.BANG {
		return !functions.IsTruthy(right), nil
	}

	if unary.Operator.TokenType == token.MINUS {
		err := checkNumberOperand(unary.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}
	return nil, nil
}

func (interpreter *Interpreter) Stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	if functions.TypeOf(obj) == "float64" {
		result := functions.FormatWithFixedPrecision(obj.(float64))

		if result[len(result)-2:] == ".0" {
			result = result[0 : len(result)-2]
		}

		return result
	}
	return fmt.Sprintf("%v", obj)
}

func (interpreter *Interpreter) Interpret(expression interfaces.Expr) error {
	value, err := interpreter.evaluate(expression)
	if err != nil {
		return err
	}
	fmt.Println(interpreter.Stringify(value))
	return nil
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

/****************
 ****************
 *			  	*
 *	AstPrinter	*
 *             	*
 ****************
 ****************/
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
