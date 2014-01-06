package tags

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/core"
)

var (
	endUnless = &End{"unless"}
)

func UnlessFactory(p *core.Parser) (core.Tag, error) {
	condition, err := p.ReadConditionGroup()
	if err != nil {
		return nil, err
	}
	p.SkipPastTag()
	condition.Inverse()
	return &Unless{NewCommon(), condition, nil}, nil
}

func EndUnlessFactory(p *core.Parser) (core.Tag, error) {
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

func (u *Unless) Render(data map[string]interface{}) []byte {
	if u.condition.IsTrue(data) {
		return u.Common.Render(data)
	}
	if u.elseCondition != nil {
		return u.elseCondition.Render(data)
	}
	return nil
}

func (u *Unless) Name() string {
	return "unless"
}

func (u *Unless) Type() core.TagType {
	return core.ContainerTag
}
