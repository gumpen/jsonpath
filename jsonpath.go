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

		// normal property
		if !strings.Contains(s, "[") {
			p.tokens = append(p.tokens, Token{Type: TokenTypeKey, Value: s})
			continue
		}

		// contains bracket

		start := strings.Index(s, "[")
		end := strings.LastIndex(s, "]")

		if end < 0 || start >= end {
			return fmt.Errorf("invalid query")
		}
		p.tokens = append(p.tokens, Token{Type: TokenTypeKey, Value: s[0:start]})

		exp := s[start+1 : end]

		// union
		if strings.Contains(exp, ",") {
			p.tokens = append(p.tokens, Token{Type: TokenTypeUnion, Value: exp})
			continue
		}

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
		case TokenTypeUnion:
			unionResults := make([]interface{}, 0)
			nums := strings.Split(t.Value, ",")
			for _, n := range nums {
				r, err := p.findValue(n, result)
				if err != nil {
					return nil, err
				}
				unionResults = append(unionResults, r)
			}

			result = unionResults
		case TokenTypeIndex:
			r, err := p.findValue(t.Value, result)
			if err != nil {
				return nil, err
			}
			result = r
		case TokenTypeKey:
			if resultElem, ok := result.([]interface{}); ok {
				res := make([]interface{}, 0)
				for _, e := range resultElem {
					r, err := p.findValue(t.Value, e)
					if err != nil {
						return nil, err
					}

					res = append(res, r)
				}

				result = res
				continue
			}

			r, err := p.findValue(t.Value, result)
			if err != nil {
				return nil, err
			}
			result = r
		default:
			return nil, fmt.Errorf("invalid token")
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
