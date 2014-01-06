// liquid template parser
package liquid

import (
	"crypto/sha1"
	"fmt"
	"github.com/karlseguin/liquid/core"
	"io"
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
func Parse(data []byte, config *core.Configuration) (*Template, error) {
	if config == nil {
		config = defaultConfig
	}
	cache := config.GetCache()
	if cache == nil {
		return buildTemplate(data, config)
	}
	hasher := sha1.New()
	hasher.Write(data)
	key := fmt.Sprintf("%x", hasher.Sum(nil))

	template := cache.Get(key)
	if template == nil {
		var err error
		template, err = buildTemplate(data, config)
		if err != nil {
			return nil, err
		}
		cache.Set(key, template)
	}
	return template.(*Template), nil
}

// Parse the string into a liquid template
func ParseString(data string, config *core.Configuration) (*Template, error) {
	return Parse([]byte(data), config)
}

// Turn the contents of the specified file into a liquid template
func ParseFile(path string, config *core.Configuration) (*Template, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data, config)
}

func (t *Template) Render(writer io.Writer, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	for _, code := range t.Code {
		code.Render(writer, data)
	}
}

func buildTemplate(data []byte, config *core.Configuration) (*Template, error) {
	parser := core.NewParser(data)
	template := new(Template)
	if err := extractTokens(parser, template, config); err != nil {
		return nil, err
	}
	return template, nil
}

func extractTokens(parser *core.Parser, container core.Tag, config *core.Configuration) error {
	stack := make([]core.Tag, 0, 0)
	parentContainer := container
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
			tag, err := newTag(parser, config)
			if err != nil {
				return err
			}
			switch tag.Type() {
			case core.ContainerTag:
				container.AddCode(tag)
				stack = append(stack, container)
				container = tag
				parentContainer = tag
			case core.EndTag:
				if tag.Name() != parentContainer.Name() {
					return parser.Error("unexpected end tag", start)
				}
				l := len(stack) - 1
				container = stack[l]
				parentContainer = nil
				stack = stack[0:l]
				parser.SkipPastTag()
			case core.SiblingTag:
				if err := parentContainer.AddSibling(tag); err != nil {
					return err
				}
				container = tag
			case core.StandaloneTag:
				container.AddCode(tag)
			}
		} else {
			break
		}
	}
	return nil
}
