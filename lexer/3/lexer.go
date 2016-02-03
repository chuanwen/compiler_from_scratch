package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	tokVariable int = iota
	tokNumber
	tokString
	tokAssignment // one of =  +=  -=  *=  /=
	tokMathOp     // one of the math OP: + - * / % ^ <<  >>
	tokCmpOP      // one of the comparison OP: > < >= <= != ==
	tokBooleanOP  // && || !
	tokKeyword    // one of the keywords: if else for break ...
	tokFuncSymbol // one of special symbols: ( ) [  ] {  } ?
	tokNewLine
	tokComment
	tokEOF
	tokOther
)

type Token struct {
	typ int
	val string
}

func (t *Token) String() string {
	var typeNames = [...]string{
		"Variable",
		"Number",
		"String",
		"Assign",
		"MathOp",     // one of the math OP: + - * / % ^ <<  >>
		"CmpOP",      // one of the comparison OP: > < >= <= != ==
		"BooleanOP",  // && || !
		"Keyword",    // one of the keywords: if else for break ...
		"FuncSymbol", // one of special symbols: ( ) [  ] {  } ?
		"NewLine",
		"Comment",
		"EOF",
		"Other",
	}
	typeName := typeNames[t.typ]
	if t.val[len(t.val)-1] != '\n' {
		return fmt.Sprintf("%s %s | ", typeName, t.val)
	} else {
		return fmt.Sprintf("%s %s", typeName, t.val)
	}
}

var KEYWORDS = map[string]bool{
	"if":       true,
	"else":     true,
	"for":      true,
	"break":    true,
	"function": true,
	"return":   true,
	"struct":   true,
	"int":      true,
	"float":    true,
	"char":     true,
	"typedef":  true,
	"type":     true,
	"auto":     true,
}

func IsWhite(x byte) bool {
	return x == ' ' || x == '\t' || x == '\r'
}

func IsAlpha(x byte) bool {
	return (x >= 'a' && x <= 'z') || (x >= 'A' && x <= 'Z')
}

func IsDigit(x byte) bool {
	return x >= '0' && x <= '9'
}

func IsAlphaNum(x byte) bool {
	return IsAlpha(x) || IsDigit(x)
}

func GetName(data []byte) (advance int, tokenType int) {
	n := len(data)
	advance = 1
	for ; advance < n; advance++ {
		if !IsAlphaNum(data[advance]) {
			break
		}
	}
	tokenType = tokVariable
	if KEYWORDS[string(data[:advance])] {
		tokenType = tokKeyword
	}
	return
}

func GetNumber(data []byte) (advance int, tokenType int) {
	n := len(data)
	advance = 1
	tokenType = tokNumber
	ndot := 0
	for ; advance < n && ndot < 2; advance++ {
		if !IsDigit(data[advance]) && data[advance] != '.' {
			break
		}
		if data[advance] == '.' {
			ndot++
		}
	}
	return
}

func IsNewLine(x byte) bool {
	return x == '\n'
}

func IsFuncSymbol(x byte) bool {
	return strings.IndexByte("()[]{}?.:;", x) != -1
}

func IsBoundary(x byte) bool {
	return IsWhite(x) || IsAlphaNum(x) || x == '(' || IsNewLine(x) || x == ';'
}

func IsAssignOP(x, y, z byte) int {
	if x == '=' && IsBoundary(y) {
		return 1
	}
	if strings.IndexByte("+-*/", x) != -1 && y == '=' && IsBoundary(z) {
		return 2
	}
	return 0
}

func IsCmpOP(x, y, z byte) int {
	if (x == '>' || x == '<') && IsBoundary(y) {
		return 1
	}
	if (x == '>' || x == '<' || x == '!') && y == '=' && IsBoundary(z) {
		return 2
	}
	return 0
}

func IsMathOP(x, y, z byte) int {
	if strings.IndexByte("+-*/^", x) != -1 && IsBoundary(y) {
		return 1
	}
	if x == y && (x == '<' || x == '>') && IsBoundary(z) {
		return 2
	}
	return 0
}

func GetOther(data []byte) (advance int, tokenType int) {
	advance, tokenType = GetComments(data)
	if advance != 0 {
		return
	}
	advance, tokenType = GetOP(data)
	if advance != 0 {
		return
	}
	return 1, tokOther
}

func GetComments(data []byte) (advance int, tokenType int) {
	n := len(data)
	x := data[0]
	y := data[1]
	if x == '/' && y == '*' {
		for i := 2; i < (n - 1); i++ {
			if data[i] == '*' && data[i+1] == '/' {
				return i + 2, tokComment
			}
		}
		return -1, tokComment
	}
	if x == '/' && y == '/' {
		for i := 2; i < n; i++ {
			if data[i] == '\n' {
				return i + 1, tokComment
			}
		}
		return -1, tokComment
	}
	return 0, 0
}

func GetOP(data []byte) (advance int, tokenType int) {
	x := data[0]
	y := data[1]
	z := data[2]
	if i := IsAssignOP(x, y, z); i != 0 {
		return i, tokAssignment
	}
	if i := IsCmpOP(x, y, z); i != 0 {
		return i, tokCmpOP
	}
	if i := IsMathOP(x, y, z); i != 0 {
		return i, tokMathOp
	}
	return 0, 0
}

func MySplit(data []byte, atEOF bool) (advance int, tokenType int, token []byte, err error) {
	var i, j int
	n := len(data)
	for i = 0; i < n; i++ {
		if !IsWhite(data[i]) {
			break
		}
	}

	if i == n {
		return i, 0, nil, nil
	}

	data = data[i:]
	look := data[0]

	switch {
	case IsAlpha(look):
		j, tokenType = GetName(data)
	case IsDigit(look):
		j, tokenType = GetNumber(data)
	case IsFuncSymbol(look):
		j = 1
		tokenType = tokFuncSymbol
	case IsNewLine(look):
		j = 1
		tokenType = tokNewLine
	default:
		n := len(data)
		if n < 3 {
			return 0, 0, nil, nil
		}
		j, tokenType = GetOther(data)
		if j == -1 {
			return 0, tokenType, nil, nil
		}
	}
	advance = i + j
	return advance, tokenType, data[0:j], nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("No input source file")
		return
	}

	sourceFile := os.Args[1]
	fd, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := NewScanner(fd)
	scanner.Split(MySplit)
	for scanner.Scan() {
		token := scanner.Token()
		fmt.Print(token)
	}
}
