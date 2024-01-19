package evaluator

import (
	"monkeylang/ast"
	"monkeylang/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}

    NULL = &object.Boolean{}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        // statements
        case *ast.Program:
            return evalStatements(node.Statements)
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        case *ast.PrefixExpression:
            right := Eval(node.Right)
            return evalPrefixExpression(node.Operator, right)
        // expression
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value} // returns an integer object of our internal representation of our language.
        case *ast.Boolean:
            return nativeBoolToBooleanObject(node.Value) 
    }

    return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
    case "!" :
        return evalBangOperatorExpression(right)
    case "-" :
        return evalPrefixMinusOperator(right)
    default:
        return NULL
    }
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
        return FALSE
    }
}

func evalPrefixMinusOperator(right object.Object) object.Object {
    if(right.Type() != object.INTEGER_OBJ) {
        return NULL
    } 

    value := right.(*object.Integer).Value

    return &object.Integer{Value: -value}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
    if input {
        return TRUE
    }

    return FALSE
}

func evalStatements(stmts []ast.Statement) object.Object{
    var result object.Object 

    for _, statement := range stmts {
        result = Eval(statement)
    }

    return result
}

