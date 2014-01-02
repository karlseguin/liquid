package liquid

import (
	"errors"
	"strings"
	"fmt"
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/tags"
)

// A tag factory creates a tag based on the raw data
type TagFactory func(data, all []byte) (core.Token, error)

var TagEnds = map[string]core.Tag {
	"endcomment": tags.EndComment,
	"endraw": tags.EndRaw,
}

var Tags = map[string]TagFactory{
	"comment": tags.CommentFactory,
	"raw": tags.RawFactory,
	"assign": tags.AssignFactory,
}

func tagExtractor(all []byte) (core.Token, error) {
	//strip out leading and trailing {% %}
	data := all[2 : len(all)-2]
	start := 0
	if start = core.SkipSpaces(data); start == -1 {
		return nil, nil
	}
	l := len(data)
	i := start
	for ;i < l; i++ {
		if data[i] == ' ' || data[i] == '%' { break }
	}
	tagName := strings.ToLower(string(data[start:i]))
	if end, exists := TagEnds[tagName]; exists {
		return end, nil
	}

	if factory, exists := Tags[tagName]; exists {
		if i < l {
			data = data[i+1:]
		}
		return factory(data, all)
	}
	return nil, errors.New(fmt.Sprintf("unknown tag %q at %q", tagName, all))
}
