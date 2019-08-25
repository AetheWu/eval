package eval

import (
	"fmt"
)

func handleEQL(n *node, e *environment) interface{} {

	if n.childen == nil || len(n.childen) < 2 {
		panic(fmt.Sprintf("op: %s need 2 or more args", n.tok.String()))
	}

	for i := 1; i < len(n.childen); i++ {
		if n.childen[0].value(e) != n.childen[i].value(e) {
			return false
		}
	}

	return true
}

func handleLSS(n *node, e *environment) interface{} {
	if n.childen == nil || len(n.childen) < 2 {
		panic(fmt.Sprintf("op: %s need 2 or more args", n.tok.String()))
	}

	args := argsToFloat(n, e)

	for i := 1; i < len(args); i++ {
		if args[0] >= args[i] {
			return false
		}
	}

	return true
}

func handleGTR(n *node, e *environment) interface{} {
	if n.childen == nil || len(n.childen) < 2 {
		panic(fmt.Sprintf("op: %s need 2 or more args", n.tok.String()))
	}

	args := argsToFloat(n, e)

	for i := 1; i < len(args); i++ {
		if args[0] <= args[i] {
			return false
		}
	}

	return true
}

func handleNOT(n *node, e *environment) interface{} {
	var val interface{}

	return val
}

func handleADD(n *node, e *environment) interface{} {
	if n.childen == nil || len(n.childen) < 2 {
		panic(fmt.Sprintf("op: %s need 2 or more args", n.tok.String()))
	}

	var output float64
	args := argsToFloat(n, e)

	for i := 0; i < len(args); i++ {
		output += args[i]
	}

	return output
}

func handleSUB(n *node, e *environment) interface{} {
	var val float64
	if len(n.childen) <= 0 {
		panic(fmt.Sprintf("operator %s need 1 or more args", n.tok.String()))
	} else {
		args := argsToFloat(n, e)
		if len(args) == 1 {
			val = -args[0]
		} else {
			val = args[0]
			for i := 1; i < len(n.childen); i++ {
				val -= args[i]
			}
		}
	}
	return val
}

func handleMUL(n *node, e *environment) interface{} {
	var val float64
	if len(n.childen) < 2 {
		panic(fmt.Sprintf("operator %s need 2 or more args", n.tok.String()))
	} else {
		args := argsToFloat(n, e)
		val = args[0]
		for i := 1; i < len(n.childen); i++ {
			val *= args[i]
		}
	}
	return val
}

func handleQUO(n *node, e *environment) interface{} {
	var val float64
	if len(n.childen) < 2 {
		panic(fmt.Sprintf("operator %s need 2 or more args", n.tok.String()))
	} else {
		args := argsToFloat(n, e)
		val = args[0]
		for i := 1; i < len(n.childen); i++ {
			val /= args[i]
		}
	}
	return val
}

func handleVari(n *node, e *environment) interface{} {
	var val interface{}
	tmp := e.GetVar(n.tok.String())
	switch tmp.(type) {
	case int, float64, string:
		if n.executeType != SELFVAL {
			panic(callNotProcedureError(n))
		}
		val = tmp
	case *callStruct:
		if n.executeType != CALL {
			panic(callNotProcedureError(n))
		}
		call := tmp.(*callStruct)
		args := nodeChildren2Args(n, e)
		val = runCall(call, args)
	// case *node:
	// 	n := tmp.(*node)
	// 	val = n.value(e)
	default:
		panic(fmt.Sprintf("variable:%s type error", n.tok.String()))
	}
	return val
}

func handleDEFINE(n *node, e *environment) interface{} {
	if len(n.childen) < 2 {
		panic("define format error")
	}

	varName := n.childen[0].tok.String()

	if n.childen[0].childen == nil || len(n.childen[0].childen) == 0 {
		varVal := n.childen[1].value(e)
		e.CreateVar(varName, varVal)
		return nil
	}
	var val *callStruct

	if e.LookupVar(varName) {
		//重复定义检测
		panic(fmt.Sprintf("vari: %s has be defined repeat", varName))
	} else {
		val = newCall()
		val.name = varName
		for _, nc := range n.childen[0].childen {
			// args[i] = nc.tok.String()
			val.args = append(val.args, nc.tok.String())
		}

		var body = make([]*node, len(n.childen)-1)
		for i := 1; i < len(n.childen); i++ {
			body[i-1] = n.childen[i]
		}
		val.body = body
		// handleInnerDefine(body, val.env)
		extendEnv(val.env, e)
		e.CreateVar(varName, val)
	}
	return val
}

func handleLambda(n *node, e *environment) interface{} {
	val := newCall()
	val.name = "lambda"

	if len(n.childen) != 2 {
		panic("lambda syntax format error")
	}

	val.args = append(val.args, n.childen[0].tok.String())

	if len(n.childen[0].childen) > 0 {
		for _, nc := range n.childen[0].childen {
			// args[i] = nc.tok.String()
			val.args = append(val.args, nc.tok.String())
		}
	}

	val.body = []*node{n.childen[1]}
	// handleInnerDefine(body, val.env)
	extendEnv(val.env, e)
	e.CreateVar(val.name, val)

	return val
}

func handleSET(n *node, e *environment) {
	if len(n.childen) != 2 {
		panic("SET syntax format error")
	}

	varName := n.childen[0].tok.String()
	val := n.childen[1].value(e)

	e.SetVar(varName, val)
}

func handleIF(n *node, e *environment) interface{} {
	if len(n.childen) != 3 {
		panic(errIFArgsNum)
	}

	if n.childen[0].value(e).(bool) {
		return n.childen[1].value(e)
	}

	return n.childen[2].value(e)
}

func handleBEGIN(n *node, e *environment) interface{} {
	var output interface{}

	if len(n.childen) < 1 {
		panic("BEGIN syntax format error")
	}

	for _, child := range n.childen {
		output = child.value(e)
	}

	return output
}
