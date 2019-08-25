package eval

import (
	"fmt"
)

type executeKind int

//type of exp
const (
	SELFVAL executeKind = iota
	CALL
	ANNOTATION
)

func ast(sentence []Token) *node {
	root := new(node)
	p := root

	opStack := NewStack()
	pStack := NewStack()

	for i := 0; i < len(sentence); i++ {
		// if pStack.IsEmpty() {
		// 	if sentence[i].Kind != LPAREN {
		// 		panic(fmt.Sprintf("extra symbol:%s", sentence[i].String()))
		// 	}
		// }

		if sentence[i].Kind == LPAREN {
			if i+1 < len(sentence) {
				pStack.Push(sentence[i])
				if opStack.IsEmpty() {
					p.tok = sentence[i+1]
					p.executeType = CALL
				} else {
					node := token2Node(sentence[i+1])
					p.childen = append(p.childen, node)
					node.executeType = CALL
					p = node
				}
				opStack.Push(p)
				i++
				continue
			} else {
				panic("lack of right paren")
			}
		} else if sentence[i].Kind == RPAREN {
			if pStack.IsEmpty() {
				panic("extra right paren")
			} else {
				pStack.Pop()
				opStack.Pop()
				if opStack.IsEmpty() {
					p = nil
				} else {
					p = opStack.Top().(*node)
				}
			}
		} else {
			node := token2Node(sentence[i])
			node.executeType = SELFVAL

			if opStack.IsEmpty() {
				root = node
			} else {
				p.childen = append(p.childen, node)
			}
		}
	}
	return root
}

func getVarName(tree *node) string {
	if len(tree.childen) < 2 {
		panic("define syntax error")
	}
	return tree.childen[0].tok.String()
}

func getVarValue(tree *node, env *environment) interface{} {
	return tree.childen[1].value(env)
}

func eval(sentence []Token, e *environment) interface{} {
	syntaxTree := ast(sentence)
	rel := syntaxTree.value(e)
	return rel
}

func expsEval(nodes []*node, env *environment) interface{} {
	var output interface{}
	for _, node := range nodes {
		output = node.value(env)
	}
	return output
}

func runCall(cal *callStruct, args []interface{}) interface{} {
	// defer destroyEnv(call.env)

	call := copyCall(cal)

	if len(call.args) != len(args) {
		panic(fmt.Sprintf("%s args num is %d", call.name, len(call.args)))
	}

	for i := 0; i < len(args); i++ {
		call.env.CreateVar(call.args[i], args[i])
	}

	return expsEval(call.body, call.env)
}

func expType(sentence []Token) executeKind {
	if sentence[0].Kind == SEMICOLON {
		return ANNOTATION
	}

	// if sentence[0].Kind != RPAREN {
	// 	return VARI
	// }

	return CALL
}
