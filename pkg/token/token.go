// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package token

import (
	"strconv"
	"unicode"
)

// Type represents the type of a token which will be emitted by the lexer.
type Type int

// Various types of tokens emitted by the lexer.
const (
	// Special tokens
	Illegal Type = iota
	Eof
	Comment

	literalBeg
	// Identifiers and basic type literals
	Identifier // main
	Number     // 3.14
	String     // "abc"
	literalEnd

	operatorBeg
	// Operators and delimiters
	Addition       // +
	Subtraction    // -
	Multiplication // *
	Quotient       // /
	Remainder      // %

	And        // &
	Or         // |
	Xor        // ^
	ShiftLeft  // <<
	ShiftRight // >>
	AndNot     // &^

	AdditionAssign       // +=
	SubtractionAssign    // -=
	MultiplicationAssign // *=
	QuotientAssign       // /=
	RemainderAssign      // %=

	AndAssign        // &=
	OrAssign         // |=
	XorAssign        // ^=
	ShiftLeftAssign  // <<=
	ShiftRightAssign // >>=
	AndNotAssign     // &^=

	LogicalAnd // &&
	LogicalOr  // ||

	Equal       // ==
	LessThan    // <
	GreaterThan // >
	Assign      // =
	Define      // :=
	Not         // !

	NotEqual         // !=
	LessThanEqual    // <=
	GreaterThanEqual // >=

	LeftParen // (
	LeftBrack // [
	LeftBrace // {
	Template  // '
	Comma     // ,
	Period    // .

	RightParen // )
	RightBrack // ]
	RightBrace // }
	Semicolon  // ;
	Colon      // :
	operatorEnd

	keywordBeg
	// Keywords
	For
	If
	Else

	Let
	Obj
	Func

	Break
	Continue
	Return
	keywordEnd
)

var tokens = [...]string{
	Illegal: "ILLEGAL",

	Eof:     "EOF",
	Comment: "COMMENT",

	Identifier: "IDENT",
	Number:     "FLOAT",
	String:     "STRING",

	Addition:       "+",
	Subtraction:    "-",
	Multiplication: "*",
	Quotient:       "/",
	Remainder:      "%",

	And:        "&",
	Or:         "|",
	Xor:        "^",
	ShiftLeft:  "<<",
	ShiftRight: ">>",
	AndNot:     "&^",

	AdditionAssign:       "+=",
	SubtractionAssign:    "-=",
	MultiplicationAssign: "*=",
	QuotientAssign:       "/=",
	RemainderAssign:      "%=",

	AndAssign:        "&=",
	OrAssign:         "|=",
	XorAssign:        "^=",
	ShiftLeftAssign:  "<<=",
	ShiftRightAssign: ">>=",
	AndNotAssign:     "&^=",

	LogicalAnd: "&&",
	LogicalOr:  "||",

	Equal:       "==",
	LessThan:    "<",
	GreaterThan: ">",
	Assign:      "=",
	Not:         "!",

	NotEqual:         "!=",
	LessThanEqual:    "<=",
	GreaterThanEqual: ">=",
	Define:           ":=",

	LeftParen: "(",
	LeftBrack: "[",
	LeftBrace: "{",
	Template:  "'",
	Comma:     ",",
	Period:    ".",

	RightParen: ")",
	RightBrack: "]",
	RightBrace: "}",
	Semicolon:  ";",
	Colon:      ":",

	For:  "for",
	If:   "if",
	Else: "else",

	Let:  "let",
	Obj:  "obj",
	Func: "func",

	Break:    "break",
	Continue: "continue",
	Return:   "return",
}

func token(s string) Type {
	for t, val := range tokens {
		if val == s {
			return Type(t)
		}
	}

	return Illegal
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Type) String() string {
	s := ""
	if 0 <= tok && tok < Type(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// InsertSemi returns a boolean depending on wether a semicolon
// should be inserted after a token of type tok. It returns true if
// a semicolon should be inserted, and false if should not.
//
func (tok Type) InsertSemi() bool {
	if tok.IsLiteral() {
		return true
	}

	switch tok {
	case RightParen, RightBrack, RightBrace, Break, Continue, Return:
		return true
	default:
		return false
	}
}

// IsLiteral returns a boolean depending on wether the type of tok is
// a valid literal. Literals are tokens of with a value greater than
// literal_beg but less than literal_end.
//
func (tok Type) IsLiteral() bool {
	return literalBeg < tok && tok < literalEnd
}

// IsOperator returns a boolean depending on wether the type of tok is
// a valid operator. Operators are tokens of with a value greater than
// operator_beg but less than operator_end.
//
func (tok Type) IsOperator() bool {
	return operatorBeg < tok && tok < operatorEnd
}

// IsKeyword returns a boolean depending on wether the type of tok is
// a valid keyword. Keywords are tokens of with a value greater than
// keyword_beg but less than keyword_end.
//
func (tok Type) IsKeyword() bool {
	return keywordBeg < tok && tok < keywordEnd
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
}

// IsKeyword returns a boolean depending on wether name is a valid
// keyword. A string is a keyword if it is present in the keywords
// map.
func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

// IsIdentifier returns a boolean depending of wether name is a valid
// identifier. A string is a valid identifier if it's first letter is
// an unicode letter(gc = L) or an underscore, while the rest of the
// characters are letters, underscores, or decimal digits(0-9).
//
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && (i == 0 || !unicode.IsDigit(c)) && c != '_' {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}

// IsOperator returns a boolean depending on wether name is a valid
// operator or not. If the string belongs in the list of mash operators,
// it is a valid operator.
func IsOperator(s string) bool {
	t := token(s)
	return t.IsOperator()
}

// Lookup checks if name is a keyword, and returns the token type of the
// keyword if it is. Otherwise, it returns IDENT.
func Lookup(name string) Type {
	if tok, ok := keywords[name]; ok {
		return tok
	}

	return Identifier
}

// Token represtents a single token which will be emitted by the lexer.
type Token struct {
	Type     Type     // type of the token
	Literal  string   // literal in source
	Position Position // position in source
}
