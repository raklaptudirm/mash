// Copyright © 2021 Rak Laptudirm <raklaptudirm@gmail.com>
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

package lexer

import (
	"unicode"

	"github.com/raklaptudirm/mash/pkg/token"
)

type stateFunc func(*lexer) stateFunc

func (l *lexer) run() {
	for state := lexBase; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
}

func lexBase(l *lexer) stateFunc {
	r := l.peek()
	if unicode.IsSpace(r) {
		l.consumeSpace()
	}

	if isAlphabet(r) {
		l.consumeWord()

		word := l.literal()
		if token.IsKeyword(word) {
			l.emit(token.Lookup(word))
			return lexStmt
		}

		l.backup()
	}

	return lexCmd
}

func lexStmt(l *lexer) stateFunc {
	l.consume()

	switch {
	case unicode.IsSpace(l.ch):
		// ignore whitespace
		l.consumeSpace()

	// literals
	case isIdentStart(l.ch):
		// identifier
		l.consumeIdent()
		l.emit(token.Lookup(l.literal()))
	case unicode.IsDigit(l.ch):
		// number
		return lexNum
	case l.ch == '"':
		// format string
		l.consumeString()
		l.emit(token.STRING)

	// operators
	case token.IsOperator(string(l.ch)):
		return lexStmtOp

	// special
	case l.ch == '#':
		// line comment
		l.consumeComment()
		l.emit(token.COMMENT)
	case l.ch == eof:
		l.emit(token.EOF)
		return nil
	default:
		// rune not supported
		l.emit(token.ILLEGAL)
	}

	return lexStmt
}

func lexNum(l *lexer) stateFunc {
	for unicode.IsDigit(l.peek()) {
		l.consume()
	}

	l.emit(token.FLOAT)
	return lexStmt
}

func lexStmtOp(l *lexer) stateFunc {
	var t token.TokenType
	switch l.ch {
	case '+':
		t = token.ADD
	case '-':
		t = token.SUB
	case '*':
		t = token.MUL
	case '/':
		t = token.QUO
	case '%':
		t = token.REM
	case '&':
		t = token.AND
	case '|':
		t = token.OR
	case '^':
		t = token.XOR
	case '<':
		t = token.LSS
	case '>':
		t = token.GTR
	case '=':
		t = token.ASSIGN
	case '!':
		t = token.NOT
	case '(':
		t = token.LPAREN
	case '[':
		t = token.LPAREN
	case '{':
		t = token.LBRACE
	case ',':
		t = token.COMMA
	case ')':
		t = token.RPAREN
	case ']':
		t = token.RBRACK
	case '}':
		t = token.RBRACE
	case ';':
		t = token.SEMICOLON
	case ':':
		t = token.COLON
	}

	l.emit(t)
	return lexStmt
}

func lexCmd(l *lexer) stateFunc {
	return nil
}

func isAlphabet(r rune) bool {
	return r > 'A' && r < 'Z' || r > 'a' && r < 'z'
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdent(r rune) bool {
	return isIdentStart(r) || unicode.IsDigit(r)
}
