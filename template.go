package liquid

import (
	"errors"
	"fmt"
	"io/ioutil"
)

type TokenExtractor func(data []byte) (Token, error)

type Token interface {
}

type Template struct {
	tokens []Token
}

func Parse(data []byte) (*Template, error) {
	tokens, err := extractTokens(data)
	if err != nil {
		return nil, err
	}
	return &Template{
		tokens: tokens,
	}, nil
}

func ParseString(data string) (*Template, error) {
	return Parse([]byte(data))
}

func ParseFile(path string) (*Template, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func extractTokens(data []byte) ([]Token, error) {
	start := 0
	isLiteral := false
	var err error
	var extractor TokenExtractor
	tokens := make([]Token, 0, 5)

	for i, l := 0, len(data)-1; i < l; i++ {
		current := data[i]
		if current == '{' {
			next := data[i+1]
			if next == '{' || next == '%' {
				if extractor != nil {
					if isLiteral == false {
						return nil, errors.New(fmt.Sprintf("Invalid escape sequence %q", data[start:i]))
					}
					if tokens, err = doExtraction(extractor, data[start:i], tokens); err != nil {
						return nil, err
					}
				}
				start = i
				isLiteral = false
				if next == '{' {
					extractor = outputExtractor
				} else {

				}
			}
		} else if current == '}' && i > 0 {
			prev := data[i-1]
			if prev == '}' || prev == '%' {
				if tokens, err = doExtraction(extractor, data[start:i], tokens); err != nil {
					return nil, err
				}
				extractor = nil
			}
		} else if extractor == nil {
			extractor = literalExtractor
			isLiteral = true
			start = i
		}
	}
	if tokens, err = doExtraction(extractor, data[start:len(data)], tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

func doExtraction(extractor TokenExtractor, data []byte, tokens []Token) ([]Token, error) {
	token, err := extractor(data)
	if err != nil {
		return nil, err
	}
	return append(tokens, token), nil
}
