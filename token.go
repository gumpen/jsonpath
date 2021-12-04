package jsonpath

type TokenType int

const (
	TokenTypeRoot TokenType = iota
	TokenTypeIndex
	TokenTypeKey
	TokenTypeUnion
)

type Token struct {
	Type  TokenType
	Value string
}
