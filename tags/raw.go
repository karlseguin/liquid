package tags

import (
	"github.com/karlseguin/liquid/core"
)

var theOneRaw = new(Raw)

func RawFactory(data []byte) (core.Token, error) {
	return theOneRaw, nil
}

type Raw struct {
}

func (r *Raw) Name() string {
	return "raw"
}

func (r *Raw) IsEnd() bool {
	return false
}

func (r *Raw) Render(data map[string]interface{}) []byte {
	return []byte{}
}

func (r *Raw) AddToken(token core.Token) {
}

func (r *Raw) AddTag(tag core.Tag) (bool, bool) {
	if tag.Name() == "endraw" {
		return true, false
	}
	return false, true
}
