package main

import (
	"bufio"
	"fmt"
	"os"
)

func IsWhite(x byte) bool {
	return x == ' ' || x == '\t' || x == '\n'
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

	for j = i; j < n; j++ {
		if IsWhite(data[j]) {
			break
		}
	}
	return j, data[i:j], nil
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
		fmt.Println(scanner.Text())
	}
}
