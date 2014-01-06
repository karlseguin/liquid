package tags

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/core"
	"io"
)

var (
	endUnless = &End{"unless"}
)

func UnlessFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	condition, err := p.ReadConditionGroup()
	if err != nil {
		return nil, err
	}
	p.SkipPastTag()
	condition.Inverse()
	return &Unless{NewCommon(), condition, nil}, nil
}

func EndUnlessFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endUnless, nil
}

type Unless struct {
	*Common
	condition     core.Verifiable
	elseCondition *Else
}

func (u *Unless) AddSibling(tag core.Tag) error {
	e, ok := tag.(*Else)
	if ok == false {
		return errors.New(fmt.Sprintf("%q does not belong as a sibling of an unless"))
	}
	u.elseCondition = e
	return nil
}

func (u *Unless) Render(writer io.Writer, data map[string]interface{}) {
	if u.condition.IsTrue(data) {
		u.Common.Render(writer, data)
	} else if u.elseCondition != nil {
		u.elseCondition.Render(writer, data)
	}
}

func (u *Unless) Name() string {
	return "unless"
}

func (u *Unless) Type() core.TagType {
	return core.ContainerTag
}
