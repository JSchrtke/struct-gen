package main

import "errors"

type parser struct{}

func createParser() (*parser, error) {
	p := parser{}
	return &p, nil
}

type parsedObject struct{}

func (p *parser) parseStructogram(s string) (*parsedObject, error) {
	err := errors.New("Parsing error, structogram string is empty!")
	return &parsedObject{}, err
}
