run:
	go run main.go

tl:
	go test ./lexer

tp:
	go test ./parser

te:
	go test ./evaluator

t: 
	go test ./parser ./evaluator ./lexer ./object
