package myast

type ParserContext struct {
	isAlternate bool
	skipChars   bool
	stackDepth  int
	anchorStart bool
}

func NewParserContext() *ParserContext {
	return &ParserContext{false, true, 0, false}
}
