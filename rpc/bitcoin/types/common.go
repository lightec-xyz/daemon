package types

type Param map[string]interface{}

func (p *Param) Add(key string, value interface{}) {
	(*p)[key] = value
}
