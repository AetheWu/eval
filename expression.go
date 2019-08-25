package eval

import "fmt"

type node struct {
	tok         Token
	childen     []*node
	executeType executeKind
}

func (n *node) kind() kind {
	return n.tok.Kind
}

func newNode(scan *StrScan) *node {
	root := new(node)
	var children []*node

	root.tok = NewToken(scan)

	root.childen = children
	return root
}

func token2Node(t Token) *node {
	root := new(node)
	root.childen = []*node{}
	root.tok = t
	return root
}

func (n *node) value(e *environment) interface{} {
	var rel interface{}
	switch n.kind() {

	//built-in type
	case FLOAT, INT, STRING:
		if n.executeType != SELFVAL {
			panic(callNotProcedureError(n))
		}
		rel = n.tok.Val

	//normal operations
	case ADD:
		rel = handleADD(n, e)
	case SUB:
		rel = handleSUB(n, e)
	case MUL:
		rel = handleMUL(n, e)
	case QUO:
		rel = handleQUO(n, e)

	case EQL:
		rel = handleEQL(n, e)
	case LSS:
		rel = handleLSS(n, e)
	case GTR:
		rel = handleGTR(n, e)
	case NOT:
		rel = handleNOT(n, e)

	//vari and call
	case IDENT:
		rel = handleVari(n, e)

	//keywords
	case DEFINE:
		rel = handleDEFINE(n, e)
	case SET:
		handleSET(n, e)
	case IF:
		rel = handleIF(n, e)
	case LAMBDA:
		rel = handleLambda(n, e)
	case BEGIN:
		rel = handleBEGIN(n, e)
	}
	return rel
}

func nodeChildren2Args(n *node, e *environment) []interface{} {
	var output = make([]interface{}, len(n.childen))

	for i, child := range n.childen {
		output[i] = child.value(e)
	}

	return output
}

func typeCheck(n *node) {

	if n.childen == nil {
		panic(fmt.Sprintf("call: %s need more args", n.tok.String()))
	}

	if len(n.childen) > 0 {
		tmp0 := n.childen[0].kind()
		for _, child := range n.childen {
			tmp := child.kind()
			if child.kind() != FLOAT && child.kind() != INT && tmp0 != tmp {
				panic(fmt.Sprintf("call: %s args type conflict", n.tok.String()))
			}
		}
	}
}

func argsToFloat(n *node, e *environment) []float64 {

	var output = make([]float64, len(n.childen))

	for i := 0; i < len(n.childen); i++ {
		switch n.childen[i].value(e).(type) {
		case int:
			output[i] = float64(n.childen[i].value(e).(int))
		case float64:
			output[i] = n.childen[i].value(e).(float64)
		default:
			panic(fmt.Sprintf("op %s args type is int or float64", n.tok.String()))
		}
	}

	return output
}

func argsToInt(n *node, e *environment) []int {

	var output = make([]int, len(n.childen))

	for i := 0; i < len(n.childen); i++ {
		output[i] = n.childen[i].value(e).(int)
	}

	return output
}

func isArgsFloat(n *node, e *environment) bool {
	for _, child := range n.childen {
		if _, ok := child.value(e).(float64); !ok {
			return false
		}
	}
	return true
}

func isArgsInt(n *node, e *environment) bool {
	for _, child := range n.childen {
		if _, ok := child.value(e).(int); !ok {
			return false
		}
	}
	return true
}
