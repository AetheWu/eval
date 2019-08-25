package eval

import (
	"io/ioutil"
	"strconv"
	"strings"
	"text/scanner"
)

type Lexer struct {
	CodeSource string
	Sentences  [][]Token
}

//StrScan scan string to symbol
type StrScan struct {
	scan  *scanner.Scanner
	token rune
}

//Next continue the scan
func (p *StrScan) Next() {
	p.token = p.scan.Scan()
}

//Text return the current symbol
func (p *StrScan) Text() string {
	return p.scan.TokenText()
}

func (p *StrScan) Pos() scanner.Position {
	return p.scan.Pos()
}

//IsEnd return true if scan ends
func (p *StrScan) IsEnd() bool {
	return p.token == scanner.EOF
}

//Init call to run string scan
func (p *StrScan) Init(input string) {
	if p.scan == nil {
		p.scan = new(scanner.Scanner)
	}
	p.scan.Init(strings.NewReader(input))
	// p.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	p.Next()
}

//IsInt return true if symbol is int type
func (p *StrScan) IsInt() bool {
	return p.token == scanner.Int
}

//IsFloat return true if symbol is float type
func (p *StrScan) IsFloat() bool {
	return p.token == scanner.Float
}

//IsString return true if symbol is string type
func (p *StrScan) IsString() bool {
	return p.token == scanner.String
}

//NewStrScan create a new StrScan
func NewStrScan() *StrScan {
	p := new(StrScan)
	return p
}

//NewToken trans str to symbol
func NewToken(scan *StrScan) Token {
	tok := new(Token)
	initOnce.Do(initTokenMap)

	if k, ok := tokenMap[scan.Text()]; ok {
		tok.Kind = k
		tok.Str = scan.Text()
	} else {
		if scan.IsInt() {
			val, _ := strconv.Atoi(scan.Text())
			tok = createToken(scan.Text(), INT, val)
		} else if scan.IsFloat() {
			val, _ := strconv.ParseFloat(scan.Text(), 64)
			tok = createToken(scan.Text(), FLOAT, val)
		} else if scan.IsString() {
			tok = createToken(scan.Text(), STRING, scan.Text())
		} else {
			tok.Kind = IDENT
			tok.Str = scan.Text()
		}
	}
	tok.Pos = scan.Pos()
	return *tok
}

func createToken(s string, k kind, val interface{}) *Token {
	tok := new(Token)

	tok.Kind = k
	tok.Str = s
	tok.Val = val
	return tok
}

func CodeSplit(input string) []string {
	var tmpStr, output []string
	BraceStack := NewStack()

	for i := 0; i < len(input); i++ {
		if input[i] == ';' {
			for i < len(input) && input[i] != '\n' && input[i] != '\r' {
				tmpStr = append(tmpStr, string(input[i]))
				i++
			}
			output = append(output, strings.Join(tmpStr, ""))
			tmpStr = []string{}
			continue
		}

		if !BraceStack.IsEmpty() {
			tmpStr = append(tmpStr, string(input[i]))
			if input[i] == ')' {
				BraceStack.Pop()
				if BraceStack.IsEmpty() {
					output = append(output, strings.Join(tmpStr, ""))
					tmpStr = []string{}
					continue
				}
			}
		} else {
			if input[i] == ' ' || input[i] == '\r' || input[i] == '\n' || input[i] == '\t' {
				if len(tmpStr) > 0 {
					if input[i] == '\r' || input[i] == '\n' || input[i] == ' ' {
						output = append(output, strings.Join(tmpStr, ""))
						tmpStr = []string{}
					}
				}
				continue
			}

			tmpStr = append(tmpStr, string(input[i]))

			if i == len(input)-1 && len(tmpStr) > 0 {
				output = append(output, strings.Join(tmpStr, ""))
			}
		}

		if input[i] == '(' {
			BraceStack.Push(input[i])
		}

	}
	if !BraceStack.IsEmpty() {
		panic("lack of right paren")
	}

	return output
}

func buildSentence(expStr string) []Token {
	scan := NewStrScan()
	var sentence []Token
	for scan.Init(expStr); !scan.IsEnd(); scan.Next() {
		sentence = append(sentence, NewToken(scan))
	}
	return sentence
}

//FileLoad trans source code to string
func FileLoad(filename string) string {
	sourceByte, err := ioutil.ReadFile(filename)

	if err != nil {
		panic("source file open failed")
	}

	return string(sourceByte)
}