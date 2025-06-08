package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
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
	ok, err := matchLine(line, pattern)
	// ok, err := matchChars(string(line), pattern)
	fmt.Println("ok=", ok)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

/* My implementaion without library support
func matchChars(checkString string, reExp string) (bool, error) {
	l, r := 0, 1
	var nextExp string
	// fmt.Println(checkString, reExp)
	if len(checkString) != 0 && len(reExp) == 0 {
		return true, nil
	}
	if len(checkString) == 0 && len(reExp) != 0 {
		return false, nil
	}
	for r <= len(reExp) {
		ok := true
		nextExp = reExp[l:r]
		// println(l, r, nextExp)
		if nextExp == "\\" {
			r++
			continue
		} else if nextExp == "\\d" || nextExp == "\\w" {
			ok, _ = matchLine([]byte(checkString), nextExp)
			if !ok {
				return false, nil
			}
			chars := "01234556789"
			if nextExp == "\\w" {
				chars = "0120123456789abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ3456789"
			}
			index := bytes.IndexAny([]byte(checkString), chars)
			// println("index", index, chars)
			nextMatch, _ := matchChars(checkString[index+1:], reExp[r:])
			return ok && nextMatch, nil
		} else if nextExp == " " {
			if checkString[0:1] != " " {
				return false, nil
			}
			nextMatch, _ := matchChars(checkString[1:], reExp[r:])
			return ok && nextMatch, nil
		} else if matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, nextExp); matched {
			r++
			if r > len(reExp) {
				if !bytes.Equal([]byte(checkString[0:len(reExp)]), []byte(reExp)) {
					return false, nil
				} else {
					return true, nil
				}
			}
		} else {
			return matchLine([]byte(checkString), reExp)
		}
	}
	return true, nil
}

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
			toCheck := pattern[1 : len(pattern)-1]
			ok = bytes.ContainsAny(line, toCheck)
		} else if matched, _ := regexp.MatchString(`^\[\^[a-zA-Z]+\]$`, pattern); matched {
			endIndex := len(pattern) - 1
			res := 0
			searchLength := endIndex - 2
			for i := 2; i < endIndex; i++ {
				char := pattern[i]
				ok = bytes.ContainsAny(line, string(char))
				if ok {
					res += 1
				}
			}
			if res == searchLength {
				ok = false
			} else {
				ok = true
			}
		} else if matched, _ := regexp.MatchString(`^[a-zA-Z]$`, pattern); matched {
			//pattern cab be a single alphabet "a" or "A"
			// fmt.Println("matched alphabet")
			ok = bytes.ContainsAny(line, pattern)
		} else {
			return false, fmt.Errorf("unsupported pattern: %q", pattern)
		}
	}
	// fmt.Fprintln(os.Stdout, "matched status = ", ok)
	return ok, nil
}
*/

// Someones implementation with library that just looks so clean. This feels like cheating though.
func matchLine(line []byte, pattern string) (bool, error) {
	patternRuneCount := utf8.RuneCountInString(pattern)
	if patternRuneCount < 1 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool
	if validatePatternHasCharacterClasses(pattern) {
		ok = containsCharacterClass(string(line), pattern)
	} else {
		ok = bytes.ContainsAny(line, pattern)
	}

	return ok, nil
}

func validatePatternHasCharacterClasses(p string) bool {
	// [^a] or [ab]
	if containsCharacterClass(p, `^\[.*.\]$`) {
		return true
	}

	// ^staringString
	if containsCharacterClass(p, `^\^.+$`) {
		return true
	}

	// endingString$
	if containsCharacterClass(p, `^.+\$$`) {
		return true
	}

	// one or more someChar+
	if containsCharacterClass(p, `^.+\+.*$`) {
		return true
	}

	// number
	if strings.Contains(p, `\d`) {
		return true
	}

	// alphanums
	if strings.Contains(p, `\w`) {
		return true
	}

	return false
}

func containsCharacterClass(s string, p string) bool {
	return regexp.MustCompile(p).MatchString(s)
}
