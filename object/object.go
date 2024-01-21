package object

import "fmt"

const (
    INTEGER_OBJ = "INTEGER"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ    = "NULL"
    RETURN_VALUE_OBJ = "RETURN_VALUE"
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



