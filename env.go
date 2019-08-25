package eval

import "fmt"

type environment struct {
	valMap  map[string]interface{}
	baseEnv *environment
}

type callStruct struct {
	name string
	args []string
	body []*node
	env  *environment
}

func newCall() *callStruct {
	call := new(callStruct)
	call.args = []string{}
	call.body = []*node{}
	call.env = initEnv()
	return call
}

func initEnv() *environment {
	env := new(environment)
	env.valMap = make(map[string]interface{})
	return env
}

func (e *environment) LookupVar(name string) bool {
	if _, ok := e.valMap[name]; ok {
		return true
	}
	return false
}

func (e *environment) GetVar(name string) interface{} {
	var output interface{}
	if item, ok := e.valMap[name]; ok {
		output = item
	} else {
		if e.baseEnv != nil {
			output = e.baseEnv.GetVar(name)
		} else {
			panic(fmt.Sprintf("variable: %s is not defined", name))
		}
	}
	return output
}

func (e *environment) SetVar(name string, val interface{}) {
	if _, ok := e.valMap[name]; ok {
		e.valMap[name] = val
	} else {
		if e.baseEnv != nil {
			e.baseEnv.SetVar(name, val)
		} else {
			panic(fmt.Sprintf("variable:%s is not defined", name))
		}
	}
}

func (e *environment) CreateVar(name string, vari interface{}) {

	if _, ok := e.valMap[name]; ok {
		panic(fmt.Sprintf("define:%s has repeated", name))
	} else {
		e.valMap[name] = vari
	}
}

func extendEnv(newEnv, baseEnv *environment) {
	newEnv.baseEnv = baseEnv
}

func destroyEnv(e *environment) {
	e.valMap = make(map[string]interface{})
}

func copyCall(call *callStruct) *callStruct {
	output := newCall()
	output.name = call.name
	output.args = call.args
	extendEnv(output.env, call.env.baseEnv)
	output.body = call.body

	return output
}
