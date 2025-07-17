package myast

type ParserContext struct {
	isAlternate bool
	skipChars   bool
}

func NewParserContext() *ParserContext {
	return &ParserContext{false, true}
}
