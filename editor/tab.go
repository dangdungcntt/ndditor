package editor

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"path"
)

var _ layout.Element = (*Tab)(nil)

// Tab represents the content of a Tab. It is a doubly linked list of lines.
type Tab struct {
	layout.BaseElement
	name      string
	path      string
	cursorPos layout.Point
	lineIndex int
	lines     []*Line
}

// NewTab creates a new Tab
func NewTab(name string, lines ...*Line) *Tab {
	if len(lines) > 0 {
		lines[0].moveCursorTo(0)
	} else {
		lines = append(lines, NewEmptyLine(64))
	}
	return &Tab{
		name:  name,
		lines: lines,
	}
}

// NewTabFromPath creates a new Tab from a file
func NewTabFromPath(filePath string) (*Tab, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("error reading file: %s", err)
		}
		tab := NewTab(path.Base(filePath), NewEmptyLine(64))
		tab.SetPath(filePath)
		return tab, nil
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("%s is not a file", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		println(err)
	}
	defer func() {
		_ = file.Close()
	}()
	scanner := bufio.NewScanner(file)
	var lines []*Line
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, NewLine([]rune(line)))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	tab := NewTab("", lines...)
	tab.SetPath(filePath)
	return tab, nil
}

// SetPath sets the save path of the tab
func (s *Tab) SetPath(p string) {
	s.path = p
	s.name = path.Base(p)
}

// GetPath returns the save path of the tab
func (s *Tab) GetPath() string {
	return s.path
}

// InsertNewline inserts a newline at the current cursor position
func (s *Tab) InsertNewline() {
	renderSize := s.GetRenderSize()
	line := s.lines[s.lineIndex]
	s.lineIndex++
	s.cursorPos.Y++
	if s.cursorPos.Y >= renderSize.Height {
		s.cursorPos.Y = renderSize.Height - 1
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

// InsertRune inserts a rune at the current cursor position
func (s *Tab) InsertRune(r rune) {
	s.lines[s.lineIndex].Insert(r)
	s.cursorPos.X++
}

// Backspace deletes the character before the cursor
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

// Delete deletes the character after the cursor
func (s *Tab) Delete() {
	if s.cursorPos.X < s.lines[s.lineIndex].Len() {
		s.lines[s.lineIndex].DeleteAfterCursor()
	} else if s.lineIndex < len(s.lines)-1 {
		s.lines[s.lineIndex].Append(s.lines[s.lineIndex+1])
		copy(s.lines[s.lineIndex+1:], s.lines[s.lineIndex+2:])
		s.lines = s.lines[:len(s.lines)-1]
	}
}

// MoveCursor moves the cursor in the active tab
func (s *Tab) MoveCursor(dx, dy int) {
	maxCursorY := min(s.GetRenderSize().Height-1, s.cursorPos.Y+(len(s.lines)-s.lineIndex-1))
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

// GetName returns the name of the window
func (s *Tab) GetName() string {
	return fmt.Sprintf("Tab(%s)", s.name)
}

// GetPreferredSize returns the preferred size of the window
func (s *Tab) GetPreferredSize() layout.Size {
	return layout.Size{} // Auto size
}

// Render renders the tab to the screen
func (s *Tab) Render(screen tcell.Screen, mountPoint layout.Point) layout.Size {
	renderSize := s.GetRenderSize()
	minLine := s.lineIndex - s.cursorPos.Y
	maxLine := s.lineIndex + (renderSize.Height - s.cursorPos.Y)

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

	return renderSize
}

// Save saves the tab
func (s *Tab) Save() error {
	if s.path == "" {
		return errors.New("tab has no path")
	}
	tmpPath := s.path + ".tmp"
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = tmpFile.Close()
	}()
	for i, line := range s.lines {
		bytes := line.Bytes()
		if i < len(s.lines)-1 {
			bytes = append(bytes, '\n')
		}
		_, err = tmpFile.Write(bytes)
		if err != nil {
			return err
		}
	}
	err = tmpFile.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, s.path)
}
