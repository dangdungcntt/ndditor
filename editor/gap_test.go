package editor

import (
	"github.com/test-go/testify/require"
	"iter"
	"strings"
	"testing"
)

func visualizeGapBuffer(g *GapBuffer) string {
	s := strings.Builder{}
	for i, r := range g.buf {
		if i >= g.gapStart && i < g.gapEnd {
			s.WriteString("_")
		} else {
			s.WriteRune(r)
		}
	}

	return s.String()
}

type pair[K any, V any] struct {
	K K
	V V
}

func collect2Entries[K any, V any](seq iter.Seq2[K, V]) []pair[K, V] {
	var res []pair[K, V]
	for k, r := range seq {
		res = append(res, pair[K, V]{k, r})
	}
	return res
}

func TestGap(t *testing.T) {
	g := NewGapBuffer(5)
	g.Insert('a')
	g.Insert('b')
	g.Insert('c')
	g.Insert('d')

	require.Equal(t, 4, g.Len())
	require.Equal(t, "abcd_", visualizeGapBuffer(g))
	g.moveGapTo(0)
	require.Equal(t, "_abcd", visualizeGapBuffer(g))
	g.moveGapTo(1)
	require.Equal(t, "a_bcd", visualizeGapBuffer(g))
	g.Insert('e')
	require.Equal(t, "aebcd", visualizeGapBuffer(g))
	require.Equal(t, 2, g.gapStart)
	require.Equal(t, 2, g.gapEnd)
	g.Insert('f')
	require.Equal(t, "aef____bcd", visualizeGapBuffer(g))
	require.Equal(t, 3, g.gapStart)
	require.Equal(t, 7, g.gapEnd)
	g.moveCursorDelta(-1)
	g.DeleteBeforeCursor()
	require.Equal(t, "a_____fbcd", visualizeGapBuffer(g))
	require.Equal(t, "fbcd", string(g.CutAfterCursor()))
	require.Equal(t, "a_________", visualizeGapBuffer(g))
	require.Equal(t, "", string(g.CutAfterCursor()))
}

func TestGapRunes(t *testing.T) {
	g := NewGapBuffer(5)
	g.Insert('a')
	g.Insert('b')
	g.Insert('c')
	g.Insert('d')
	require.Equal(t, []pair[int, rune]{
		{0, 'a'}, {1, 'b'}, {2, 'c'}, {3, 'd'},
	}, collect2Entries(g.Runes()))
	g.moveGapTo(0)
	require.Equal(t, []pair[int, rune]{
		{0, 'a'}, {1, 'b'}, {2, 'c'}, {3, 'd'},
	}, collect2Entries(g.Runes()))
}
