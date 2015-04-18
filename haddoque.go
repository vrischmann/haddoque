package haddoque

type Parser interface {
	Parse(query string) (Statement, error)
}

type parser struct {
}

func NewQueryParser() Parser {
	return &parser{}
}

func (p *parser) Parse(query string) (Statement, error) {
	return nil, nil
}

type Statement interface {
	Exec(object interface{}) (interface{}, error)
}
