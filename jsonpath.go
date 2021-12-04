package jsonpath

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const separator = "."

type Path struct {
	query  string
	tokens []Token
}

func NewPath(q string) *Path {
	return &Path{
		query: q,
	}
}

func (p *Path) Parse() error {
	err := p.tokenize()
	if err != nil {
		return err
	}

	return nil
}

func (p *Path) tokenize() error {
	splitted := strings.Split(p.query, separator)

	if len(splitted) > 0 && splitted[0] != "$" {
		return fmt.Errorf("invalid query")
	}

	p.tokens = append(p.tokens, Token{Type: TokenTypeRoot, Value: splitted[0]})

	for _, s := range splitted[1:] {
		if !strings.Contains(s, "[") {
			p.tokens = append(p.tokens, Token{Type: TokenTypeKey, Value: s})
			continue
		}

		start := strings.Index(s, "[")
		end := strings.LastIndex(s, "]")

		if end < 0 || start >= end {
			return fmt.Errorf("invalid query")
		}

		p.tokens = append(p.tokens, Token{Type: TokenTypeKey, Value: s[0:start]})
		p.tokens = append(p.tokens, Token{Type: TokenTypeIndex, Value: s[start+1 : end]})

	}

	// p.tokens = splitted

	return nil
}

func (p *Path) Execute(data interface{}) (interface{}, error) {
	output, err := p.find(data)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (p *Path) find(data interface{}) (interface{}, error) {
	var result = data
	for _, t := range p.tokens {
		switch t.Type {
		case TokenTypeRoot:
			result = data
		default:
			r, err := p.findValue(t.Value, result)
			if err != nil {
				return nil, err
			}
			result = r
		}

		if result == nil {
			return nil, nil
		}
	}

	return result, nil
}

func (p *Path) findValue(q string, data interface{}) (interface{}, error) {
	t := reflect.TypeOf(data)

	switch t.Kind() {
	case reflect.Slice:
		v := reflect.ValueOf(data)
		qn, err := strconv.Atoi(q)
		if err != nil {
			return nil, err
		}

		return v.Index(qn).Interface(), nil
	case reflect.Map:
		v := reflect.ValueOf(data)
		for _, k := range v.MapKeys() {
			if q == k.String() {
				return v.MapIndex(k).Interface(), nil
			}
		}
	}

	return nil, nil
}
