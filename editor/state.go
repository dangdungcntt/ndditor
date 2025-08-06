package editor

import (
	"fmt"
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
	"time"
)

const (
	// ModeView is the view mode
	ModeView = iota
	// ModeInsert is the insert mode
	ModeInsert
	// ModeCommand is the command mode
	ModeCommand
)

var _ layout.Element = (*State)(nil)
var _ CursorEventListener = (*State)(nil)

// State represents the state of the editor
type State struct {
	layout.BaseElement
	mode           int
	errorMessage   string
	pendingCommand *Line
	cursorX        int
	finished       bool
}

// NewState creates a new state
func NewState() *State {
	s := &State{
		mode: ModeView,
	}
	s.initEventListeners()
	return s
}

func (s *State) initEventListeners() {
	OnEvent(func(e KeyEvent) {
		switch e.Ev.Key() {
		case tcell.KeyEscape:
			if !s.IsMode(ModeView) {
				s.SetMode(ModeView)
			}
		case tcell.KeyEnter:
			if s.IsMode(ModeCommand) {
				cmd := s.GetCommand()
				s.SetMode(ModeView)
				EmitEvent(SubmittedCommandEvent{Command: cmd})
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if s.IsMode(ModeCommand) {
				s.Delete()
				return
			}
		default:
			if e.Ev.Rune() == 0 {
				return
			}
			if s.IsMode(ModeView) {
				switch e.Ev.Rune() {
				case 'i':
					s.SetMode(ModeInsert)
				case ':':
					s.SetMode(ModeCommand)
				}
				return
			}
			if s.IsMode(ModeCommand) {
				s.AppendToCommand(e.Ev.Rune())
			}
		}
	})
}

// IsMode returns true if the mode is m
func (s *State) IsMode(m int) bool {
	return s.mode == m
}

// SetMode sets the mode
func (s *State) SetMode(m int) {
	s.errorMessage = ""
	s.mode = m
	s.pendingCommand = NewEmptyLine(64)
	s.cursorX = 0
	EmitEvent(ModeChangedEvent{Mode: m})
}

// ToastMessage displays a message for a short amount of time
func (s *State) ToastMessage(msg string) {
	s.errorMessage = msg
	go func() {
		time.Sleep(1500 * time.Millisecond)
		s.errorMessage = ""
		EmitEvent(StateChangedEvent{})
	}()
}

// AppendToCommand appends a rune to the pending command
func (s *State) AppendToCommand(r rune) {
	s.pendingCommand.Insert(r)
	s.cursorX++
}

// WriteToCommand writes a string to the pending command
func (s *State) WriteToCommand(str string) {
	for _, r := range str {
		s.pendingCommand.Insert(r)
		s.cursorX++
	}
}

// GetCommand returns the pending command as a string
func (s *State) GetCommand() string {
	return string(s.pendingCommand.Bytes())
}

// Delete deletes the character before the cursor
func (s *State) Delete() {
	s.pendingCommand.DeleteBeforeCursor()
	if s.cursorX > 0 {
		s.cursorX--
	}
}

// IsFinished returns true if the state is finished
func (s *State) IsFinished() bool {
	return s.finished
}

// SetFinished sets the finished flag
func (s *State) SetFinished() {
	s.finished = true
}

// GetName returns the name of the state
func (s *State) GetName() string {
	return "State"
}

// GetPreferredSize returns the preferred size of the state
func (s *State) GetPreferredSize() layout.Size {
	return layout.Size{
		Height: 1,
	}
}

// Render renders the state
func (s *State) Render(screen tcell.Screen, point layout.Point) layout.Size {
	renderSize := s.GetRenderSize()
	layout.DrawText(screen, point, point.AddSize(renderSize), s.getInfoLine())
	if s.IsMode(ModeCommand) {
		screen.ShowCursor(point.X+s.cursorX+1, point.Y)
	}
	return renderSize
}

// MoveCursor moves the cursor
func (s *State) MoveCursor(dx int, _ int) {
	if dx == 0 {
		return
	}
	s.cursorX += dx
	if s.cursorX < 0 {
		s.cursorX = 0
	} else if s.cursorX >= s.pendingCommand.Len() {
		s.cursorX = s.pendingCommand.Len()
	}
	s.pendingCommand.moveCursorTo(s.cursorX)
}

func (s *State) getInfoLine() string {
	if s.errorMessage != "" {
		return s.errorMessage
	}
	if s.IsMode(ModeCommand) {
		return fmt.Sprintf(":%s", s.GetCommand())
	}
	mode := "VIEW"
	if s.IsMode(ModeInsert) {
		mode = "INSERT"
	}
	return fmt.Sprintf("-- %s --", mode)
}
