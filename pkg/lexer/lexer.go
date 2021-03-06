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

package lexer

import (
	"errors"
	"unicode/utf8"

	"laptudirm.com/x/mash/pkg/token"
)

// Various error values returned by lexer.consume.
var (
	ErrNUL = errors.New("illegal character NUL")
	ErrBOM = errors.New("illegal byte order mark")
	ErrEnc = errors.New("illegal utf-8 encoding")
)

// lexer represents a mash source string and related lexing information.
type lexer struct {
	src string // source string
	ch  rune   // current character
	wd  int    // character width

	insertSemi bool

	Tokens TokenStream // lexer token channel

	err ErrorHandler // lexer errors handling function

	offset   int // offset of the start of the token in the source
	rdOffset int // offsen of the current rune in the source

	start token.Position // position in the source of the start of the token
	prev  token.Position // previous position in the source
	pos   token.Position // position in the source of the current rune

	ErrCount int // number of errors encountered
}

const (
	eof = -1     // end of file
	bom = 0xFEFF // byte order mark
)

// ErrorHandler is a function which accepts a position in the source and an
// error from the lexer and properly handles it.
//
type ErrorHandler func(token.Position, error)

type TokenStream chan token.Token

// Lex starts the lexing of src, using err to handle any lexer errors, and
// returns the lexer's token channel.
//
func Lex(src string, err ErrorHandler) TokenStream {
	origin := token.Position{
		Line: 1,
		Col:  1,
	}

	l := &lexer{
		src: src,

		Tokens: make(TokenStream),

		err: err,

		start: origin,
		pos:   origin,
	}
	go l.run()

	return l.Tokens
}

// emit emits a token of type t with the current position and literal to
// the lexer's token channel. It also resets the lexer position and offset
// variables.
func (l *lexer) emit(t token.Type) {
	l.Tokens <- token.Token{
		Type:     t,
		Literal:  l.literal(),
		Position: l.start,
	}

	l.ignore()
}

// error call's the lexer's error handler, if there is one, with the err
// and the current position, and increases the lexer's ErrorCount by 1.
//
func (l *lexer) error(err error) {
	l.ErrCount++
	if l.err != nil {
		l.err(l.pos, err)
	}
}

// peek returns the byte right after the current rune. It returns eof if
// there are no more bytes after the current rune.
//
func (l *lexer) peek() rune {
	if l.atEnd() {
		return eof
	}

	return rune(l.src[l.rdOffset])
}

// consume consumes the next rune, incresing rdOffset and pos by it's
// width, and sets ch to the consumed rune. It sets ch to eof if it is at
// the end of the source.
//
func (l *lexer) consume() {
	if l.atEnd() {
		l.ch = eof
		l.wd = 0
		return
	}

	r, w := rune(l.src[l.rdOffset]), 1
	if r == 0 {
		// null rune is illegal in source
		l.error(ErrNUL)
		goto advance
	}

	if r < utf8.RuneSelf {
		// r is a single byte rune
		goto advance
	}

	// decode multi-byte rune
	r, w = utf8.DecodeRuneInString(l.src[l.rdOffset:])

	if r == utf8.RuneError && w == 1 {
		// illegal unicode encoding
		l.error(ErrEnc)
		goto advance
	}

	// bom is only legal as the first rune
	if r == bom && l.offset > 0 {
		l.error(ErrBOM)
	}

advance:
	l.ch = r
	l.wd = w

	l.prev = l.pos

	l.rdOffset += w
	l.pos.Col += w

	if r == '\n' {
		l.pos.NextLine()
	}
}

func (l *lexer) backup() {
	l.rdOffset -= l.wd
	l.pos = l.prev
}

// literal returns a sub-string from the source from offset to rdOffset.
//
func (l *lexer) literal() string {
	return l.src[l.offset:l.rdOffset]
}

// ignore sets start to pos and offset to rdOffset.
//
func (l *lexer) ignore() {
	l.offset = l.rdOffset
	l.start = l.pos
}

// atEnd returns true if the rdOffset is greater than the length of the
// source.
func (l *lexer) atEnd() bool {
	return l.rdOffset >= len(l.src)
}
