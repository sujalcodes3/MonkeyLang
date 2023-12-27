# This file containst the various elements the this code base Contains..

let <identifier> = <expression>

## Identifiers 
    Basically variable names

## Keywords 
    Things that seem like Identifiers but are part of the language itself.

## Expressions
    expression produce values

## Statements
    does not produce values

### Difference b/w Statements and expression
    for e.g. : let a = 5; doesn't produce a value but add(5, 5) produces a value.

## IntegerLiterals 
    for e.g : 5, 10 etc    

## Lexer
    The following is the output from the lexer. Our parser traverses through this data to create the ast.
    >> let x = -1

    {Type:LET Literal:let}
    {Type:IDENT Literal:x}
    {Type:= Literal:=}
    {Type:- Literal:-}
    {Type:INT Literal:1}
