package editor

import "iter"

// Line represents a single line of text.
// Internally, the line uses a gap buffer to optimize insertions and deletions in the middle of the line.
// Line also is a doubly linked list node.
type Line struct {
	data     []rune
	gapStart int
	gapEnd   int
}

func NewEmptyLine(size int) *Line {
	return &Line{
		data:     make([]rune, size),
		gapStart: 0,
		gapEnd:   size,
	}
}

// NewLine creates a new Line with the given content. If cursorAtStart is true,
// the cursor will be placed at the start of the line, otherwise it will be placed
// at the end.
func NewLine(content []rune, cursorAtStart bool) *Line {
	contentLen := len(content)
	// TODO: optimize size based on content length
	size := contentLen * 2
	buf := make([]rune, size)
	var gapStart int
	var gapEnd int
	if cursorAtStart {
		gapStart = 0
		gapEnd = contentLen
		copy(buf[size-contentLen:], content)
	} else {
		gapStart = contentLen
		gapEnd = size
		copy(buf, content)
	}

	return &Line{
		data:     buf,
		gapStart: gapStart,
		gapEnd:   gapEnd,
	}
}

func (g *Line) Len() int {
	return len(g.data) - (g.gapEnd - g.gapStart)
}

func (g *Line) moveGapTo(pos int) {
	if pos < g.gapStart {
		// Move gap left
		for i := g.gapStart - 1; i >= pos; i-- {
			g.data[g.gapEnd-1] = g.data[i]
			g.gapEnd--
			g.gapStart--
		}
	} else if pos > g.gapStart {
		// Move gap right
		for i := g.gapStart; i < pos; i++ {
			g.data[g.gapStart] = g.data[g.gapEnd]
			g.gapEnd++
			g.gapStart++
		}
	}
}

func (g *Line) Insert(r rune) {
	if g.gapStart == g.gapEnd {
		g.grow()
	}
	g.data[g.gapStart] = r
	g.gapStart++
}

func (g *Line) Append(o *Line) {
	g.moveCursorTo(g.Len())
	for _, r := range o.Runes() {
		g.Insert(r)
	}
}

func (g *Line) DeleteBeforeCursor() {
	if g.gapStart > 0 {
		g.gapStart--
	}
}

func (g *Line) DeleteAfterCursor() {
	if g.gapEnd < len(g.data) {
		g.gapEnd++
	}
}

func (g *Line) CursorPos() int {
	return g.gapStart
}

func (g *Line) CutAfterCursor() []rune {
	if len(g.data)-g.gapEnd > 0 {
		res := g.data[g.gapEnd:]
		g.gapEnd = len(g.data)
		return res
	}
	return nil
}

func (g *Line) moveCursorDelta(delta int) {
	g.moveCursorTo(g.gapStart + delta)
}

func (g *Line) moveCursorTo(pos int) {
	if pos < 0 {
		pos = 0
	}
	if pos > g.Len() {
		pos = g.Len()
	}
	g.moveGapTo(pos)
}

func (g *Line) grow() {
	newBuf := make([]rune, len(g.data)*2)
	copy(newBuf, g.data[:g.gapStart])
	copy(newBuf[len(newBuf)-(len(g.data)-g.gapEnd):], g.data[g.gapEnd:])
	g.gapEnd = len(newBuf) - (len(g.data) - g.gapEnd)
	g.data = newBuf
}

func (g *Line) Runes() iter.Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		if g.gapStart == g.gapEnd {
			return
		}

		i := 0
		for i < g.gapStart {
			if !yield(i, g.data[i]) {
				return
			}
			i++
		}
		i = g.gapEnd
		l := len(g.data)
		delta := g.gapEnd - g.gapStart
		for i < l {
			if !yield(i-delta, g.data[i]) {
				return
			}
			i++
		}
	}
}
