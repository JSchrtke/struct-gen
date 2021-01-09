package main

import "testing"

func checkOk(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Did not expect any errors, but got %s", err.Error())
	}
}

func checkErrorMsg(t *testing.T, err error, expectedMsg string) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error but was nil")
	}
	if err.Error() != expectedMsg {
		t.Errorf(
			"Expected error with msg %s, but got %s",
			expectedMsg,
			err.Error(),
		)
	}
}

func checkNodeCount(t *testing.T, n []Node, count int) {
	t.Helper()
	if len(n) != count {
		t.Errorf("Wrong node count, expected %d, but got %d", count, len(n))
	}
}

func TestEmptyStructogramNameCausesError(t *testing.T) {
	tokens := makeTokens("name()")
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:6, expected 'string', but got 'closeParentheses'")
}

func TestStructogramsHaveNames(t *testing.T) {
	tokens := makeTokens(`name("test name")`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)

	if structogram.name != "test name" {
		t.Errorf(
			"Diagram has wrong name, expected: %s, but was: %s",
			"test name", structogram.name,
		)
	}
}

func TestNamesCanNotBeNested(t *testing.T) {
	tokens := makeTokens("name(name())")
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:6, expected 'string', but got 'name'")
}

func TestNameHasToBeFirstToken(t *testing.T) {
	tokens := makeTokens(`instruction("something")name("a name")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:1, expected 'name', but got 'instruction'")
}

func TestNameValueHasToBeEnclosedByParentheses(t *testing.T) {

	tokens := makeTokens(`name"a name"`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:5, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a"(`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:9, expected 'closeParentheses', but got 'openParentheses'")
}

func TestInstructionValueHasToBeEnclosedByParentheses(t *testing.T) {
	tokens := makeTokens(`name("some name")instruction"something")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:29, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a")instruction("b"(`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:25, expected 'closeParentheses', but got 'openParentheses'")
}

func TestInstructionsCanNotBeEmpty(t *testing.T) {
	tokens := makeTokens(`name("test structogram")instruction()`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:37, expected 'string', but got 'closeParentheses'")
}

func TestInstuctionsCanNotBeNested(t *testing.T) {
	tokens := makeTokens(`name("a")instruction(instruction())`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:22, expected 'string', but got 'instruction'")
}

func TestStructogramCanHaveInstructions(t *testing.T) {
	tokens := makeTokens(`name("a")instruction("something")`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)
	checkNodeCount(t, structogram.nodes, 1)
	instructionNode := structogram.nodes[0]
	if instructionNode.nodeType != "instruction" {
		t.Errorf("Wrong node type, expected %s, but got %s",
			"instruction", instructionNode.nodeType,
		)
	}
	if instructionNode.value != "something" {
		t.Errorf("Wrong node value, expected %s, but got %s",
			"something", instructionNode.value,
		)
	}
}

func TestStructogramsCanHaveMultipleInstructions(t *testing.T) {
	tokens := makeTokens(`name("a")instruction("b")instruction("c")`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)
	checkNodeCount(t, structogram.nodes, 2)
	instructionNode := structogram.nodes[0]
	if instructionNode.nodeType != "instruction" {
		t.Errorf("Wrong node type, expected %s, but got %s",
			"instruction", instructionNode.nodeType,
		)
	}
	if instructionNode.value != "b" {
		t.Errorf("Wrong node value, expected %s, but got %s",
			"b", instructionNode.value,
		)
	}
	instructionNode = structogram.nodes[1]
	if instructionNode.nodeType != "instruction" {
		t.Errorf("Wrong node type, expected %s, but got %s",
			"instruction", instructionNode.nodeType,
		)
	}
	if instructionNode.value != "c" {
		t.Errorf("Wrong node value, expected %s, but got %s",
			"c", instructionNode.value,
		)
	}
}

func TestParserCanHandleInvalidTokens(t *testing.T) {
	tokens := makeTokens(`name("a")asd`)
	_, err := parseTokens(tokens)
	// TODO it's not called 'identifier', it's actually a keyword
	checkErrorMsg(t, err, "1:10, expected 'identifier', but got 'invalid'")
}

func TestParserIgnoresWhitespaceTokens(t *testing.T) {
	tokens := makeTokens(`name("a")` + "\n " + `instruction("b")`)
	_, err := parseTokens(tokens)
	checkOk(t, err)
}

func TestIfTokenValuesAreEnclosedByParentheses(t *testing.T) {
	tokens := makeTokens(`name("a")if"b")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:12, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a")if("b"`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:16, expected 'closeParentheses', but got 'EOF'")
}

func TestIfTokenValueCanNotBeEmpty(t *testing.T) {
	tokens := makeTokens(`name("a")if()`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:13, expected 'string', but got 'closeParentheses'")
}

func TestIfTokenHasToHaveBody(t *testing.T) {
	tokens := makeTokens(`name("a")if("b")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:17, expected 'openBrace', but got 'EOF'")

	// The only valid tokens inside of an if-body are keywords or whitespace.
	// Whitespace should get entirely ignored, and anything that is not a
	// keyword, so either a string or EOF should cause an error.
	// The only exception are openParentheses, which are legal if they
	// are preceeded by a keyword
	tokens = makeTokens(`name("a")if("b"){`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:18, expected 'keyword', but got 'EOF'")

	tokens = makeTokens(`name("a")if("b"){"c"`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:18, expected 'keyword', but got 'string'")

	tokens = makeTokens(`name("a")if("b"){name}`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:18, expected 'keyword', but got 'name'")
}

func TestIfTokenCanHaveWhitespaceBetweenConditionAndBody(t *testing.T) {
	tokens := makeTokens(`name("a")if("b")` + "\n " + `{instruction("c")}`)
	_, err := parseTokens(tokens)
	checkOk(t, err)
}

func TestInstructionTokenInsideIfBodyBehavesTheSameAsOutside(t *testing.T) {
	tokens := makeTokens(`name("a") if("b") {instruction}`)
	_, err := parseTokens(tokens)
	checkErrorMsg(
		t, err, "1:31, expected 'openParentheses', but got 'closeBrace'",
	)

	tokens = makeTokens(`name("a") if("b") {instruction(}`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:32, expected 'string', but got 'closeBrace'")

	tokens = makeTokens(`name("a") if("b") {instruction("c"}`)
	_, err = parseTokens(tokens)
	checkErrorMsg(
		t, err, "1:35, expected 'closeParentheses', but got 'closeBrace'",
	)
}

func TestCanParseMultipleInstructionsInsideIfBody(t *testing.T) {
	tokens := makeTokens(`name("a") if("b") {instruction("c") instruction("d")}`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)

	checkNodeCount(t, structogram.nodes, 1)
	ifNode := structogram.nodes[0]
	if ifNode.nodeType != "if" {
		t.Errorf("Wronge node type, expected %s, but got %s", "if", ifNode.nodeType)
	}
	if ifNode.value != "b" {
		t.Errorf("Wrong node value, expected %s, but got %s", "b", ifNode.value)
	}

	ifBody := ifNode.nodes
	checkNodeCount(t, ifBody, 2)
	instructionNode := ifBody[0]
	if instructionNode.nodeType != "instruction" {
		t.Errorf("Wronge node type, expected %s, but got %s",
			"instruction", instructionNode.nodeType)
	}
	if instructionNode.value != "c" {
		t.Errorf("Wronge node value, expected %s, but got %s", "c", instructionNode.value)
	}
	instructionNode = ifBody[1]
	if instructionNode.nodeType != "instruction" {
		t.Errorf("Wronge node type, expected %s, but got %s",
			"instruction", instructionNode.nodeType)
	}
	if instructionNode.value != "d" {
		t.Errorf("Wronge node value, expected %s, but got %s", "d", instructionNode.value)
	}
}

func TestCanParseNestedIfs(t *testing.T) {
	tokens := makeTokens(`name("a") if("b") {if("c"){instruction("d")}}`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)

	checkNodeCount(t, structogram.nodes, 1)

	// first if
	ifNode := structogram.nodes[0]
	if ifNode.nodeType != "if" {
		t.Errorf("Wrong node type, expected %s, but got %s",
			"if",
			ifNode.nodeType,
		)
	}
	if ifNode.value != "b" {
		t.Errorf("Wrong node value, expected %s, but got %s", "b", ifNode.value)
	}
	ifBody := ifNode.nodes
	checkNodeCount(t, ifBody, 1)

	// nested if
	ifNode = ifBody[0]
	if ifNode.nodeType != "if" {
		t.Errorf("Wrong node type, expected %s, but got %s",
			"if",
			ifNode.nodeType,
		)
	}
	if ifNode.value != "c" {
		t.Errorf("Wrong node value, expected %s, but got %s", "c", ifNode.value)
	}
	ifBody = ifNode.nodes
	checkNodeCount(t, ifBody, 1)
	instructionNode := ifBody[0]
	if instructionNode.nodeType != "instruction" {
		t.Errorf(
			"Wrong node type, expected %s, but got %s",
			"instruction",
			instructionNode.nodeType,
		)
	}
	if instructionNode.value != "d" {
		t.Errorf("Wrong node value, expected %s, but got %s", "d", ifNode.value)
	}
}
