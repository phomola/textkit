// Copyright 2019-2020 Petr Homola. All rights reserved.
// Use of this source code is governed by the AGPL v3.0
// that can be found in the LICENSE file.

// Package textkit provides a tokeniser and morphological analyser.
package textkit

import (
	"io/ioutil"
	"os"
	"strings"
)

// LexicalEntry is a lexical entry.
type LexicalEntry struct {
	// The entry's lemma.
	Lemma string
	// The morphological tag.
	Tag string
}

// MorphologicalLexicon is a morphological lexicon.
type MorphologicalLexicon struct {
	entries map[string][]*LexicalEntry
}

// AddEntry adds an entry to the lexicon.
func (lex *MorphologicalLexicon) AddEntry(form, lemma, tag string) {
	form = strings.ToLower(form)
	list := lex.entries[form]
	list = append(list, &LexicalEntry{lemma, tag})
	lex.entries[form] = list
}

// Analyse analyses a word form. Returns a list of lexical entries.
func (lex *MorphologicalLexicon) Analyse(form string) []*LexicalEntry {
	return lex.entries[form]
}

// NewMorphologicalLexicon returns a new morphological lexicon.
func NewMorphologicalLexicon(filename string) (*MorphologicalLexicon, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	lex := &MorphologicalLexicon{make(map[string][]*LexicalEntry)}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line != "" {
			comps := strings.Split(line, "\t")
			lex.AddEntry(comps[0], comps[1], comps[2])
		}
	}
	return lex, nil
}
