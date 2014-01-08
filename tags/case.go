package tags

import (
	"errors"
	"fmt"
	"github.com/karlseguin/liquid/core"
	"io"
)

var (
	endCase = &End{"case"}
)

type CaseSibling interface {
	Condition() core.Verifiable
	AddLeftValue(value core.Value)
	core.Code
}

func CaseFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	value, err := p.ReadValue()
	if err != nil {
		return nil, err
	}
	p.SkipPastTag()
	return &Case{value, make([]CaseSibling, 0, 5)}, nil
}

func WhenFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	condition, err := p.ReadPartialCondition()
	if err != nil {
		return nil, err
	}
	p.SkipPastTag()
	return &When{NewCommon(), condition}, nil
}

func EndCaseFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	return endCase, nil
}

type Case struct {
	value      core.Value
	conditions []CaseSibling
}

func (c *Case) AddCode(code core.Code) {}

func (c *Case) AddSibling(tag core.Tag) error {
	cs, ok := tag.(CaseSibling)
	if ok == false {
		return errors.New(fmt.Sprintf("%q does not belong inside of a case"))
	}
	c.conditions = append(c.conditions, cs)
	cs.AddLeftValue(c.value)
	return nil
}

func (c *Case) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	for _, condition := range c.conditions {
		if condition.Condition().IsTrue(data) {
			return condition.Execute(writer, data)
		}
	}
	return core.Normal
}

func (c *Case) Name() string {
	return "case"
}

func (c *Case) Type() core.TagType {
	return core.ContainerTag
}

type When struct {
	*Common
	condition core.Completable
}

func (w *When) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a when")
}

func (w *When) Name() string {
	return "when"
}

func (w *When) Type() core.TagType {
	return core.SiblingTag
}

func (w *When) Condition() core.Verifiable {
	return w.condition
}

func (w *When) AddLeftValue(left core.Value) {
	w.condition.Complete(left, core.Equals)
}
