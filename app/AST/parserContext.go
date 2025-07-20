package myast

type ParserContext struct {
	isAlternate bool
	skipChars   bool
	stackDepth  int
}

func NewParserContext() *ParserContext {
	return &ParserContext{false, true, 0}
}
