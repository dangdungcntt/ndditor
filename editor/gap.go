package editor

import "iter"

type GapBuffer struct {
	buf      []rune
	gapStart int
	gapEnd   int
}

func NewGapBuffer(size int) *GapBuffer {
	return &GapBuffer{
		buf:      make([]rune, size),
		gapStart: 0,
		gapEnd:   size,
	}
}

func NewGapBufferWithContent(content []rune, cursorPosAtStart bool) *GapBuffer {
	contentLen := len(content)
	// TODO: optimize size based on content length
	size := contentLen * 2
	buf := make([]rune, size)
	var gapStart int
	var gapEnd int
	if cursorPosAtStart {
		gapStart = 0
		gapEnd = contentLen
		copy(buf[size-contentLen:], content)
	} else {
		gapStart = contentLen
		gapEnd = size
		copy(buf, content)
	}

	return &GapBuffer{
		buf:      buf,
		gapStart: gapStart,
		gapEnd:   gapEnd,
	}
}

func (g *GapBuffer) Len() int {
	return len(g.buf) - (g.gapEnd - g.gapStart)
}

func (g *GapBuffer) moveGapTo(pos int) {
	if pos < g.gapStart {
		// Move gap left
		for i := g.gapStart - 1; i >= pos; i-- {
			g.buf[g.gapEnd-1] = g.buf[i]
			g.gapEnd--
			g.gapStart--
		}
	} else if pos > g.gapStart {
		// Move gap right
		for i := g.gapStart; i < pos; i++ {
			g.buf[g.gapStart] = g.buf[g.gapEnd]
			g.gapEnd++
			g.gapStart++
		}
	}
}

func (g *GapBuffer) Insert(r rune) {
	if g.gapStart == g.gapEnd {
		g.grow()
	}
	g.buf[g.gapStart] = r
	g.gapStart++
}

func (g *GapBuffer) Append(o *GapBuffer) {
	g.moveCursor(g.Len())
	for _, r := range o.Runes() {
		g.Insert(r)
	}
}

func (g *GapBuffer) DeleteBeforeCursor() {
	if g.gapStart > 0 {
		g.gapStart--
	}
}

func (g *GapBuffer) DeleteAfterCursor() {
	if g.gapEnd < len(g.buf) {
		g.gapEnd++
	}
}

func (g *GapBuffer) CursorPos() int {
	return g.gapStart
}

func (g *GapBuffer) CutAfterCursor() []rune {
	if len(g.buf)-g.gapEnd > 0 {
		res := g.buf[g.gapEnd:]
		g.gapEnd = len(g.buf)
		return res
	}
	return nil
}

func (g *GapBuffer) moveCursorDelta(delta int) {
	g.moveCursor(g.gapStart + delta)
}

func (g *GapBuffer) moveCursor(pos int) {
	if pos < 0 {
		pos = 0
	}
	if pos > g.Len() {
		pos = g.Len()
	}
	g.moveGapTo(pos)
}

func (g *GapBuffer) grow() {
	newBuf := make([]rune, len(g.buf)*2)
	copy(newBuf, g.buf[:g.gapStart])
	copy(newBuf[len(newBuf)-(len(g.buf)-g.gapEnd):], g.buf[g.gapEnd:])
	g.gapEnd = len(newBuf) - (len(g.buf) - g.gapEnd)
	g.buf = newBuf
}

func (g *GapBuffer) Runes() iter.Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		if g.gapStart == g.gapEnd {
			return
		}

		i := 0
		for i < g.gapStart {
			if !yield(i, g.buf[i]) {
				return
			}
			i++
		}
		i = g.gapEnd
		l := len(g.buf)
		delta := g.gapEnd - g.gapStart
		for i < l {
			if !yield(i-delta, g.buf[i]) {
				return
			}
			i++
		}
	}
}
