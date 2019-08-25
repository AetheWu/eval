package eval

import (
	"bufio"
	"fmt"
	"os"
)

func EvalLoopDriver(input string) {
	env := initEnv()
	sentences := CodeSplit(input)
	var rel interface{}
	for _, exp := range sentences {
		sentence := buildSentence(exp)
		switch expType(sentence) {
		case ANNOTATION:
			continue
		case CALL:
			rel = executeCall(sentence, env)
		}
		evalDisplay(rel)
	}

}

func RunSourceFile(filename string) {
	input := FileLoad(filename)
	EvalLoopDriver(input)
}

func Interpreter() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			Interpreter()
		}
	}()

	var env = initEnv()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("input:")
		input, _ := reader.ReadString('\n')
		sentence := buildSentence(input)
		eval(sentence, env)
	}
}

func executeVari(sentence []Token, env *environment) interface{} {
	var output interface{}
	tok := sentence[0]

	switch tok.Kind {
	case INT, FLOAT, STRING:
		output = tok.Val
	case IDENT:
		output = env.GetVar(tok.String())
	}

	return output
}

func executeCall(sentence []Token, env *environment) interface{} {
	node := ast(sentence)
	return node.value(env)
}
