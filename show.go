package eval

import (
	"fmt"
	"reflect"
	"strconv"
)

func printSentence(sentence []Token) {
	for _, tok := range sentence {
		fmt.Printf("%s ", tok.String())
	}
	fmt.Printf("\n")
}

func printSingleNode(n *node, s int) {
	space := " "
	for i := 0; i < s; i++ {
		space += " "
	}
	fmt.Printf("%s %s, nodeType:%d \n", space, n.tok.String(), n.executeType)
	if n.childen == nil {
		return
	}

	for _, child := range n.childen {
		printSingleNode(child, s+1)
	}
}

func formatValue(x interface{}) string {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Invalid:
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return v.String()
		// return strconv.Quote(v.String())
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return ""
	}
}

func evalDisplay(x interface{}) {
	str := formatValue(x)
	if str != "" {
		fmt.Println(str)
	}
}

func PrintNode(filename string) {
	input := FileLoad(filename)
	sentences := CodeSplit(input)

	for i := 0; i < len(sentences); i++ {
		sentence := buildSentence(sentences[i])
		node := ast(sentence)
		printSingleNode(node, 0)
	}
}

func PrintSen(filename string) {
	input := FileLoad(filename)
	sentences := CodeSplit(input)

	for i := 0; i < len(sentences); i++ {
		sentence := buildSentence(sentences[i])
		fmt.Printf("%d: ", i)
		printSentence(sentence)
	}
}
