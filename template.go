// liquid template parser
package liquid

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
)

type TokenExtractor func(data []byte) (Token, error)

type Token interface {
	Render(data interface{}) []byte
}

// A compiled liquid template
type Template struct {
	Tokens []Token
}

// Parse the bytes into a Liquid template
func Parse(data []byte, config *Configuration) (*Template, error) {
	if config == nil {
		config = defaultConfig
	}
	if config.cache == nil {
		return buildTemplate(data)
	}
	hasher := sha1.New()
	hasher.Write(data)
	key := fmt.Sprintf("%x", hasher.Sum(nil))

	template := config.cache.Get(key)
	if template == nil {
		var err error
		template, err = buildTemplate(data)
		if err != nil {
			return nil, err
		}
		config.cache.Set(key, template)
	}
	return template, nil
}

// Parse the string into a liquid template
func ParseString(data string, config *Configuration) (*Template, error) {
	return Parse([]byte(data), config)
}

// Turn the contents of the specified file into a liquid template
func ParseFile(path string, config *Configuration) (*Template, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data, config)
}

func (t *Template) Render(data interface{}) []byte {
	buffer := new(bytes.Buffer)
	for _, token := range t.Tokens {
		buffer.Write(token.Render(data))
	}
	return buffer.Bytes()
}

func buildTemplate(data []byte) (*Template, error) {
	tokens, err := extractTokens(data)
	if err != nil {
		return nil, err
	}
	return &Template{Tokens: tokens}, nil
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
