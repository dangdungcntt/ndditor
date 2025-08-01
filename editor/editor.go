package editor

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strings"
)

const (
	ModeView = iota
	ModeInsert
	ModeCommand
)

type Editor struct {
	screen              tcell.Screen
	mode                int
	viewportH           int
	viewportW           int
	cursorX             int
	cursorY             int
	pendingCommand      []rune
	lines               []*GapBuffer
	lineIndex           int
	finished            bool
	defaultStyle        tcell.Style
	cursorPositionStyle tcell.Style
}

func NewEditor(screen tcell.Screen) *Editor {
	return &Editor{
		screen:              screen,
		mode:                ModeView,
		cursorX:             0,
		cursorY:             0,
		lines:               []*GapBuffer{NewGapBuffer(64)}, // start with one empty line
		lineIndex:           0,
		defaultStyle:        tcell.StyleDefault,
		cursorPositionStyle: tcell.StyleDefault.Reverse(true),
	}
}

func (e *Editor) Run() {
	e.refreshScreen()

	for {
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.screen.Sync()
			e.refreshScreen()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				return
			case tcell.KeyEscape:
				if e.mode != ModeView {
					e.mode = ModeView
					break
				}
			case tcell.KeyEnter:
				if e.mode == ModeCommand {
					e.executeCommand()
					break
				}
				if e.mode == ModeInsert {
					e.insertNewline()
					break
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if e.mode == ModeCommand {
					e.pendingCommand = e.pendingCommand[:len(e.pendingCommand)-1]
					break
				}
				if e.mode == ModeInsert {
					e.backspace()
					break
				}
			case tcell.KeyDelete:
				if e.mode == ModeInsert {
					e.delete()
					break
				}
			case tcell.KeyLeft:
				e.moveCursor(-1, 0)
			case tcell.KeyRight:
				e.moveCursor(1, 0)
			case tcell.KeyUp:
				e.moveCursor(0, -1)
			case tcell.KeyDown:
				e.moveCursor(0, 1)
			default:
				if ev.Rune() == 0 {
					break
				}
				if e.mode == ModeView {
					switch ev.Rune() {
					case 'i':
						e.mode = ModeInsert
					case ':':
						e.mode = ModeCommand
					}
					break
				}
				if e.mode == ModeCommand {
					e.pendingCommand = append(e.pendingCommand, ev.Rune())
					break
				}
				if e.mode == ModeInsert {
					e.insertRune(ev.Rune())
				}
			}
			e.refreshScreen()
		}
		if e.finished {
			return
		}
	}
}

func (e *Editor) insertRune(r rune) {
	e.lines[e.lineIndex].Insert(r)
	e.cursorX++
}

func (e *Editor) backspace() {
	if e.cursorX > 0 {
		e.lines[e.lineIndex].DeleteBeforeCursor()
		e.cursorX--
	} else if e.lineIndex > 0 {
		aboveLine := e.lines[e.lineIndex-1]
		e.cursorX = aboveLine.Len()
		aboveLine.Append(e.lines[e.lineIndex])
		aboveLine.moveCursor(e.cursorX)
		copy(e.lines[e.lineIndex:], e.lines[e.lineIndex+1:])
		e.lines = e.lines[:len(e.lines)-1]
		e.lineIndex--
		e.cursorY--
		if e.cursorY < 0 {
			e.cursorY = 0
		}
	}
}

func (e *Editor) delete() {
	if e.cursorX < e.lines[e.lineIndex].Len() {
		e.lines[e.lineIndex].DeleteAfterCursor()
	} else if e.lineIndex < len(e.lines)-1 {
		e.lines[e.lineIndex].Append(e.lines[e.lineIndex+1])
		copy(e.lines[e.lineIndex+1:], e.lines[e.lineIndex+2:])
		e.lines = e.lines[:len(e.lines)-1]
	}
}

func (e *Editor) insertNewline() {
	line := e.lines[e.lineIndex]
	e.lineIndex++
	e.cursorY++
	if e.cursorY >= e.viewportH {
		e.cursorY = e.viewportH - 1
	}
	newLineContent := line.CutAfterCursor()
	var newLine *GapBuffer
	if len(newLineContent) == 0 {
		newLine = NewGapBuffer(64)
	} else {
		newLine = NewGapBufferWithContent(newLineContent, true)
	}
	e.lines = append(e.lines, nil)
	if e.lineIndex == len(e.lines)-1 {
		e.lines[e.lineIndex] = newLine
	} else {
		copy(e.lines[e.lineIndex+1:], e.lines[e.lineIndex:])
		e.lines[e.lineIndex] = newLine
	}
	e.cursorX = 0
}

func (e *Editor) moveCursor(dx, dy int) {
	maxCursorY := min(e.viewportH-1, e.cursorY+(len(e.lines)-e.lineIndex-1))
	e.lineIndex += dy
	if e.lineIndex < 0 {
		e.lineIndex = 0
	} else if e.lineIndex >= len(e.lines) {
		e.lineIndex = len(e.lines) - 1
	}
	e.cursorY += dy

	if e.cursorY < 0 {
		e.cursorY = 0
	} else if e.cursorY > maxCursorY {
		e.cursorY = maxCursorY
	}

	e.cursorX += dx
	if e.cursorX < 0 {
		e.cursorX = 0
	}
	maxLen := e.lines[e.lineIndex].Len()
	if e.cursorX > maxLen {
		e.cursorX = maxLen
	}
	e.lines[e.lineIndex].moveCursor(e.cursorX)
}

func (e *Editor) refreshScreen() {
	// TODO: can I only redraw the changed lines?
	e.screen.Clear()
	screenW, screenH := e.screen.Size()
	e.viewportW = screenW
	e.viewportH = screenH - 1
	for x, r := range e.getInfoLine() {
		e.screen.SetContent(x, screenH-1, r, nil, tcell.StyleDefault)
	}

	minLine := e.lineIndex - e.cursorY
	maxLine := e.lineIndex + (e.viewportH - e.cursorY)

	showCursor := true
	screenLine := 0
	for y, line := range e.lines {
		if y < minLine || y >= maxLine {
			continue
		}
		for x, r := range line.Runes() {
			if x == e.cursorX && y == e.lineIndex {
				showCursor = false
				e.screen.SetContent(x, screenLine, r, nil, e.cursorPositionStyle)
				continue
			}
			e.screen.SetContent(x, screenLine, r, nil, e.defaultStyle)
		}
		screenLine++
	}
	if showCursor {
		e.screen.ShowCursor(e.cursorX, e.cursorY)
	} else {
		e.screen.HideCursor()
	}
	e.screen.Show()
}

func (e *Editor) getInfoLine() string {
	if e.mode == ModeCommand {
		return fmt.Sprintf(":%s", string(e.pendingCommand))
	}
	mode := "VIEW"
	if e.mode == ModeInsert {
		mode = "INSERT"
	}
	return fmt.Sprintf("-- %s --", mode)
}

func (e *Editor) executeCommand() {
	cmd := string(e.pendingCommand)
	e.pendingCommand = nil
	e.mode = ModeView
	if strings.Contains(cmd, "q") {
		e.finished = true
		return
	}

	// TODO : implement commands write
}
