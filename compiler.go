package main

import (
	"fmt"
	"log"
	"strings"
)

type token struct {
	kind  string
	value string
}

func tokenize(input []rune) []token {
	tokens := []token{}

	for curr := 0; curr < len(input); curr++ {
		switch {
		case input[curr] == '(':
			tokens = append(tokens, token{"paren", "("})
		case input[curr] == ')':
			tokens = append(tokens, token{"paren", ")"})
		case input[curr] == ' ':

		case '0' <= input[curr] && input[curr] <= '9':
			value := []rune{input[curr]}
			for '0' <= input[curr+1] && input[curr+1] <= '9' {
				curr++
				value = append(value, input[curr])
			}
			tokens = append(tokens, token{"number", string(value)})

		case 'a' <= input[curr] && input[curr] <= 'z':
			value := []rune{input[curr]}
			for 'a' <= input[curr+1] && input[curr+1] <= 'z' {
				curr++
				value = append(value, input[curr])
			}
			tokens = append(tokens, token{"name", string(value)})

		default:
		}
	}

	return tokens
}

type node struct {
	kind   string
	value  string
	name   string
	body   []node
	params []node
}

type ast node

var pc int

func parse(tokens []token) ast {
	pc = 0

	ast := ast{
		kind: "Program",
		body: []node{},
	}

	for pc < len(tokens) {
		ast.body = append(ast.body, walk(tokens))
	}

	return ast
}

func walk(tokens []token) node {
	tok := tokens[pc]
	pc++

	switch {
	case tok.kind == "number":
		return node{kind: "NumberLiteral", value: tok.value}

	case tok.kind == "paren" && tok.value == "(":
		tok = tokens[pc]
		nod := node{
			kind:   "CallExpression",
			name:   tok.value,
			params: []node{},
		}

		pc++
		tok = tokens[pc]
		for tok.kind != "paren" || tok.value != ")" {
			nod.params = append(nod.params, walk(tokens))
			tok = tokens[pc]
		}
		pc++
		return nod
	}

	log.Fatal(tokens, " ", pc, " ", tok.kind)
	return node{}
}

func genCode(nod node) string {
	switch nod.kind {
	case "Program":
		var rst []string
		for _, child := range nod.body {
			rst = append(rst, genCode(child))
		}
		return strings.Join(rst, "\n")
	case "CallExpression":
		var rst []string
		for _, child := range nod.params {
			rst = append(rst, genCode(child))
		}
		r := strings.Join(rst, ", ")
		return nod.name + "(" + r + ")"
	case "NumberLiteral":
		return nod.value
	default:
		log.Fatal(nod.kind)
		return ""
	}
}

func compile(input string) string {
	tokens := tokenize([]rune(input))
	ast := parse(tokens)
	out := genCode(node(ast))
	return out
}

func main() {
	program := "(add 100 (subtract 110 600))"
	out := compile(program)
	fmt.Println(out)
}
