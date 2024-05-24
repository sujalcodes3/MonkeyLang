package object

import (
	"bytes"
	"fmt"
	"monkeylang/ast"
	"strings"

    "hash/fnv"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
    HASH_OBJ         = "HASH"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

// begin Integer Data Type -> satisfies the Object interface
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// end Integer Data Type

// begin Boolean Data Type -> satisfies the Object interface
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// end Boolean Data Type

// begin NULL Data Type -> satisfies the Object interface
type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// end NULL Data Type

// begin RETURNVALUE Data Type -> satisfies the Object interface
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// end RETURNVALUE Data Type

// begin Error Data Type -> satisfies the Error interface
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// end Error Data Type

// begin Function Data Type -> satisfies the Object interface

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
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

// begin String Data Type

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// end String Data Type

// BuiltinFunction

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// end BuiltinFunction

// start Array
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

//end Array

type Hashable interface { 
    HashKey() HashKey
    /* only implemented by 
    *object.String
    *object.Boolean
    *object.Intege
    */
}

type HashKey struct {
    Type ObjectType
    Value uint64
}

func (b * Boolean) HashKey() HashKey {
    var value uint64 

    if b.Value {
        value = 1
    } else {
        value = 0
    }

    return HashKey{Type: b.Type(), Value: value}
}

func (i * Integer) HashKey() HashKey {
    return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s * String) HashKey() HashKey {
    h := fnv.New64a()
    h.Write([]byte(s.Value))

    return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// start HASH

type HashPair struct {
    Key Object
    Value Object
}

type Hash struct {
    Pairs map[HashKey]HashPair
}
    
func (h * Hash) Type() ObjectType { return HASH_OBJ }

func (h * Hash) Inspect() string {
    var out bytes.Buffer

    pairs := []string{}

    for _, pair := range h.Pairs {
        pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", ")) 
    out.WriteString("}")

    return out.String()
}
