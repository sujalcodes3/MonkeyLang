package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkeylang/evaluator"
	"monkeylang/lexer"
	"monkeylang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        if !scanned {
            return 
        }

        line := scanner.Text()
        l := lexer.New(line)
        p := parser.New(l)

        program := p.ParseProgram()

        if len(p.Errors()) != 0 {
            printParserErrors(out, p.Errors())
            continue
        }
        
        evaluated := evaluator.Eval(program)

        if evaluated != nil {
            _, err := io.WriteString(out, evaluated.Inspect())
            if err != nil {
                fmt.Printf("Error writing the output: %s", err)
            }
            _, err =  io.WriteString(out, "\n")
            if err != nil {
                fmt.Printf("Error writing the output: %s", err)
            }
        }
    }
}


func printParserErrors (out io.Writer, errors []string) {
    for _, msg := range errors {
        _, err := io.WriteString(out, "\t" + msg + "\n")

        if err != nil {
            fmt.Printf("Error writing errors to output: %s", err)
        }
    }
}
