package myast

import (
	"fmt"
	"os"
	"regexp"
)

// used for both building AST with regExp
// as well as for checkingParseTree with checkString
type Parser struct {
	i       int
	pattern []rune
	context *ParserContext
}

func NewParser(pattern []rune) *Parser {
	return &Parser{0, pattern, NewParserContext()}
}

func (p Parser) getCurrent() (string, bool) {
	if p.i >= len(p.pattern) {
		return "", true
	}
	return string(p.pattern[p.i]), false
}

func (p *Parser) advance() {
	p.i++
}

func (p *Parser) currentIndex() int {
	return p.i
}

func (p *Parser) moveIndex(pos int) {
	p.i = pos
}

func (p *Parser) Parse0() RegexpNode {
	node := p.Parse1()
	current, isEnd := p.getCurrent()
	if isEnd {
		return node
	}
	if current == "|" {
		p.advance()
		return Alternate{node, p.Parse0()}
	}
	return node
}

func (p *Parser) Parse1() RegexpNode {
	node := p.Parse2()
	current, isEnd := p.getCurrent()
	if isEnd {
		return node
	}
	// println("p1", current)
	matched, _ := regexp.MatchString(`^[a-zA-Z\s]$`, current)

	if current == "]" || current == ")" {
		return node
	}

	if matched {
		if p.context.isAlternate {
			return Alternate{node, p.Parse1()}
		}
		return Concat{node, p.Parse1()}
	}

	if current == "\\" || current == "(" || current == "^" || current == "$" || current == "." {
		return Concat{node, p.Parse1()}
	}

	return node
}

func (p *Parser) Parse2() RegexpNode {
	node := p.Parse3()
	if _, ok := node.(Undefined); ok {
		println("Undefined at pos", p.currentIndex())
		return Undefined{}
	}
	current, isEnd := p.getCurrent()
	if isEnd {
		return node
	}
	if current == "?" {
		p.advance()
		return Optional{node}
	}
	if current == "+" {
		p.advance()
		return Repeat{node}
	}
	return node
}

