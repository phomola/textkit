// Copyright 2019-2020 Petr Homola. All rights reserved.
// Use of this source code is governed by the AGPL v3.0
// that can be found in the LICENSE file.

package textkit

import (
	"fmt"
	"strings"
)

// A token type.
type TokenType int

const (
	Word TokenType = iota
	Number
	String
	Symbol
	EOL
	EOF
)

// Location specifies the location of a token.
type Location struct {
	File   string
	Line   int
	Column int
}

func (loc Location) String() string {
	return fmt.Sprintf("%s:%d", loc.File, loc.Line)
}

var _ fmt.Stringer = Location{}

// A token.
type Token struct {
	// The token's type.
	Type TokenType
	// The form of the token as a slice of runes.
	Form []rune
	// The location where the token is located.
	Loc Location
	// An associated tag.
	Tag string
}

// A tokeniser which takes into account comments and special characters in identifiers.
type Tokeniser struct {
	CommentPrefix string
	StringRune    rune
	IdentChars    string
	KeepEOLs      bool
}

func isWhiteChar(c rune) bool {
	return c == ' ' || c == '\r' || c == '\n' || c == '\t'
}

func (t *Tokeniser) isAlpha(c rune) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= 128 || strings.ContainsRune(t.IdentChars, c)
}

func isNum(c rune) bool {
	return c >= '0' && c <= '9'
}

const (
	global = iota
	word
	number
	qstring
)

// Tokenises a text.
func (t *Tokeniser) Tokenise(text, file string) []*Token {
	runes := []rune(text)
	commentPrefixRunes := []rune(t.CommentPrefix)
	var tokens []*Token
	i, line, col, colstart, state, numtag := 0, 1, 1, 1, global, ""
	//var sb strings.Builder
	var form []rune
	for {
		if state == global {
			for i < len(runes) {
				r := runes[i]
				if len(t.CommentPrefix) > 0 {
					if (len(runes) - i) >= len(commentPrefixRunes) {
						if string(runes[i:i+len(commentPrefixRunes)]) == t.CommentPrefix {
							for i < len(runes) && runes[i] != '\n' {
								i++
							}
							r = '\n'
						}
					}
				}
				if !isWhiteChar(r) {
					break
				}
				if r == '\n' {
					if t.KeepEOLs {
						tokens = append(tokens, &Token{EOL, nil, Location{file, line, col}, ""})
					}
					line++
					col = 1
				} else {
					col++
				}
				i++
			}
		}
		if i == len(runes) {
			break
		}
		r := runes[i]
		switch state {
		case word:
			if t.isAlpha(r) || isNum(r) {
				if numtag == "" {
					form = append(form, r) //sb.WriteRune(r)
				} else {
					numtag += string(r)
				}
				col++
				i++
			} else {
				if numtag == "" {
					tokens = append(tokens, &Token{Word, form, Location{file, line, colstart}, ""})
				} else {
					tokens = append(tokens, &Token{Number, form, Location{file, line, colstart}, numtag})
				}
				state = global
			}
		case number:
			if isNum(r) {
				form = append(form, r) //sb.WriteRune(r)
				col++
				i++
			} else {
				if t.isAlpha(r) {
					numtag += string(r)
					col++
					i++
					state = word
				} else {
					tokens = append(tokens, &Token{Number, form, Location{file, line, colstart}, ""})
					state = global
				}
			}
		case qstring:
			if r == t.StringRune {
				tokens = append(tokens, &Token{String, form, Location{file, line, colstart}, ""})
				state = global
				col++
				i++
			} else {
				form = append(form, r) //sb.WriteRune(r)
				if r == '\n' {
					line++
					col = 1
				} else {
					col++
				}
				i++
			}
		case global:
			if t.isAlpha(r) {
				state = word
				colstart = col
				numtag = ""
				form = nil             //sb.Reset()
				form = append(form, r) //sb.WriteRune(r)
				col++
				i++
			} else if isNum(r) {
				state = number
				colstart = col
				numtag = ""
				form = nil             //sb.Reset()
				form = append(form, r) //sb.WriteRune(r)
				col++
				i++
			} else if r == t.StringRune {
				form = nil //sb.Reset()
				state = qstring
				colstart = col
				numtag = ""
				col++
				i++
			} else {
				tokens = append(tokens, &Token{Symbol, []rune{r}, Location{file, line, col}, ""})
				col++
				i++
			}
		}
	}
	switch state {
	case word:
		if numtag == "" {
			tokens = append(tokens, &Token{Word, form, Location{file, line, colstart}, ""})
		} else {
			tokens = append(tokens, &Token{Number, form, Location{file, line, colstart}, numtag})
		}
	case number:
		tokens = append(tokens, &Token{Number, form, Location{file, line, colstart}, ""})
	case qstring:
		tokens = append(tokens, &Token{String, form, Location{file, line, colstart}, ""})
	}
	tokens = append(tokens, &Token{EOF, nil, Location{file, line, col}, ""})
	return tokens
}
