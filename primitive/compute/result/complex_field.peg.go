// Code generated by peg. DO NOT EDIT.
package result

//go:generate peg -inline ./api/compute/result/complex_field.peg

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleComplexField
	rulearray
	rulearray_contents
	ruleitem
	rulestring
	ruledquote_string
	rulesquote_string
	rulevalue
	rulews
	rulecomma
	rulelf
	rulecr
	ruleescdquote
	ruleescsquote
	rulesquote
	ruleobracket
	rulecbracket
	ruleoparen
	rulecparen
	rulenumber
	rulenegative
	ruledecimal_point
	ruletextdata
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	rulePegText
	ruleAction4
	ruleAction5
	ruleAction6
)

var rul3s = [...]string{
	"Unknown",
	"ComplexField",
	"array",
	"array_contents",
	"item",
	"string",
	"dquote_string",
	"squote_string",
	"value",
	"ws",
	"comma",
	"lf",
	"cr",
	"escdquote",
	"escsquote",
	"squote",
	"obracket",
	"cbracket",
	"oparen",
	"cparen",
	"number",
	"negative",
	"decimal_point",
	"textdata",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"PegText",
	"Action4",
	"Action5",
	"Action6",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type ComplexField struct {
	arrayElements

	Buffer string
	buffer []rune
	rules  [32]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *ComplexField) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *ComplexField) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *ComplexField
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *ComplexField) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *ComplexField) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.pushArray()
		case ruleAction1:
			p.popArray()
		case ruleAction2:
			p.pushArray()
		case ruleAction3:
			p.popArray()
		case ruleAction4:
			p.addElement(buffer[begin:end])
		case ruleAction5:
			p.addElement(buffer[begin:end])
		case ruleAction6:
			p.addElement(buffer[begin:end])

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *ComplexField) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 ComplexField <- <(array !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rulearray]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					if !matchDot() {
						goto l2
					}
					goto l0
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
				add(ruleComplexField, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 array <- <((ws* obracket Action0 array_contents cbracket Action1) / (ws* oparen Action2 array_contents comma? ws* cparen Action3))> */
		func() bool {
			position3, tokenIndex3 := position, tokenIndex
			{
				position4 := position
				{
					position5, tokenIndex5 := position, tokenIndex
				l7:
					{
						position8, tokenIndex8 := position, tokenIndex
						if !_rules[rulews]() {
							goto l8
						}
						goto l7
					l8:
						position, tokenIndex = position8, tokenIndex8
					}
					if !_rules[ruleobracket]() {
						goto l6
					}
					{
						add(ruleAction0, position)
					}
					if !_rules[rulearray_contents]() {
						goto l6
					}
					if !_rules[rulecbracket]() {
						goto l6
					}
					{
						add(ruleAction1, position)
					}
					goto l5
				l6:
					position, tokenIndex = position5, tokenIndex5
				l11:
					{
						position12, tokenIndex12 := position, tokenIndex
						if !_rules[rulews]() {
							goto l12
						}
						goto l11
					l12:
						position, tokenIndex = position12, tokenIndex12
					}
					if !_rules[ruleoparen]() {
						goto l3
					}
					{
						add(ruleAction2, position)
					}
					if !_rules[rulearray_contents]() {
						goto l3
					}
					{
						position14, tokenIndex14 := position, tokenIndex
						if !_rules[rulecomma]() {
							goto l14
						}
						goto l15
					l14:
						position, tokenIndex = position14, tokenIndex14
					}
				l15:
				l16:
					{
						position17, tokenIndex17 := position, tokenIndex
						if !_rules[rulews]() {
							goto l17
						}
						goto l16
					l17:
						position, tokenIndex = position17, tokenIndex17
					}
					if !_rules[rulecparen]() {
						goto l3
					}
					{
						add(ruleAction3, position)
					}
				}
			l5:
				add(rulearray, position4)
			}
			return true
		l3:
			position, tokenIndex = position3, tokenIndex3
			return false
		},
		/* 2 array_contents <- <(ws* (item ws* (comma ws* item ws*)*)?)> */
		func() bool {
			{
				position20 := position
			l21:
				{
					position22, tokenIndex22 := position, tokenIndex
					if !_rules[rulews]() {
						goto l22
					}
					goto l21
				l22:
					position, tokenIndex = position22, tokenIndex22
				}
				{
					position23, tokenIndex23 := position, tokenIndex
					if !_rules[ruleitem]() {
						goto l23
					}
				l25:
					{
						position26, tokenIndex26 := position, tokenIndex
						if !_rules[rulews]() {
							goto l26
						}
						goto l25
					l26:
						position, tokenIndex = position26, tokenIndex26
					}
				l27:
					{
						position28, tokenIndex28 := position, tokenIndex
						if !_rules[rulecomma]() {
							goto l28
						}
					l29:
						{
							position30, tokenIndex30 := position, tokenIndex
							if !_rules[rulews]() {
								goto l30
							}
							goto l29
						l30:
							position, tokenIndex = position30, tokenIndex30
						}
						if !_rules[ruleitem]() {
							goto l28
						}
					l31:
						{
							position32, tokenIndex32 := position, tokenIndex
							if !_rules[rulews]() {
								goto l32
							}
							goto l31
						l32:
							position, tokenIndex = position32, tokenIndex32
						}
						goto l27
					l28:
						position, tokenIndex = position28, tokenIndex28
					}
					goto l24
				l23:
					position, tokenIndex = position23, tokenIndex23
				}
			l24:
				add(rulearray_contents, position20)
			}
			return true
		},
		/* 3 item <- <(array / string / (<value> Action4))> */
		func() bool {
			position33, tokenIndex33 := position, tokenIndex
			{
				position34 := position
				{
					position35, tokenIndex35 := position, tokenIndex
					if !_rules[rulearray]() {
						goto l36
					}
					goto l35
				l36:
					position, tokenIndex = position35, tokenIndex35
					{
						position38 := position
						{
							position39, tokenIndex39 := position, tokenIndex
							{
								position41 := position
								if !_rules[ruleescdquote]() {
									goto l40
								}
								{
									position42 := position
								l43:
									{
										position44, tokenIndex44 := position, tokenIndex
										{
											position45, tokenIndex45 := position, tokenIndex
											if !_rules[ruletextdata]() {
												goto l46
											}
											goto l45
										l46:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[rulesquote]() {
												goto l47
											}
											goto l45
										l47:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[rulelf]() {
												goto l48
											}
											goto l45
										l48:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[rulecr]() {
												goto l49
											}
											goto l45
										l49:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[ruleobracket]() {
												goto l50
											}
											goto l45
										l50:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[rulecbracket]() {
												goto l51
											}
											goto l45
										l51:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[ruleoparen]() {
												goto l52
											}
											goto l45
										l52:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[rulecparen]() {
												goto l53
											}
											goto l45
										l53:
											position, tokenIndex = position45, tokenIndex45
											if !_rules[rulecomma]() {
												goto l44
											}
										}
									l45:
										goto l43
									l44:
										position, tokenIndex = position44, tokenIndex44
									}
									add(rulePegText, position42)
								}
								if !_rules[ruleescdquote]() {
									goto l40
								}
								{
									add(ruleAction5, position)
								}
								add(ruledquote_string, position41)
							}
							goto l39
						l40:
							position, tokenIndex = position39, tokenIndex39
							{
								position55 := position
								if !_rules[rulesquote]() {
									goto l37
								}
								{
									position56 := position
								l57:
									{
										position58, tokenIndex58 := position, tokenIndex
										{
											position59, tokenIndex59 := position, tokenIndex
											{
												position61 := position
												if buffer[position] != rune('\\') {
													goto l60
												}
												position++
												if buffer[position] != rune('\'') {
													goto l60
												}
												position++
												add(ruleescsquote, position61)
											}
											goto l59
										l60:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[ruleescdquote]() {
												goto l62
											}
											goto l59
										l62:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[ruletextdata]() {
												goto l63
											}
											goto l59
										l63:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[rulelf]() {
												goto l64
											}
											goto l59
										l64:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[rulecr]() {
												goto l65
											}
											goto l59
										l65:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[ruleobracket]() {
												goto l66
											}
											goto l59
										l66:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[rulecbracket]() {
												goto l67
											}
											goto l59
										l67:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[ruleoparen]() {
												goto l68
											}
											goto l59
										l68:
											position, tokenIndex = position59, tokenIndex59
											if !_rules[rulecparen]() {
												goto l58
											}
										}
									l59:
										goto l57
									l58:
										position, tokenIndex = position58, tokenIndex58
									}
									add(rulePegText, position56)
								}
								if !_rules[rulesquote]() {
									goto l37
								}
								{
									add(ruleAction6, position)
								}
								add(rulesquote_string, position55)
							}
						}
					l39:
						add(rulestring, position38)
					}
					goto l35
				l37:
					position, tokenIndex = position35, tokenIndex35
					{
						position70 := position
						{
							position71 := position
							{
								position72, tokenIndex72 := position, tokenIndex
								{
									position74 := position
									if buffer[position] != rune('-') {
										goto l72
									}
									position++
									add(rulenegative, position74)
								}
								goto l73
							l72:
								position, tokenIndex = position72, tokenIndex72
							}
						l73:
							if !_rules[rulenumber]() {
								goto l33
							}
						l75:
							{
								position76, tokenIndex76 := position, tokenIndex
								if !_rules[rulenumber]() {
									goto l76
								}
								goto l75
							l76:
								position, tokenIndex = position76, tokenIndex76
							}
							{
								position77, tokenIndex77 := position, tokenIndex
								{
									position79 := position
									if buffer[position] != rune('.') {
										goto l77
									}
									position++
									add(ruledecimal_point, position79)
								}
								if !_rules[rulenumber]() {
									goto l77
								}
							l80:
								{
									position81, tokenIndex81 := position, tokenIndex
									if !_rules[rulenumber]() {
										goto l81
									}
									goto l80
								l81:
									position, tokenIndex = position81, tokenIndex81
								}
								goto l78
							l77:
								position, tokenIndex = position77, tokenIndex77
							}
						l78:
							add(rulevalue, position71)
						}
						add(rulePegText, position70)
					}
					{
						add(ruleAction4, position)
					}
				}
			l35:
				add(ruleitem, position34)
			}
			return true
		l33:
			position, tokenIndex = position33, tokenIndex33
			return false
		},
		/* 4 string <- <(dquote_string / squote_string)> */
		nil,
		/* 5 dquote_string <- <(escdquote <(textdata / squote / lf / cr / obracket / cbracket / oparen / cparen / comma)*> escdquote Action5)> */
		nil,
		/* 6 squote_string <- <(squote <(escsquote / escdquote / textdata / lf / cr / obracket / cbracket / oparen / cparen)*> squote Action6)> */
		nil,
		/* 7 value <- <(negative? number+ (decimal_point number+)?)> */
		nil,
		/* 8 ws <- <' '> */
		func() bool {
			position87, tokenIndex87 := position, tokenIndex
			{
				position88 := position
				if buffer[position] != rune(' ') {
					goto l87
				}
				position++
				add(rulews, position88)
			}
			return true
		l87:
			position, tokenIndex = position87, tokenIndex87
			return false
		},
		/* 9 comma <- <','> */
		func() bool {
			position89, tokenIndex89 := position, tokenIndex
			{
				position90 := position
				if buffer[position] != rune(',') {
					goto l89
				}
				position++
				add(rulecomma, position90)
			}
			return true
		l89:
			position, tokenIndex = position89, tokenIndex89
			return false
		},
		/* 10 lf <- <'\n'> */
		func() bool {
			position91, tokenIndex91 := position, tokenIndex
			{
				position92 := position
				if buffer[position] != rune('\n') {
					goto l91
				}
				position++
				add(rulelf, position92)
			}
			return true
		l91:
			position, tokenIndex = position91, tokenIndex91
			return false
		},
		/* 11 cr <- <'\r'> */
		func() bool {
			position93, tokenIndex93 := position, tokenIndex
			{
				position94 := position
				if buffer[position] != rune('\r') {
					goto l93
				}
				position++
				add(rulecr, position94)
			}
			return true
		l93:
			position, tokenIndex = position93, tokenIndex93
			return false
		},
		/* 12 escdquote <- <'"'> */
		func() bool {
			position95, tokenIndex95 := position, tokenIndex
			{
				position96 := position
				if buffer[position] != rune('"') {
					goto l95
				}
				position++
				add(ruleescdquote, position96)
			}
			return true
		l95:
			position, tokenIndex = position95, tokenIndex95
			return false
		},
		/* 13 escsquote <- <('\\' '\'')> */
		nil,
		/* 14 squote <- <'\''> */
		func() bool {
			position98, tokenIndex98 := position, tokenIndex
			{
				position99 := position
				if buffer[position] != rune('\'') {
					goto l98
				}
				position++
				add(rulesquote, position99)
			}
			return true
		l98:
			position, tokenIndex = position98, tokenIndex98
			return false
		},
		/* 15 obracket <- <'['> */
		func() bool {
			position100, tokenIndex100 := position, tokenIndex
			{
				position101 := position
				if buffer[position] != rune('[') {
					goto l100
				}
				position++
				add(ruleobracket, position101)
			}
			return true
		l100:
			position, tokenIndex = position100, tokenIndex100
			return false
		},
		/* 16 cbracket <- <']'> */
		func() bool {
			position102, tokenIndex102 := position, tokenIndex
			{
				position103 := position
				if buffer[position] != rune(']') {
					goto l102
				}
				position++
				add(rulecbracket, position103)
			}
			return true
		l102:
			position, tokenIndex = position102, tokenIndex102
			return false
		},
		/* 17 oparen <- <'('> */
		func() bool {
			position104, tokenIndex104 := position, tokenIndex
			{
				position105 := position
				if buffer[position] != rune('(') {
					goto l104
				}
				position++
				add(ruleoparen, position105)
			}
			return true
		l104:
			position, tokenIndex = position104, tokenIndex104
			return false
		},
		/* 18 cparen <- <')'> */
		func() bool {
			position106, tokenIndex106 := position, tokenIndex
			{
				position107 := position
				if buffer[position] != rune(')') {
					goto l106
				}
				position++
				add(rulecparen, position107)
			}
			return true
		l106:
			position, tokenIndex = position106, tokenIndex106
			return false
		},
		/* 19 number <- <([a-z] / [A-Z] / [0-9])> */
		func() bool {
			position108, tokenIndex108 := position, tokenIndex
			{
				position109 := position
				{
					position110, tokenIndex110 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex = position110, tokenIndex110
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l112
					}
					position++
					goto l110
				l112:
					position, tokenIndex = position110, tokenIndex110
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l108
					}
					position++
				}
			l110:
				add(rulenumber, position109)
			}
			return true
		l108:
			position, tokenIndex = position108, tokenIndex108
			return false
		},
		/* 20 negative <- <'-'> */
		nil,
		/* 21 decimal_point <- <'.'> */
		nil,
		/* 22 textdata <- <([a-z] / [A-Z] / [0-9] / ' ' / '!' / '#' / '$' / '&' / '%' / '*' / '+' / '-' / '.' / '/' / ':' / ';' / [<->] / '?' / '\\' / '^' / '_' / '`' / '{' / '|' / '}' / '~')> */
		func() bool {
			position115, tokenIndex115 := position, tokenIndex
			{
				position116 := position
				{
					position117, tokenIndex117 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex = position117, tokenIndex117
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l119
					}
					position++
					goto l117
				l119:
					position, tokenIndex = position117, tokenIndex117
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l120
					}
					position++
					goto l117
				l120:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune(' ') {
						goto l121
					}
					position++
					goto l117
				l121:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('!') {
						goto l122
					}
					position++
					goto l117
				l122:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('#') {
						goto l123
					}
					position++
					goto l117
				l123:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('$') {
						goto l124
					}
					position++
					goto l117
				l124:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('&') {
						goto l125
					}
					position++
					goto l117
				l125:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('%') {
						goto l126
					}
					position++
					goto l117
				l126:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('*') {
						goto l127
					}
					position++
					goto l117
				l127:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('+') {
						goto l128
					}
					position++
					goto l117
				l128:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('-') {
						goto l129
					}
					position++
					goto l117
				l129:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('.') {
						goto l130
					}
					position++
					goto l117
				l130:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('/') {
						goto l131
					}
					position++
					goto l117
				l131:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune(':') {
						goto l132
					}
					position++
					goto l117
				l132:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune(';') {
						goto l133
					}
					position++
					goto l117
				l133:
					position, tokenIndex = position117, tokenIndex117
					if c := buffer[position]; c < rune('<') || c > rune('>') {
						goto l134
					}
					position++
					goto l117
				l134:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('?') {
						goto l135
					}
					position++
					goto l117
				l135:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('\\') {
						goto l136
					}
					position++
					goto l117
				l136:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('^') {
						goto l137
					}
					position++
					goto l117
				l137:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('_') {
						goto l138
					}
					position++
					goto l117
				l138:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('`') {
						goto l139
					}
					position++
					goto l117
				l139:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('{') {
						goto l140
					}
					position++
					goto l117
				l140:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('|') {
						goto l141
					}
					position++
					goto l117
				l141:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('}') {
						goto l142
					}
					position++
					goto l117
				l142:
					position, tokenIndex = position117, tokenIndex117
					if buffer[position] != rune('~') {
						goto l115
					}
					position++
				}
			l117:
				add(ruletextdata, position116)
			}
			return true
		l115:
			position, tokenIndex = position115, tokenIndex115
			return false
		},
		/* 24 Action0 <- <{ p.pushArray() }> */
		nil,
		/* 25 Action1 <- <{ p.popArray() }> */
		nil,
		/* 26 Action2 <- <{ p.pushArray() }> */
		nil,
		/* 27 Action3 <- <{ p.popArray() }> */
		nil,
		nil,
		/* 29 Action4 <- <{ p.addElement(buffer[begin:end]) }> */
		nil,
		/* 30 Action5 <- <{ p.addElement(buffer[begin:end]) }> */
		nil,
		/* 31 Action6 <- <{ p.addElement(buffer[begin:end]) }> */
		nil,
	}
	p.rules = _rules
}