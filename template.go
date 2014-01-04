// liquid template parser
package liquid

import (
	"bytes"
	"crypto/sha1"
	// "errors"
	"fmt"
	// "strings"
	"github.com/karlseguin/liquid/core"
	"io/ioutil"
)

// A compiled liquid template
type Template struct {
	Code []core.Code
}

func (t *Template) AddCode(code core.Code) {
	t.Code = append(t.Code, code)
}

func (t *Template) AddSibling(tag core.Tag) error {
	return nil
}

func (t *Template) Type() core.TagType {
	return core.ContainerTag
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
	if data == nil {
		data = make(map[string]interface{})
	}
	buffer := new(bytes.Buffer)
	for _, code := range t.Code {
		buffer.Write(code.Render(data))
	}
	return buffer.Bytes()
}

func buildTemplate(data []byte) (*Template, error) {
	parser := core.NewParser(data)
	template := new(Template)
	if err := extractTokens(parser, template); err != nil {
		return nil, err
	}
	return template, nil
}

func extractTokens(parser *core.Parser, container core.Tag) error {
	stack := make([]core.Tag, 0, 0)
	for parser.HasMore() {
		pre, markupType := parser.ToMarkup()
		if len(pre) > 0 {
			container.AddCode(newLiteral(pre))
		}
		if markupType == core.OutputMarkup {
			code, err := newOutput(parser)
			if err != nil {
				return err
			}
			if code != nil {
				container.AddCode(code)
			}
		} else if markupType == core.TagMarkup {
			start := parser.Position
			tag, err := newTag(parser)
			if err != nil {
				return err
			}
			switch tag.Type() {
			case core.ContainerTag:
				stack = append(stack, container)
				container = tag
			case core.EndTag:
				if tag.Name() != container.Name() {
					return parser.Error("unexpected end tag", start)
				}
				l := len(stack) - 1
				container = stack[l]
				stack = stack[0:l]
			case core.SiblingTag:
				if err := container.AddSibling(tag); err != nil {
					return err
				}
			case core.StandaloneTag:
				container.AddCode(tag)
			}
		} else {
			break
		}
	}
	return nil
}
