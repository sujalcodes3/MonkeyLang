package evaluator

import (
	"monkeylang/lexer"
	"monkeylang/object"
	"monkeylang/parser"
	"testing"
)

func TestEvalBangOperator(t * testing.T) {
    tests := [] struct {
        input string
        expected bool 
    } {
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true}, 
        {"!!false", false},
        {"!!5", true},
    }
    

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}



// for integer expressions
func TestEvalIntegerExpression(t * testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"5", 5},
        {"4", 4},
        {"-4", -4},
        {"-5", -5},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input) 
        testIntegerObject(t, evaluated, tt.expected)
    }
}
// for boolean expressions
func TestEvalBooleanExpression(t * testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"true", true},
        {"false", false},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input) 
        testBooleanObject(t, evaluated, tt.expected)
    }
}



// for integer expressions
func testIntegerObject(t * testing.T, evaluated object.Object, expected int64) bool {
    result, ok := evaluated.(*object.Integer)

    if !ok {
        t.Errorf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
        return false
    }
        
    if result.Value != expected {
        t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
        return false
    }

    return true
}

// for boolean expressions
func testBooleanObject(t * testing.T, evaluated object.Object, expected bool) bool {
    result, ok := evaluated.(*object.Boolean)

    if !ok {
        t.Errorf("object is not a boolean, got= %T (%+v)", evaluated, evaluated)
        return false
    }

    if result.Value != expected {
        t.Errorf("object has wrong value, got=%t, want= %t", result.Value, expected) 
        return false
    }

    return true
}

func testEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    return Eval(program)
} 
