package repl

import (
	"fmt"
	"io"
	"monkeylang/evaluator"
	"monkeylang/lexer"
	"monkeylang/object"
	"monkeylang/parser"
	"os"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
    dat, err := os.ReadFile("scratch/input.mkl")
    if err != nil {
        panic(err)
    }
    env := object.NewEnvironment()
    l := lexer.New(string(dat))
    p := parser.New(l)

    program := p.ParseProgram()

    if len(p.Errors()) != 0 {
        printParserErrors(out, p.Errors())
    }
    
    evaluated := evaluator.Eval(program, env)

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

func printParserErrors (out io.Writer, errors []string) {
    for _, msg := range errors {
        _, err := io.WriteString(out, "\t" + msg + "\n")

        if err != nil {
            fmt.Printf("Error writing errors to output: %s", err)
        }
    }
}
/*func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()

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

        evaluated := evaluator.Eval(program, env)

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
}*/