func (p *Parser) Parse3() RegexpNode {
	current, isEnd := p.getCurrent()
	// println(current, p.currentIndex())
	if isEnd {
		return nil
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z\s,]$`, current); matched {
		node := Literal{[]rune(current)[0]}
		p.advance()
		return node
	}
	if current == "." {
		p.advance()
		return Wildcard{}
	}
	if current == "\\" {
		p.advance()
		current, _ := p.getCurrent()
		p.advance()
		if current == "d" {
			node := Digit{}
			return node
		} else {
			node := AlphaNum{}
			return node
		}
	}
	if current == "[" {
		p.advance()
		currentPosNeg, _ := p.getCurrent()
		if currentPosNeg == "^" {
			p.advance()
		}
		p.context.isAlternate = true
		node := p.Parse0()
		p.context.isAlternate = false
		current, _ := p.getCurrent()
		if current == "]" {
			var node2 RegexpNode
			if currentPosNeg == "^" {
				node2 = NegativeCharacterGroup{node}
			} else {
				node2 = PositiveCharacterGroup{node}
			}
			p.advance()
			return node2
		}
	}
	if current == "(" {
		p.advance()
		// p.context.isAlternate = true
		node := p.Parse0()
		// p.context.isAlternate = false
		current, _ := p.getCurrent()
		if current == ")" {
			p.advance()
			return node
		}
	}
	if current == "^" {
		p.advance()
		return AnchorStart{}
	}
	if current == "$" {
		p.advance()
		return AnchorEnd{}
	}

	// if current == ""
	return Undefined{}
}

func (p *Parser) CheckParseTree(node RegexpNode) (bool, error) {
	current, isEnd := p.getCurrent()
	// println(current, isEnd, node.get())
	//end of checkstring
	if isEnd {
		if _, ok := node.(AnchorEnd); ok {
			if p.i == len(p.pattern) {
				return true, nil
			}
		}

		//For handle dog with dogs?
		//Skip the last part
		if _, ok := node.(Optional); ok {
			return true, nil
		}
		return false, nil
	}

	// When the index matches the first left child it will parse the rest of the right tree
	// If true return else check till checkstring is end
	// We need contigous match so skipChar is set to false once a first match (this is done in match(literals,nums, alphanums) is found
	if c, ok := node.(Concat); ok {
		// println("Inside Concat, currentIndex = ", p.currentIndex(), c.get(), p.context.skipChars)
		firstLeftChild, err := getFirstLeftChild(node)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		var secondLeftChild string
		if firstLeftChild == "AnchorStart()" {
			secondLeftChild, err = getFirstLeftChild(c.rightNode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}

		// Greedy parsing with Backtracking for fallback
		if firstLeftChild == "Repeat()" {
			// println("Repeat Concat, index = ", p.currentIndex(), c.get())
			prevIndex := p.currentIndex()
			leftParse, _ := p.CheckParseTree(c.leftNode)
			// Repeat + has is one or more so index has to go up by 1 at the least
			if prevIndex == p.currentIndex() || !leftParse {
				return false, nil
			}
			// p.context.skipChars = false
			rightParse, _ := p.CheckParseTree(c.rightNode)
			// if we can parse the whole string while greedy parsing the repeat return true
			if rightParse {
				return true, nil
			}
			// if greedy repeat does not work try backtracking
			if !rightParse {
				for i := p.currentIndex() - 1; i > prevIndex; i-- {
					p.moveIndex(i)
					rightParse, _ = p.CheckParseTree(c.rightNode)
					if rightParse {
						return true, nil
					}

				}
			}
			// if backtracking does not work return false
			return false, nil
		}

		// Greedy parsing with Backtracking for fallback
		if firstLeftChild == "Optional()" {
			// println("Optional Concat, currentIndex = ", p.currentIndex(), c.get())
			prevIndex := p.currentIndex()

			// check optional first i.e skip lefttree parse
			rightParse, _ := p.CheckParseTree(c.rightNode)
			// if we can parse the whole string while greedy parsing the repeat return true
			if rightParse {
				return true, nil
			} else {
				// println("Skipping not mathcing OPtional")
				p.moveIndex(prevIndex)
			}

			//treat the ? like a + now (1 or more time)
			// Repeat + has is one or more so index has to go up by 1 at the least
			leftParse, _ := p.CheckParseTree(c.leftNode)
			// p.context.skipChars = false
			if prevIndex == p.currentIndex() || !leftParse {
				return false, nil
			}
			rightParse, _ = p.CheckParseTree(c.rightNode)
			// if we can parse the whole string while greedy parsing the repeat return true
			if rightParse {
				return true, nil
			}
			// if greedy repeat does not work try backtracking
			if !rightParse {
				for i := p.currentIndex() - 1; i > prevIndex; i-- {
					p.moveIndex(i)
					rightParse, _ = p.CheckParseTree(c.rightNode)
					if rightParse {
						return true, nil
					}

				}
			}
			// if backtracking does not work return false
			return false, nil
		}

		if firstLeftChild == "AnchorStart()" {
			matched, _ := regexp.MatchString(secondLeftChild, string(p.pattern[0]))
			if !matched {
				return false, nil
			}
			// p.context.skipChars = false
			success := p.checkRightSubtree(c.rightNode)
			if success {
				return true, nil
			}
		}

		if firstLeftChild == "Alternate()" {
			// println("Alternate Concat, index = ", p.currentIndex())
			leftMatch, _ := p.CheckParseTree(c.leftNode)
			rightMatch, _ := p.CheckParseTree(c.rightNode)
			return leftMatch && rightMatch, nil
		}

		if firstLeftChild == "Wildcard()" {
			// println("Alternate Concat, index = ", p.currentIndex())
			// skip check as it can go with anything
			p.advance()
			p.context.skipChars = false
			rightMatch, _ := p.CheckParseTree(c.rightNode)
			return rightMatch, nil
		}

		// println("FLC", firstLeftChild)
		var status bool = false
		//FIXME: Using loop here is not the right thing i think
		for i := p.i; i < len(p.pattern); i++ {
			// println("Retrying", p.context.skipChars, firstLeftChild)
			if matched, _ := regexp.MatchString(firstLeftChild, string(p.pattern[i])); matched {
				p.moveIndex(i + 1)
				p.context.skipChars = false
				// _, ok := c.rightNode.(Optional)
				// if ok {
				// 	//todo: skip this check
				// 	println("RN OPT")
				// }
				success := p.checkRightSubtree(c.rightNode)
				if success {
					status = true
					break
				}
				// p.context.skipChars = true
				p.moveIndex(i)
				// break
			}

			// NOTE: honestly, this is the most complicated part. Figuring out this loop
			// When i am 1 2 depth inside a recusrion (the skip chars will alreay be false),
			// and if there is no match with the current ith character we dont wanna test other we wanna return false
			// meaning false the 1st ever check befire all the recursion is false so check again
			// echo -n "I am I see 1 cat, 2 dogs and 3 cows" | ./your_program.sh -E "I see (\d (cat|dog|cow)s?(, | and )?)+$" //true
			// I will match but a will not so return false and then check again (in the main loop) thats what this means

			// if we are searching inside a matched case we cannot skip chars
			// if !p.context.skipChars {
			// 	return false, nil
			// }
		}
		return status, nil

	}

	if c, ok := node.(Repeat); ok {
		// println("Repeat Self, index = ", p.currentIndex(), c.get())
		var res bool = false
		for {
			matched, _ := p.CheckParseTree(c.node)
			res = res || matched
			if !matched {
				break
			}
		}
		// println("Repeat Self out, curentIndex = ", p.currentIndex())
		return res, nil
	}

	if _, ok := node.(Wildcard); ok {
		p.advance()
		p.context.skipChars = false
		return true, nil
	}

	//FIXME: this does not work echo -n "a" | ./your_program.sh -E "s?"
	if c, ok := node.(Optional); ok {
		// println("Optinal Self, index = ", p.currentIndex(), c.get())
		prevIndex := p.currentIndex()
		var res bool = false
		// check for 1 or more
		for {
			// println("Running OPtional loop")
			matched, _ := p.CheckParseTree(c.node)
			res = res || matched
			if !matched {
				break
			}
		}
		//if 1 or more false skipped this exp
		if !res {
			// println("skip optional")
			p.moveIndex(prevIndex)
			return false, nil
		}
		//else return true result after parsing optional
		// println("Optional Self out, curentIndex = ", p.currentIndex())
		return res, nil
	}

	if c, ok := node.(Alternate); ok {
		// println("Alternate Self, curentIndex = ", p.currentIndex())
		positionBefore := p.currentIndex()
		leftTreeParse, _ := p.CheckParseTree(c.leftNode)
		if leftTreeParse {
			// println("Alternate Self OUT1, curentIndex = ", p.currentIndex())
			return true, nil
		}
		p.moveIndex(positionBefore)
		rightTreeParse, _ := p.CheckParseTree(c.rightNode)
		// println("Alternate Self OUT2, curentIndex = ", p.currentIndex())
		return rightTreeParse, nil
	}

	if c, ok := node.(PositiveCharacterGroup); ok {
		result, _ := p.CheckParseTree(c.node)
		return result, nil
	}

	//TODO: This is not a good solution i convert to AST and back to string but this was annoying and i wanted to be done with this
	if c, ok := node.(NegativeCharacterGroup); ok {
		dict := make(map[string]bool)
		regularExp := getExpFromTree(c.node)
		for _, s := range regularExp {
			dict[s] = true
		}
		var status bool = false
		for i := 0; i < len(p.pattern); i++ {

			if !dict[string(p.pattern[i])] {
				status = true
				break
			}
		}
		return status, nil
	}

	if c, ok := node.(Literal); ok {
		if current == string(c.value) {
			p.context.skipChars = false
			p.advance()
			return true, nil
		} else {
			if p.context.skipChars {
				p.advance()
				result, _ := p.CheckParseTree(node)
				res := false || result
				if res {
					p.context.skipChars = false
				}
				return res, nil
			}
			return false, nil
		}
	}

	if _, ok := node.(AlphaNum); ok {
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]$`, current); matched {
			p.context.skipChars = false
			p.advance()
			return true, nil
		} else {
			if p.context.skipChars {
				p.advance()
				result, _ := p.CheckParseTree(node)
				res := false || result
				if res {
					p.context.skipChars = false
				}
				return res, nil
			}
			return false, nil
		}
	}

	if _, ok := node.(Digit); ok {
		if matched, _ := regexp.MatchString(`^[0-9]$`, current); matched {
			p.context.skipChars = false
			p.advance()
			return true, nil
		} else {
			if p.context.skipChars {
				p.advance()
				result, _ := p.CheckParseTree(node)
				res := false || result
				if res {
					p.context.skipChars = false
				}
				return res, nil
			}
			return false, nil
		}
	}

	return false, nil
}

