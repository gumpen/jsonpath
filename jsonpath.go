package jsonpath

import (
	"reflect"
	"strings"
)

const separator = "."

type Path struct {
	query  string
	tokens []string
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

	p.tokens = splitted

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
		switch t {
		case "$":
			result = data
		default:
			r, err := p.findValue(t, result)
			if err != nil {
				return nil, err
			}
			result = r
		}
	}

	return result, nil
}

func (p *Path) findValue(q string, data interface{}) (interface{}, error) {
	t := reflect.TypeOf(data)

	switch t.Kind() {
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
