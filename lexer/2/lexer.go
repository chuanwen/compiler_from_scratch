package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func IsWhite(x byte) bool {
	return x == ' ' || x == '\t'
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

func GetName(data[] byte) int {
	n := len(data)
	if n == 0 || !IsAlpha(data[0]) {
		return 0
	}
	i := 1
	for ; i < n; i++ {
		if !IsAlphaNum(data[i]) {
			break
		}
	}
	return i
}

func GetNumber(data[] byte) int {
	n := len(data)
	if n == 0 || !IsDigit(data[0]) {
		return 0
	}
	i := 1
	ndot := 0
	for ; i < n && ndot < 2; i++ {
		if !IsDigit(data[i]) && data[i] != '.'{
			break
		}
		if data[i] == '.' {
			ndot++;
		}
	}
	return i
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

func GetAssignOP(x, y, z byte) int {
	if x == '=' && IsBoundary(y) {
		return 1
	}
	if strings.IndexByte("+-*/", x) != -1 && y == '=' && IsBoundary(z) {
		return 2
	}
	return 0
}

func GetCmpOP(x, y, z byte) int {
	if (x == '>' || x == '<') && IsBoundary(y) {
		return 1
	}
	if (x == '>' || x == '<' || x == '!') && y == '=' && IsBoundary(z) {
		return 2
	}
	return 0
}

func GetMathOP(x, y, z byte) int {
	if strings.IndexByte("+-*/^", x) != -1 && IsBoundary(y){
		return 1
	}
	if x == y && (x == '<' || x == '>') && IsBoundary(z) {
		return 2
	}
	return 0
}

func GetOther(data []byte) int {
	ans := GetComments(data)
	if ans != 0 {
		return ans
	}
	ans = GetOP(data)
	if ans != 0 {
		return ans
	}
	return 1
}

func GetComments(data []byte) int {
	n := len(data)
	x := data[0]
	y := data[1]
	if x == '/' && y == '*' {
		for i := 2; i < (n-1); i++ {
			if data[i] == '*' && data[i+1] == '/' {
				return i + 2
			}
		}
		return -1
	}
	if x == '/' && y == '/' {
		for i := 2; i < n; i++ {
			if data[i] == '\n' {
				return i + 1
			}
		}
		return -1
	}
	return 0
}

func GetOP(data []byte) int {
	x := data[0]
	y := data[1]
	z := data[2]
	if k := GetAssignOP(x, y, z); k != 0 {
		return k
	}
	if k := GetCmpOP(x, y, z); k != 0 {
		return k
	}
	if k := GetMathOP(x, y, z); k != 0 {
		return k
	}
	return 0
}

func MySplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var i, j int
	n := len(data)
	for i = 0; i < n; i++ {
		if !IsWhite(data[i]) {
			break
		}
	}

	if i == n {
		return i, nil, nil
	}

	data = data[i:]
	look := data[0]

	switch {
	case IsAlpha(look):
		j = i + GetName(data)
	case IsDigit(look):
		j = i + GetNumber(data)
	case IsFuncSymbol(look):
		j = i + 1
	case IsNewLine(look):
		j = i + 1
	default:
		n := len(data)
		if n < 3 {
			return 0, nil, nil
		}
		j = GetOther(data)
		if j == -1 {
			return 0, nil, nil
		}
		j += i
	}
	return j, data[0:(j-i)], nil
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

	scanner := bufio.NewScanner(fd)
	scanner.Split(MySplit)
	for scanner.Scan() {
		token := scanner.Text()
		if token[len(token) - 1] != '\n' {
			fmt.Print(token, " | ")
		} else {
			fmt.Print(token)
		}
	}
}
