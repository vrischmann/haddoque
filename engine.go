package haddoque

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(query string, obj interface{}) (interface{}, error) {
	return nil, nil
}
