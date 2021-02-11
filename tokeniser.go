// Copyright 2019-2020 Petr Homola. All rights reserved.
// Use of this source code is governed by the AGPL v3.0
// that can be found in the LICENSE file.

package textkit

import (
	"strings"
)

// A token type.
type TokenType int

const (
	Word TokenType = iota
	Number
	String
	Symbol
	EOF
)

// A token.
type Token struct {
	// The token's type.
	Type TokenType
	// The form of the token as a string.
	Form string
	// The line where the token is located.
	Line int
	// The column where the token is located.
	Column int
	// An associated tag.
	Tag string
}

// A tokeniser which takes into account comments and special characters in identifiers.
type Tokeniser struct {
	CommentPrefix string
	StringChar    byte
	WordChars     string
}

func isWhiteChar(c byte) bool {
	return c == ' ' || c == '\r' || c == '\n' || c == '\t'
}

func (t *Tokeniser) isAlpha(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= 128 || strings.IndexByte(t.WordChars, c) != -1
}

func isNum(c byte) bool {
	return c >= '0' && c <= '9'
}

const (
	global = iota
	word
	number
	qstring
)

// Tokenises a text.
func (t *Tokeniser) Tokenise(s string) []*Token {
	var tokens []*Token
	i, line, col, colstart, state, numtag := 0, 1, 1, 1, global, ""
	var sb strings.Builder
	for {
		if state == global {
			for i < len(s) {
				c := s[i]
				if len(t.CommentPrefix) > 0 {
					if (len(s) - i) >= len(t.CommentPrefix) {
						if s[i:i+len(t.CommentPrefix)] == t.CommentPrefix {
							for i < len(s) && s[i] != '\n' {
								i++
							}
							c = '\n'
						}
					}
				}
				if !isWhiteChar(c) {
					break
				}
				if c == '\n' {
					line++
					col = 1
				} else {
					col++
				}
				i++
			}
		}
		if i == len(s) {
			break
		}
		c := s[i]
		switch state {
		case word:
			if t.isAlpha(c) || isNum(c) {
				if numtag == "" {
					sb.WriteByte(c)
				} else {
					numtag += string(c)
				}
				col++
				i++
			} else {
				if numtag == "" {
					tokens = append(tokens, &Token{Word, sb.String(), line, colstart, ""})
				} else {
					tokens = append(tokens, &Token{Number, sb.String(), line, colstart, numtag})
				}
				state = global
			}
		case number:
			if isNum(c) {
				sb.WriteByte(c)
				col++
				i++
			} else {
				if t.isAlpha(c) {
					numtag += string(c)
					col++
					i++
					state = word
				} else {
					tokens = append(tokens, &Token{Number, sb.String(), line, colstart, ""})
					state = global
				}
			}
		case qstring:
			if c == t.StringChar {
				tokens = append(tokens, &Token{String, sb.String(), line, colstart, ""})
				state = global
				col++
				i++
			} else {
				sb.WriteByte(c)
				if c == '\n' {
					line++
					col = 1
				} else {
					col++
				}
				i++
			}
		case global:
			if t.isAlpha(c) {
				state = word
				colstart = col
				numtag = ""
				sb.Reset()
				sb.WriteByte(c)
				col++
				i++
			} else if isNum(c) {
				state = number
				colstart = col
				numtag = ""
				sb.Reset()
				sb.WriteByte(c)
				col++
				i++
			} else if c == t.StringChar {
				sb.Reset()
				state = qstring
				colstart = col
				numtag = ""
				col++
				i++
			} else {
				tokens = append(tokens, &Token{Symbol, string([]byte{c}), line, col, ""})
				col++
				i++
			}
		}
	}
	switch state {
	case word:
		if numtag == "" {
			tokens = append(tokens, &Token{Word, sb.String(), line, colstart, ""})
		} else {
			tokens = append(tokens, &Token{Number, sb.String(), line, colstart, numtag})
		}
	case number:
		tokens = append(tokens, &Token{Number, sb.String(), line, colstart, ""})
	case qstring:
		tokens = append(tokens, &Token{String, sb.String(), line, colstart, ""})
	}
	tokens = append(tokens, &Token{EOF, "", line, col, ""})
	return tokens
}
