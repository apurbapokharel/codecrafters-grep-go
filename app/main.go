package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	myast "github.com/codecrafters-io/grep-starter-go/app/AST"
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

	// echo -n "log" | ./your_program.sh -E "^log"
	if len(os.Args) == 3 {
		pattern := os.Args[2]
		line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
			os.Exit(2)
		}
		// Build the pattern into a ParseTree
		regExpParser := myast.NewParser([]rune(pattern))
		node := regExpParser.Parse0()
		node.Log()
		// Once the tree is built check the ParsedPattern against the checkString
		checkStringParser := myast.NewParser([]rune(string(line)))
		ok, err := checkStringParser.CheckParseTree(node)
		fmt.Println("result", ok)

		// ok, err := matchLine2(line, pattern)
		// ok, err := matchChars(string(line), pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		if !ok {
			os.Exit(1)
		}
		// ./your_program.sh -E "carrot" fruits.txt
		// ./your_program.sh -E "search_pattern" file1.txt file2.txt
	} else if len(os.Args) >= 4 {
		length := len(os.Args) - 3
		fileArray := make([]string, 0, length)

		for i := range length {
			fileArray = append(fileArray, os.Args[3+i])
		}
		for _, file := range fileArray {
			if file == " " {
				continue
			}
			file, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			var lines []string
			scanner := bufio.NewScanner(file)

			// Read each line and append to the slice
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}

			noMatch := true
			for _, line := range lines {
				// parse the pattern to a ParseTree
				regExpParser := myast.NewParser([]rune(os.Args[2]))

				node := regExpParser.Parse0()
				// node.Log()
				// check the presence of the pattern inside the checkString
				checkStringParser := myast.NewParser([]rune(line))
				ok, err := checkStringParser.CheckParseTree(node)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(2)
				}
				if ok {
					fmt.Println(line)
					noMatch = false
				}
			}
			if noMatch {
				os.Exit(1)
			}
		}
	}
	// default exit code is 0 which means success
}

func matchChars(checkString string, reExp string) (bool, error) {
	return true, nil
}

// My implementaion without library
// support supports determinstic cases only
/*
func matchChars(checkString string, reExp string) (bool, error) {
	l, r := 0, 1
	var nextExp string
	// if all reExp match done than return true
	if len(checkString) != 0 && len(reExp) == 0 {
		return true, nil
	}
	// run out of checkString to match with reExp
	if len(checkString) == 0 && len(reExp) != 0 {
		return false, nil
	}
	// check if "|" exits
	if index := bytes.IndexAny([]byte(reExp), "|"); index != -1 {
		startParen := bytes.IndexAny([]byte(reExp), "(")
		endParen := bytes.IndexAny([]byte(reExp), ")")
		patternString := reExp[startParen+1 : endParen]
		remString := reExp[endParen+1:]
		stringsToCheck := bytes.Split([]byte(patternString), []byte("|"))
		res := false
		for i := 0; i < len(stringsToCheck); i++ {
			reExpToCheck := reExp[:startParen] + string(stringsToCheck[i]) + remString
			println(checkString, reExpToCheck)
			// funRes, _ := matchChars(checkString, reExpToCheck)
			res = res || bytes.Contains([]byte(checkString), []byte(reExpToCheck))
		}
		return res, nil
	}
	// check start of string
	if string(reExp[0]) == "^" {
		return strings.HasPrefix(checkString, reExp[1:]), nil
	}
	// check end of string
	if string(reExp[len(reExp)-1]) == "$" {
		return strings.HasSuffix(checkString, reExp[:len(reExp)-1]), nil
	}
	// check for one or more times match // refactor to support wildcard
	if index := bytes.IndexAny([]byte(reExp), "+"); index != -1 {
		skipChar := reExp[index-1]

		lenStringBefore := len(reExp[:index-1])
		preIndexMatch, _ := matchChars(checkString[:lenStringBefore+1], reExp[:index-1])
		if index == len(reExp)-1 || !preIndexMatch {
			return preIndexMatch, nil
		}
		stopChar := reExp[index+1]
		i := lenStringBefore
		var res bool
		// check the skip char
		// if skip char is . skip until stop char is encountered
		// else skip the skip char
		// break if the checkString has nothing after the skip char
		for true {
			if string(skipChar) == "." || checkString[i] == skipChar {
				i++
				// j++
			} else if checkString[i] != skipChar {
				res = false
				break
			}

			if i == len(checkString) || checkString[i] == stopChar {
				res = true
				break
			}
		}
		if !res {
			return false, nil
		}
		postIndexMatch := bytes.Contains([]byte(checkString[i:]), []byte(reExp[index+1:]))
		return postIndexMatch, nil

	}
	// check for zero or more times
	// does not support aczzct and c?at
	if index := bytes.IndexAny([]byte(reExp), "?"); index != -1 {
		checkChar := reExp[index-1]
		var checkCharAfter byte
		if index+1 < len(reExp) {
			checkCharAfter = reExp[index+1]
		}

		// 1.check from start to ___*__(checkChar)?___*___ char from regExp subset of checkstring

		// a.check if the start index match
		startChar := reExp[0]
		checkIndex := bytes.IndexAny([]byte(checkString), string(startChar))
		if checkIndex == -1 {
			return false, nil
		}
		// b.once index match check if the rem substring until checkChar match
		var preIndexMatch bool
		remLength := len(reExp) - len(reExp[index-1:])
		preIndexMatch, _ = matchChars(checkString[checkIndex:checkIndex+remLength], reExp[:index-1])
		if !preIndexMatch {
			return false, nil
		}
		// 2.check if the checkchar match or match the postSubString
		checkIndex = checkIndex + remLength
		var postIndexMatch bool
		if checkString[checkIndex] == checkChar || checkString[checkIndex] == checkCharAfter {
			remLength := len(reExp[index+1:])
			postIndexMatch = false
			if remLength == 0 {
				postIndexMatch = true
			} else {
				checkIndex = len(checkString) - remLength
				postIndexMatch, _ = matchChars(checkString[checkIndex:], reExp[index+1:])
			}
		} else {
			return false, nil
		}
		if !postIndexMatch {
			return false, nil
		}
		return true, nil
	}
	// WILDCARD: match 1 times
	if index := bytes.IndexAny([]byte(reExp), "."); index != -1 {
		lenPre := len(reExp) - len(reExp[index:])
		preSubArray := bytes.Equal([]byte(checkString[:lenPre]), []byte(reExp[:index]))
		postSubArray := bytes.Contains([]byte(checkString[lenPre+1:]), []byte(reExp[index+1:]))
		return preSubArray && postSubArray, nil
	}
	// Remaining cases
	for r <= len(reExp) {
		ok := true
		nextExp = reExp[l:r]
		println(l, r, nextExp)
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
					// if !bytes.Contains([]byte(checkString), []byte(reExp)) {
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

// Someones implementation that i extended with library that just looks so clean but is cheating though.
/*
func matchLine2(line []byte, pattern string) (bool, error) {
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
	// .
	if containsCharacterClass(p, `^.*\..*$`) {
		return true
	}

	// ?
	if containsCharacterClass(p, `^.*.?.*$`) {
		return true
	}

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
*/
