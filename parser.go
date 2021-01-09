package main

import (
	"errors"
	"fmt"
)

type Token struct {
	tokenType string
	value     string
	line      int
	column    int
}

type Structogram struct {
	name         string
	instructions []string
	nodes        []Node
}

type Node struct {
	nodeType string
	value    string
	nodes    []Node
}

type Parser struct {
	tokenIndex int
	tokens     []Token
}

func (p *Parser) next() Token {
	return p.tokens[p.tokenIndex]
}

func (p *Parser) readNext() Token {
	t := p.next()
	p.tokenIndex++
	return t
}

func newTokenValueError(expected string, actual Token) error {
	return errors.New(
		fmt.Sprintf(
			"%d:%d, expected '%s', but got '%s'",
			actual.line,
			actual.column,
			expected,
			actual.tokenType,
		),
	)
}

func isKeyword(s string) bool {
	return s == "instruction" || s == "if"
}

func (p *Parser) parseIf(parsed Structogram) ([]Node, error) {
	var nodes []Node
	var ifNode Node

	ifNode.nodeType = p.readNext().value
	if p.next().tokenType != "openParentheses" {
		return nil, newTokenValueError("openParentheses", p.next())
	}
	// Discard the openParentheses
	p.readNext()

	if p.next().tokenType != "string" {
		return nil, newTokenValueError("string", p.next())
	}
	ifNode.value = p.readNext().value

	if p.next().tokenType != "closeParentheses" {
		return nil, newTokenValueError("closeParentheses", p.next())
	}
	p.readNext()
	if p.next().tokenType == "whitespace" {
		p.readNext()
	}
	if p.next().tokenType != "openBrace" {
		return nil, newTokenValueError("openBrace", p.next())
	}
	p.readNext()

	if !isKeyword(p.next().tokenType) {
		return nil, newTokenValueError("keyword", p.next())
	}

	// Parsing of the if body
	for p.next().tokenType != "closeBrace" {
		tok := p.readNext()
		switch tok.tokenType {
		case "instruction":
			var instructionNode Node
			instructionNode.nodeType = tok.tokenType

			if p.next().tokenType != "openParentheses" {
				return nil, newTokenValueError(
					"openParentheses", p.next(),
				)
			}
			p.readNext()

			if p.next().tokenType != "string" {
				return nil, newTokenValueError("string", p.next())
			}
			instructionNode.value = p.readNext().value

			if p.next().tokenType != "closeParentheses" {
				return nil, newTokenValueError(
					"closeParentheses", p.next(),
				)
			}
			p.readNext()
			ifNode.nodes = append(ifNode.nodes, instructionNode)
		}
	}

	if p.next().tokenType != "closeBrace" {
		return nil, newTokenValueError("closeBrace", p.next())
	}
	p.readNext()
	var err error
	nodes = append(nodes, ifNode)
	return nodes, err
}

func parseTokens(tokens []Token) (Structogram, error) {
	p := Parser{
		tokenIndex: 0,
		tokens:     tokens,
	}
	var parsed Structogram
	var err error
	if p.next().tokenType != "name" {
		return parsed, newTokenValueError("name", p.next())
	}
	for {
		switch p.next().tokenType {
		case "name":
			_ = p.readNext()
			if p.next().tokenType != "openParentheses" {
				return parsed, newTokenValueError("openParentheses", p.next())
			}
			tok := p.readNext()
			if p.next().tokenType != "string" {
				return parsed, newTokenValueError("string", p.next())
			}
			tok = p.readNext()
			parsed.name = tok.value
			if p.next().tokenType != "closeParentheses" {
				return parsed, newTokenValueError("closeParentheses", p.next())
			}
			_ = p.readNext()
		case "instruction":
			_ = p.readNext()
			if p.next().tokenType != "openParentheses" {
				return parsed, newTokenValueError("openParentheses", p.next())
			}
			_ = p.readNext()
			if p.next().tokenType != "string" {
				return parsed, newTokenValueError("string", p.next())
			}
			parsed.instructions = append(parsed.instructions, p.readNext().value)
			if p.next().tokenType != "closeParentheses" {
				return parsed, newTokenValueError("closeParentheses", p.next())
			}
			_ = p.readNext()
		case "whitespace":
			// TODO Maybe have a function that runs once that strips all the
			// whitespace out of the tokens?
			// Whitespace should be completely ignored
			_ = p.readNext()
		case "if":
			parsed.nodes, err = p.parseIf(parsed)
			if err != nil {
				return parsed, err
			}
		case "invalid":
			return parsed, newTokenValueError("identifier", p.next())
		case "EOF":
			return parsed, err
		}
	}
}
