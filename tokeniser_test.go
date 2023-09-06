package textkit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEndIndents(t *testing.T) {
	r := require.New(t)

	tok := &Tokeniser{
		KeepEndIndents: true,
	}
	tokens := tok.Tokenise(`#
A
  B
C
  D
    E
  F
    G
H
  I
    J
`, "<file>")

	for _, tok := range tokens {
		t.Log(tok)
	}

	r.Equal(18, len(tokens))
}
