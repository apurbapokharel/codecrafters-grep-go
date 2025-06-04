package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	// "reflect"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
// os.Args = [/tmp/codecrafters-build-grep-go -E a]
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	// fmt.Println("type =", reflect.TypeOf(line))
	// fmt.Println("ascii=", line)
	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

// func matchLine(line []byte, pattern string) (bool, error) {
// 	if utf8.RuneCountInString(strings.Trim(pattern, "\\")) != 1 {
// 		return false, fmt.Errorf("unsupported pattern: %q", pattern)
// 	}

// 	var ok bool

// 	// You can use print statements as follows for debugging, they'll be visible when running tests.
// 	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

// 	pattern = strings.ReplaceAll(pattern, "\\d", "0123456789")
// 	pattern = strings.ReplaceAll(pattern, "\\w", "0123456789abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ")
// 	// fmt.Println("pattern =", pattern)
// 	ok = bytes.ContainsAny(line, pattern)

// 	fmt.Println("ok =", ok)
// 	return ok, nil
// }

func matchLine(line []byte, pattern string) (bool, error) {
	// fmt.Fprintln(os.Stdout, "pattern = ", pattern)
	var ok bool
	switch pattern {
	case "\\d":
		pattern = strings.ReplaceAll(pattern, "\\d", "0123456789")
		ok = bytes.ContainsAny(line, pattern)
	case "\\w":
		pattern = strings.ReplaceAll(pattern, "\\w", "0120123456789abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ3456789")
		ok = bytes.ContainsAny(line, pattern)
	default:
		if matched, _ := regexp.MatchString(`^\[[a-zA-Z]+\]$`, pattern); matched {
			//pattern can be [__*__]
			startIndex := strings.Index(pattern, "[")
			endIndex := strings.Index(pattern, "]")
			for i := startIndex + 1; i < endIndex; i++ {
				char := pattern[i]
				fmt.Println("char", string(char))
				ok = bytes.ContainsAny(line, string(char))
				if ok {
					break
				}
			}
		} else if matched, _ := regexp.MatchString(`^[a-zA-Z]$`, pattern); matched {
			//pattern cab be a single alphabet "a" or "A"
			fmt.Println("matched alphabet")
			ok = bytes.ContainsAny(line, pattern)
		} else {
			return false, fmt.Errorf("unsupported pattern: %q", pattern)
		}
	}
	fmt.Fprintln(os.Stdout, "matched status = ", ok)
	return ok, nil
}
