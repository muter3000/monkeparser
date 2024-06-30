package evaluator

import (
	"fmt"
	"github.com/muter3000/monkeparser/pkg/ast"
	"github.com/muter3000/monkeparser/pkg/object"
	"github.com/muter3000/monkeparser/pkg/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node.Value, environment)

	case *ast.CallExpression:
		function := Eval(node.Function, environment)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.LetStatement:
		val := Eval(node.Value, environment)
		if isError(val) {
			return val
		}
		environment.Set(node.Name.String(), val)

	// Blocks
	case *ast.BlockStatement:
		return evalBlockStatement(node, environment)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Environment: environment, Body: body}

	// Return
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, environment)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Environment)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func evalIdentifier(value string, environment *object.Environment) object.Object {
	val, exists := environment.Get(value)
	if !exists {
		return newError("identifier not found: %s", value)
	}
	return val
}

func evalIfExpression(ie *ast.IfExpression, environment *object.Environment) object.Object {
	pred := Eval(ie.Predicate, environment)
	if isError(pred) {
		return pred
	}
	if pred == NULL {
		return NULL
	}

	if isTruthy(pred) {
		return Eval(ie.Consequence, environment)
	}
	if ie.Alternative == nil {
		return NULL
	}
	return Eval(ie.Alternative, environment)
}

func isTruthy(pred object.Object) bool {
	switch pred {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		if pred.Type() == object.INTEGER_OBJ {
			return pred.(*object.Integer).Value != 0
		}
		return true
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	lValue := left.(*object.Boolean).Value
	rValue := right.(*object.Boolean).Value

	switch operator {
	case token.EQ:
		return nativeBoolToBooleanObject(lValue == rValue)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(lValue != rValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	lValue := left.(*object.Integer).Value
	rValue := right.(*object.Integer).Value
	switch operator {
	// Boolean logic
	case token.EQ:
		return nativeBoolToBooleanObject(lValue == rValue)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(lValue != rValue)
	case token.LT:
		return nativeBoolToBooleanObject(lValue < rValue)
	case token.GT:
		return nativeBoolToBooleanObject(lValue > rValue)
	case token.LTE:
		return nativeBoolToBooleanObject(lValue <= rValue)
	case token.GTE:
		return nativeBoolToBooleanObject(lValue >= rValue)

	// Math
	case token.PLUS:
		return &object.Integer{Value: lValue + rValue}
	case token.SUB:
		return &object.Integer{Value: lValue - rValue}
	case token.MUL:
		return &object.Integer{Value: lValue * rValue}
	case token.DIV:
		return &object.Integer{Value: lValue / rValue}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.BANG:
		return evalBangOperatorExpression(right)
	case token.SUB:
		return evalSubOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalSubOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		if right.Type() == object.INTEGER_OBJ {
			if right.(*object.Integer).Value != 0 {
				return FALSE
			}
			return TRUE
		}
		return FALSE
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalBlockStatement(block *ast.BlockStatement, environment *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, environment)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}
func evalProgram(program *ast.Program, environment *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, environment)
		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}
