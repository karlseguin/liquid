// liquid template parser
package liquid

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
	"github.com/karlseguin/liquid/core"
	"io/ioutil"
)

type TokenExtractor func(data []byte) (core.Token, error)

// A compiled liquid template
type Template struct {
	Tokens []core.Token
}

func (t *Template) AddToken(token core.Token) {
	t.Tokens = append(t.Tokens, token)
}

func (t *Template) AddTag(tag core.Tag) (bool, bool) {
	t.AddToken(tag)
	return false, false
}

func (t *Template) Name() string {
	return "root"
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

func (t *Template) Render(data map[string]interface{}) []byte {
	buffer := new(bytes.Buffer)
	for _, token := range t.Tokens {
		buffer.Write(token.Render(data))
	}
	return buffer.Bytes()
}

func buildTemplate(data []byte) (*Template, error) {
	template := new(Template)
	if err := extractTokens(data, template); err != nil {
		return nil, err
	}
	return template, nil
}

func extractTokens(data []byte, container core.Container) error {
	start := 0
	isLiteral := false
	var extractor TokenExtractor
	stack := make([]core.Container, 0, 0)

	for i, l := 0, len(data)-1; i < l; i++ {
		current := data[i]
		if current == '{' {
			next := data[i+1]
			if next == '{' || next == '%' {
				if extractor != nil {
					if isLiteral == false {
						return errors.New(fmt.Sprintf("Invalid escape sequence %q", data[start:i]))
					}
					if err := doExtraction(extractor, data[start:i], container); err != nil {
						return err
					}
				}
				start = i
				isLiteral = false
				if next == '{' {
					extractor = outputExtractor
				} else {
					extractor = tagExtractor
				}
			}
		} else if current == '}' && i > 0 {
			prev := data[i-1]
			if prev == '}' {
				if err := doExtraction(extractor, data[start:i], container); err != nil {
					return err
				}
				extractor = nil
			} else if prev == '%' {
				token, err := extractor(data[start:i])
				if err != nil {
					return err
				}
				extractor = nil
				if token != nil {
					tag := token.(core.Tag)
					name := tag.Name()
					if name == "raw" {
						length, literal := extractRaw(data[i+1:])
						if literal == nil {
							return errors.New("unclosed {% raw %} tag")
						}
						container.AddToken(literal)
						i+= length
					} else if closed, related := container.AddTag(tag); closed {
						stack, container = popContainerStack(stack)
					} else if tag.IsEnd() {
						return errors.New(fmt.Sprintf("Was not expecting a the closing tag %q", name))
					} else if related == false {
						stack = append(stack, container)
						container = tag
					}
				}
			}
		} else if extractor == nil {
			extractor = literalExtractor
			isLiteral = true
			start = i
		}
	}
	if extractor != nil {
		if err := doExtraction(extractor, data[start:len(data)], container); err != nil {
			return err
		}
	}
	return nil
}

func extractRaw(data []byte) (int, core.Token) {
	l := len(data) - 2
	for i := 0; i < l; i++ {
		if data[i] == '{' && data[i+1] == '%' {
			start := core.SkipSpaces(data[i+2:])
			if start == -1 {
				return 0, nil
			}
			start += i + 2
			for j := start ; j < l; j++ {
				if data[j] == ' ' || data[j] == '%' {
					if strings.ToLower(string(data[start:j])) == "endraw" {
						token, _ := literalExtractor(data[:i])
						for ; j < l; j++ {
							if data[j] == '}' { break }
						}
						return j+1, token
					}
				}
			}
		}
	}
	return 0, nil
}

func doExtraction(extractor TokenExtractor, data []byte, container core.Container) error {
	token, err := extractor(data)
	if err != nil {
		return err
	}
	container.AddToken(token)
	return nil
}

func popContainerStack(stack []core.Container) ([]core.Container, core.Container) {
	l := len(stack) - 1
	return stack[0:l], stack[l]
}
