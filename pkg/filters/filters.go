package filters

import (
	"errors"
)

type Filter interface {
	Judge(data []byte) (exists bool, err error)
	Reset(params *FilterParams) (err error)
	Euqal(params *FilterParams) bool
	Info() map[string]interface{}
	Delete()
}

type DefaultFilter struct {
}

func NewDefaultFilter() (Filter, error) {
	return &DefaultFilter{}, nil
}

func (filter *DefaultFilter) Judge(data []byte) (exists bool, err error) {
	return false, nil
}

func (filter *DefaultFilter) Reset(params *FilterParams) (err error) {
	if params.Kind != None {
		err = errors.New("kind not same")
		return
	}
	return nil
}

func (filter *DefaultFilter) Euqal(params *FilterParams) bool {
	return false
}

func (filter *DefaultFilter) Info() map[string]interface{} {
	result := make(map[string]interface{})
	result["info"] = &FilterParams{
		None, 0, 0, 0,
	}
	result["rate"] = 0
	return result
}

func (filter *DefaultFilter) Delete() {

}
