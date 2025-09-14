package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	myast "github.com/codecrafters-io/grep-starter-go/app/AST"
	mytree "github.com/codecrafters-io/grep-starter-go/app/PathTree"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
// os.Args = [/tmp/codecrafters-build-grep-go -E a]

func main() {
	if len(os.Args) < 3 {
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
	} else if len(os.Args) >= 4 {
		if os.Args[1] == "-r" {
			// ./your_program.sh -r -E ".*er" dir/
			pathTree, _ := BuildFileTree(os.Args[4])
			regExpParser := myast.NewParser([]rune(os.Args[3]))
			node := regExpParser.Parse0()
			status := TraverseTreeAndCheck(pathTree, node)
			if status {
				os.Exit(1)
			}
		} else {
			// ./your_program.sh -E "carrot" fruits.txt
			// ./your_program.sh -E "search_pattern" file1.txt file2.txt
			length := len(os.Args) - 3
			fileArray := make([]string, 0, length)

			for i := range length {
				fileArray = append(fileArray, os.Args[3+i])
			}
			noMatch := true
			for _, fileName := range fileArray {
				file, err := os.Open(fileName)
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
						if length > 1 {
							fmt.Printf("%s:%s\n", fileName, line)
						} else {
							fmt.Println(line)
						}
						noMatch = false
					}
				}
			}
			if noMatch {
				os.Exit(1)
			}
		}
	}
	// default exit code is 0 which means success
}

func BuildFileTree(path string) (*mytree.FileTree, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	root := mytree.NewFileTree(path, true)
	for _, file := range files {
		if file.IsDir() {
			// directoryNode := mytree.NewFileTree(path + file.Name(), true)
			children, _ := BuildFileTree(path + file.Name() + "/")
			root.Children = append(root.Children, children)
		} else {
			leafNode := mytree.NewFileTree(path+file.Name(), false)
			root.Children = append(root.Children, leafNode)
		}
		// fmt.Println(file.Name())
	}
	return root, nil
}

func TraverseTreeAndCheck(root *mytree.FileTree, astNode myast.RegexpNode) bool {
	// A queue to hold the nodes to be visited
	queue := []*mytree.FileTree{root}
	noMatch := true
	for len(queue) > 0 {
		// Dequeue the first node
		node := queue[0]
		queue = queue[1:]

		// fmt.Printf("Visiting: %s (Is directory: %t)\n", node.Pwd, node.IsDirectory)
		if !node.IsDirectory {
			// Enqueue all children of the current node
			// need to check if the current file has the checkstring
			file, err := os.Open(node.Pwd)
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

			for _, line := range lines {
				// parse the pattern to a ParseTree
				// check the presence of the pattern inside the checkString
				checkStringParser := myast.NewParser([]rune(line))
				ok, err := checkStringParser.CheckParseTree(astNode)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(2)
				}
				if ok {
					fmt.Printf("%s:%s\n", node.Pwd, line)
					noMatch = false
				}
			}
		}
		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}
	return noMatch
}
