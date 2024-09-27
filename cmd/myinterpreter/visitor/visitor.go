package visitor

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/environment"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expr"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/functions"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/statements"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

/****************
 ****************
 *			  	*
 *	Interpreter	*
 *             	*
 ****************
 ****************/
type Interpreter struct {
	environment environment.Environment
}

func (interpreter Interpreter) VisitAssignExpr(ae interfaces.Expr) (interface{}, error) {
	// fmt.Println(ae)
	assignExpr := ae.(expr.AssignExpr)
	value, err := interpreter.evaluate(assignExpr.Value)
	if err != nil {
		return nil, err
	}
	until := true
	for until {
		switch v := value.(type) {
		case interfaces.Expr:
			value, err = v.Accept(interpreter)
			if err != nil {
				return nil, err
			}
		default:
			until = false
		}
	}
	err = interpreter.environment.Assign(assignExpr.Name, value)
	if err != nil {
		return nil, err
	}
	return assignExpr.Value, nil
}

func (interpreter Interpreter) VisitVarExpr(v interfaces.Expr) (interface{}, error) {
	varExpression := v.(expr.VarExpr)
	return interpreter.environment.Get(varExpression.Token)
}

func (interpreter Interpreter) VisitVarStatement(varStmt interfaces.Statement) (interface{}, error) {
	var value interface{}
	var err error

	varStatement := varStmt.(statements.VarStatement)
	expression := varStatement.Expression

	if expression != nil {
		value, err = interpreter.evaluate(expression)
		if err != nil {
			return nil, err
		}
	}

	interpreter.environment.Define(varStatement.Name.Lexeme, value)
	return nil, nil
}

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

func (interpreter Interpreter) VisitBinaryExpr(b interfaces.Expr) (interface{}, error) {
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
		// fmt.Println("!=", left != right)
		return left != right, nil
	case token.EQUAL_EQUAL:
		return left == right, nil
	}

	return nil, nil
}

func (interpreter Interpreter) VisitGroupingExpr(g interfaces.Expr) (interface{}, error) {
	grouping := g.(expr.GroupingExpr)
	return interpreter.evaluate(grouping.Expression)
}

func (interpreter Interpreter) VisitLiteralExpr(l interfaces.Expr) (interface{}, error) {
	literal := l.(expr.LiteralExpr)
	return literal.Literal, nil
}

func (interpreter Interpreter) VisitUnaryExpr(u interfaces.Expr) (interface{}, error) {
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

func (interpreter Interpreter) VisitPrintStatement(printStmt interfaces.Statement) (interface{}, error) {
	printStatement := printStmt.(statements.PrintStatement)
	value, err := interpreter.evaluate(printStatement.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(interpreter.Stringify(value))
	return value, nil
}

func (interpreter Interpreter) VisitExpressionStatement(exprStmt interfaces.Statement) (interface{}, error) {
	expressionStatement := exprStmt.(statements.ExpressionStatement)
	return interpreter.evaluate(expressionStatement.Expression)
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

	switch obj := obj.(type) {
	case expr.LiteralExpr:
		v := obj.Literal
		return fmt.Sprintf("%v", v)
	}

	return fmt.Sprintf("%v", obj)
}

func (interpreter *Interpreter) execute(statement interfaces.Statement) error {
	_, err := statement.Accept(interpreter)
	return err
}

func (interpreter *Interpreter) Interpret(statements []interfaces.Statement) error {
	for _, statement := range statements {
		err := interpreter.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (interpreter *Interpreter) Interpret(expression interfaces.Expr) error {
// 	value, err := interpreter.evaluate(expression)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(interpreter.Stringify(value))
// 	return nil
// }

func NewInterpreter() Interpreter {
	return Interpreter{
		environment: environment.NewEnvironment(),
	}
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

// VisitAssignExpr implements interfaces.Visitor.
func (printer *AstPrinter) VisitAssignExpr(ae interfaces.Expr) (interface{}, error) {
	panic("unimplemented")
}

// VisitVarStatement implements interfaces.StatementVisitor.
func (printer *AstPrinter) VisitVarStatement(varStmt interfaces.Statement) (interface{}, error) {
	panic("unimplemented")
}

// VisitVarExpr implements interfaces.Visitor.
func (printer *AstPrinter) VisitVarExpr(v interfaces.Expr) (interface{}, error) {
	panic("unimplemented")
}

func (printer *AstPrinter) VisitExpressionStatement(exprStmt interfaces.Statement) (interface{}, error) {
	expressionStatement := exprStmt.(statements.ExpressionStatement)
	return printer.parenthesize(";", expressionStatement.Expression)
}

func (printer *AstPrinter) VisitPrintStatement(printStmt interfaces.Statement) (interface{}, error) {
	printStatement := printStmt.(statements.PrintStatement)
	return printer.parenthesize("print", printStatement.Expression)
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

func (printer *AstPrinter) Print(obj interface{}) (string, error) {
	switch obj := obj.(type) {
	case interfaces.Expr:
		result, err := obj.Accept(printer)
		if err != nil {
			return "", err
		}
		return result.(string), nil
	case interfaces.Statement:
		result, err := obj.Accept(printer)
		if err != nil {
			return "", err
		}
		return result.(string), nil
	default:
		return "", errors.NewRuntimeError(token.NewTokenNil(), fmt.Sprintf("%v", obj))
	}
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}
