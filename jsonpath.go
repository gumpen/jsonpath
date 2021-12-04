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
	result []interface{}
}

func NewPath(q string) *Path {
	return &Path{
		query:  q,
		result: make([]interface{}, 0),
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

		// slice
		if strings.Contains(exp, ":") {
			p.tokens = append(p.tokens, Token{Type: TokenTypeSlice, Value: exp})
		}

		p.tokens = append(p.tokens, Token{Type: TokenTypeIndex, Value: s[start+1 : end]})

	}

	// p.tokens = splitted

	return nil
}

func (p *Path) Execute(data interface{}) (interface{}, error) {
	output, err := p.find(p.tokens, data)
	if err != nil {
		return "", err
	}

	result, err := p.format(output)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *Path) format(data interface{}) (interface{}, error) {
	err := p.makeResult(data)
	if err != nil {
		return nil, err
	}

	if len(p.result) < 1 {
		return nil, nil
	} else if len(p.result) == 1 {
		return p.result[0], nil
	}

	return p.result, nil
}

func (p *Path) makeResult(data interface{}) error {
	if a, ok := data.([]interface{}); ok {
		for _, b := range a {
			err := p.makeResult(b)
			if err != nil {
				return err
			}
		}
	} else {
		p.result = append(p.result, data)
	}

	return nil
}

func (p *Path) find(tokens []Token, data interface{}) (interface{}, error) {
	var result = data
	for i, t := range tokens {
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

			if len(tokens) <= i+1 {
				result = unionResults
				continue
			}

			findResults := make([]interface{}, 0)
			for _, ur := range unionResults {
				findRes, err := p.find(tokens[i+1:], ur)
				if err != nil {
					return nil, err
				}

				findResults = append(findResults, findRes)
			}

			return findResults, nil
		case TokenTypeSlice:

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

func makeIndexSlice(start, end string) ([]string, error) {
	res := []string{}
	si, err := strconv.Atoi(start)
	if err != nil {
		return res, err
	}
	ei, err := strconv.Atoi(end)
	if err != nil {
		return res, err
	}

	for i := si; i < ei; i++ {
		res = append(res, strconv.Itoa(i))
	}

	return res, nil
}
