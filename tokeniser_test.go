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
abcd 1234
efgh 5678
  abcd 1234
  efgh 5678
abcd 1234
efgh 5678
`, "<file>")

	// for _, tok := range tokens {
	// 	t.Log(tok)
	// }

	r.Equal(15, len(tokens))
}
