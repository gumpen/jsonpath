package jsonpath

type TokenType string

const (
	TokenTypeRoot  TokenType = "root"
	TokenTypeIndex TokenType = "index"
	TokenTypeKey   TokenType = "key"
)

type Token struct {
	Type  TokenType
	Value string
}
