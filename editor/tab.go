package editor

import (
	"fmt"
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
)

var _ layout.Element = (*Tab)(nil)

// Tab represents the content of a Tab. It is a doubly linked list of lines.
type Tab struct {
	size      layout.Size
	name      string
	cursorPos layout.Point
	lineIndex int
	lines     []*Line
}

func NewTab(name string, lines ...*Line) *Tab {
	return &Tab{
		name:  name,
		lines: lines,
	}
}

func (s *Tab) InsertNewline() {
	line := s.lines[s.lineIndex]
	s.lineIndex++
	s.cursorPos.Y++
	if s.cursorPos.Y >= s.size.Height {
		s.cursorPos.Y = s.size.Height - 1
	}
	newLineContent := line.CutAfterCursor()
	var newLine *Line
	if len(newLineContent) == 0 {
		newLine = NewEmptyLine(64)
	} else {
		newLine = NewLine(newLineContent, true)
	}
	s.lines = append(s.lines, nil)
	if s.lineIndex == len(s.lines)-1 {
		s.lines[s.lineIndex] = newLine
	} else {
		copy(s.lines[s.lineIndex+1:], s.lines[s.lineIndex:])
		s.lines[s.lineIndex] = newLine
	}
	s.cursorPos.X = 0
}

func (s *Tab) InsertRune(r rune) {
	s.lines[s.lineIndex].Insert(r)
	s.cursorPos.X++
}

func (s *Tab) Backspace() {
	if s.cursorPos.X > 0 {
		s.lines[s.lineIndex].DeleteBeforeCursor()
		s.cursorPos.X--
	} else if s.lineIndex > 0 {
		aboveLine := s.lines[s.lineIndex-1]
		s.cursorPos.X = aboveLine.Len()
		aboveLine.Append(s.lines[s.lineIndex])
		aboveLine.moveCursorTo(s.cursorPos.X)
		copy(s.lines[s.lineIndex:], s.lines[s.lineIndex+1:])
		s.lines = s.lines[:len(s.lines)-1]
		s.lineIndex--
		s.cursorPos.Y--
		if s.cursorPos.Y < 0 {
			s.cursorPos.Y = 0
		}
	}
}

func (s *Tab) Delete() {
	if s.cursorPos.X < s.lines[s.lineIndex].Len() {
		s.lines[s.lineIndex].DeleteAfterCursor()
	} else if s.lineIndex < len(s.lines)-1 {
		s.lines[s.lineIndex].Append(s.lines[s.lineIndex+1])
		copy(s.lines[s.lineIndex+1:], s.lines[s.lineIndex+2:])
		s.lines = s.lines[:len(s.lines)-1]
	}
}

func (s *Tab) MoveCursor(dx, dy int) {
	maxCursorY := min(s.size.Height-1, s.cursorPos.Y+(len(s.lines)-s.lineIndex-1))
	s.lineIndex += dy
	if s.lineIndex < 0 {
		s.lineIndex = 0
	} else if s.lineIndex >= len(s.lines) {
		s.lineIndex = len(s.lines) - 1
	}
	s.cursorPos.Y += dy

	if s.cursorPos.Y < 0 {
		s.cursorPos.Y = 0
	} else if s.cursorPos.Y > maxCursorY {
		s.cursorPos.Y = maxCursorY
	}

	s.cursorPos.X += dx
	if s.cursorPos.X < 0 {
		s.cursorPos.X = 0
	}
	maxLen := s.lines[s.lineIndex].Len()
	if s.cursorPos.X > maxLen {
		s.cursorPos.X = maxLen
	}
	s.lines[s.lineIndex].moveCursorTo(s.cursorPos.X)
}

func (s *Tab) Render(screen tcell.Screen, mountPoint layout.Point) {
	minLine := s.lineIndex - s.cursorPos.Y
	maxLine := s.lineIndex + (s.size.Height - s.cursorPos.Y)

	showCursor := true
	screenLine := 0
	for y, line := range s.lines {
		if y < minLine || y >= maxLine {
			continue
		}
		for x, r := range line.Runes() {
			if x == s.cursorPos.X && y == s.lineIndex {
				showCursor = false
				screen.SetContent(mountPoint.X+x, mountPoint.Y+screenLine, r, nil, tcell.StyleDefault.Reverse(true))
				continue
			}
			screen.SetContent(mountPoint.X+x, mountPoint.Y+screenLine, r, nil, tcell.StyleDefault)
		}
		screenLine++
	}
	if showCursor {
		screen.ShowCursor(mountPoint.X+s.cursorPos.X, mountPoint.Y+s.cursorPos.Y)
	}
}

func (s *Tab) SetSize(size layout.Size) {
	s.size = size
}

func (s *Tab) GetSize() layout.Size {
	return s.size
}

func (s *Tab) GetName() string {
	return fmt.Sprintf("Tab(%s)", s.name)
}