func (p *Parser) checkRightSubtree(rightNode RegexpNode) bool {
	// defer func() { p.context.skipChars = true }()
	result, _ := p.CheckParseTree(rightNode)
	return result
}

func getExpFromTree(node RegexpNode) []string {
	var res []string
	var inner func(RegexpNode)

	inner = func(node RegexpNode) {
		if c, ok := node.(Literal); ok {
			// println("value", c.value)
			// println(string(c.value))
			res = append(res, string(c.value))
		}
		if c, ok := node.(Alternate); ok {
			inner(c.leftNode)
			inner(c.rightNode)
		}
	}
	inner(node)
	return res
}

// Call this when it's needed to parse AST once the 1st AST node == char of a checkstring at some index i
// [^], [] are not included as they do not need to be checked this way
func getFirstLeftChild(node RegexpNode) (string, error) {
	if c, ok := node.(Concat); ok {
		return getFirstLeftChild(c.leftNode)
	}

	if _, ok := node.(Repeat); ok {
		return "Repeat()", nil
	}

	if _, ok := node.(Optional); ok {
		return "Optional()", nil
	}

	if _, ok := node.(Alternate); ok {
		return "Alternate()", nil
	}

	if c, ok := node.(Wildcard); ok {
		return c.get(), nil
	}

	if c, ok := node.(Literal); ok {
		return c.get(), nil
	}

	if _, ok := node.(Digit); ok {
		return `^[0-9]$`, nil
	}

	if _, ok := node.(AlphaNum); ok {
		return `^[0-9a-zA-Z]$`, nil
	}

	if c, ok := node.(AnchorStart); ok {
		return c.get(), nil
	}

	// NOTE: AnchorEnd not needed as it will never be a left child it will always be a right child

	return "", fmt.Errorf("Unsupported pattern:", node.get())

}
