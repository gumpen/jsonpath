package jsonpath

type TokenType int

const (
	TokenTypeRoot TokenType = iota
	TokenTypeIndex
	TokenTypeKey
	TokenTypeUnion
	TokenTypeSlice
)

type Token struct {
	Type  TokenType
	Value string
}
