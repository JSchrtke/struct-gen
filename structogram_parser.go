package main

import (
	"errors"
	"fmt"
	"strings"
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
}

// TODO find a better name for this
type ParsedObject struct {
	name string
}

func parseTokens(tokens []Token) (ParsedObject, error) {
	var parsed ParsedObject
	var err error
	if tokens[0].tokenType == "name" {
		// The next token should be a openParentheses
		if tokens[2].tokenType == "name" {
			return parsed, errors.New(
				fmt.Sprintf(
					"%d:%d, names can not be nested",
					tokens[2].line,
					tokens[2].column,
				),
			)
		}
		if tokens[2].tokenType != "string" {
			return parsed, errors.New(
				fmt.Sprintf(
					"%d:%d, missing name",
					tokens[1].line,
					tokens[1].column,
				),
			)
		}
		parsed.name = tokens[2].value
	} else {
		return parsed, errors.New(
			fmt.Sprintf(
				"%d:%d, structogram has to start with a name",
				tokens[0].line,
				tokens[0].column,
			),
		)
	}
	return parsed, err
}

func parseToken(s, tokenStart, tokenEnd string) (content, remaining string) {
	tokenStartIndex := strings.Index(s, tokenStart)
	tokenEndIndex := strings.Index(s[tokenStartIndex:], tokenEnd) + tokenStartIndex
	content = s[tokenStartIndex+len(tokenStart) : tokenEndIndex]
	remaining = s[tokenEndIndex:]
	return
}

func parseStructogram(structogram string) (*Structogram, error) {
	if len(structogram) == 0 {
		return nil, errors.New("Parsing error, structogram string is empty!")
	}

	nameToken := "name("
	if strings.Index(structogram, nameToken) != 0 {
		return nil, errors.New("Structogram must have a name!")
	}

	parsed := Structogram{}

	var remaining string
	parsed.name, remaining = parseToken(structogram, nameToken, ")")

	if len(parsed.name) == 0 {
		return nil, errors.New("Structograms can not have empty names!")
	}
	if strings.Contains(parsed.name, nameToken) {
		return nil, errors.New("Structogram names can not be nested!")
	}

	instructionToken := "instruction("
	var instruction string
	for strings.Contains(remaining, instructionToken) {
		instruction, remaining = parseToken(remaining, instructionToken, ")")
		if len(instruction) == 0 {
			return nil, errors.New("Instructions can not be empty!")
		}
		if strings.Contains(instruction, instructionToken) {
			return nil, errors.New("Instructions can not be nested!")
		}
		parsed.instructions = append(
			parsed.instructions,
			instruction,
		)
	}

	return &parsed, nil
}
