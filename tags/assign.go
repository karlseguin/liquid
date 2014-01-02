package tags

import (
	// "errors"
	// "fmt"
	"github.com/karlseguin/liquid/core"
)

func AssignFactory(data, all []byte) (core.Token, error) {
	return nil, nil
	// var (
	// 	variable string
	// 	ok bool
	// 	values []string
	// )

	// if data, variable, ok = core.ExtractVariable(data); ok == false {
	// 	return nil, errors.New(fmt.Sprintf("invalid or missing variable in assign %q", all))
	// }
	// if data, ok = core.ExtractByte(data, '='); ok == false {
	// 	return nil, errors.New(fmt.Sprintf("missing = in assign %q", all))
	// }
	// if data, values, ok = core.ExtractValue(data); ok == false {
	// 	return nil, errors.New(fmt.Sprintf("invalid or missing value in assign %q", all))
	// }
	// return &Assign{variable, values}, nil
}

type Assign struct {
	variable string
	values []string
}

func (a *Assign) Name() string {
	return "assign"
}

func (a *Assign) IsEnd() bool {
	return false
}

func (a *Assign) Render(data map[string]interface{}) []byte {
	return []byte{}
}

func (a *Assign) AddToken(token core.Token) {
	panic("AddToken should not be called on an assign tag")
}

func (a *Assign) AddTag(tag core.Tag) (bool, bool) {
	panic("AddTag should not be called on an assign tag")
}
