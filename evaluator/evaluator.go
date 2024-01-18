package evaluator

import (
	"monkeylang/ast"
	"monkeylang/object"
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        // statements
        case *ast.Program:
            return evalStatements(node.Statements)
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        // expression
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value} // returns an integer object of our internal representation of our language.
    }

    return nil
}

func evalStatements(stmts []ast.Statement) object.Object{
    var result object.Object 

    for _, statement := range stmts {
        result = Eval(statement)
    }

    return result
}
