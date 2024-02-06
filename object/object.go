package object

import (
	"bytes"
	"fmt"
	"monkeylang/ast"
	"strings"
)

const (
    INTEGER_OBJ = "INTEGER"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ    = "NULL"
    RETURN_VALUE_OBJ = "RETURN_VALUE"
    ERROR_OBJ = "ERROR"
    FUNCTION_OBJ = "FUNCTION"
)

type ObjectType string

type Object interface {
    Type()    ObjectType
    Inspect() string
}


// begin Integer Data Type -> satisfies the Object interface
type Integer struct {
    Value int64
}

func (i * Integer) Inspect() string {
    return fmt.Sprintf("%d", i.Value)
}

func (i * Integer) Type() ObjectType {
    return INTEGER_OBJ
}
// end Integer Data Type



// begin Boolean Data Type -> satisfies the Object interface
type Boolean struct {
    Value bool 
}

func (b * Boolean) Inspect() string {
    return fmt.Sprintf("%t", b.Value)
}

func (b * Boolean) Type() ObjectType {
    return BOOLEAN_OBJ
}
// end Boolean Data Type



// begin NULL Data Type -> satisfies the Object interface
type Null struct {}

func (n * Null) Inspect() string {
    return "null" 
}

func (n * Null) Type() ObjectType {
    return NULL_OBJ
}
// end NULL Data Type




// begin RETURNVALUE Data Type -> satisfies the Object interface
type ReturnValue struct {
    Value Object
}

func (rv * ReturnValue) Type() ObjectType {
    return RETURN_VALUE_OBJ
}

func (rv * ReturnValue) Inspect() string {
    return rv.Value.Inspect()
}
// end RETURNVALUE Data Type





// begin Error Data Type -> satisfies the Error interface
type Error struct {
    Message string
}

func (e * Error) Type() ObjectType {
    return ERROR_OBJ
}
func (e * Error) Inspect() string {
    return "ERROR: " + e.Message 
}
// end Error Data Type




// begin Function Data Type -> satisfies the Object interface

type Function struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    Env        *Environment
}

func (f * Function) Type() ObjectType {
    return FUNCTION_OBJ
}

func (f * Function) Inspect() string {
    var out bytes.Buffer
    params := []string{}
    for _, p := range f.Parameters {
        params = append(params, p.String())
    }
    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(f.Body.String())
    out.WriteString("\n}")

    return out.String()
}


// end Function Data Type
