// Copyright 2019 Petr Homola. All rights reserved.
// Use of this source code is governed by the AGPL v3.0
// that can be found in the LICENSE file.

package textkit

import (
	"io/ioutil"
	"os"
	"strings"
)

type LexicalEntry struct {
	Lemma string
	Tag   string
}

type MorphologicalLexicon struct {
	entries map[string][]*LexicalEntry
}

func (lex *MorphologicalLexicon) AddEntry(form, lemma, tag string) {
	form = strings.ToLower(form)
	list := lex.entries[form]
	list = append(list, &LexicalEntry{lemma, tag})
	lex.entries[form] = list
}

func (lex *MorphologicalLexicon) Analyse(form string) []*LexicalEntry {
	return lex.entries[form]
}

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
