package editor

import (
	"fmt"
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
)

const (
	ModeView = iota
	ModeInsert
	ModeCommand
)

var _ layout.Element = (*State)(nil)

type State struct {
	windowSize layout.Size
	// Current mode of the editor
	mode           int
	pendingCommand []rune
	finished       bool
}

func (s *State) IsMode(m int) bool {
	return s.mode == m
}

func (s *State) SetMode(m int) {
	s.mode = m
}

func (s *State) AppendToCommand(r rune) {
	s.pendingCommand = append(s.pendingCommand, r)
}

func (s *State) GetCommand() []rune {
	return s.pendingCommand
}

func (s *State) DeleteLastRuneFromCommand() {
	s.pendingCommand = s.pendingCommand[:len(s.pendingCommand)-1]
}

func (s *State) ClearCommand() {
	s.pendingCommand = nil
}

func (s *State) IsFinished() bool {
	return s.finished
}

func (s *State) SetFinished() {
	s.finished = true
}

func (s *State) GetInfoLine() string {
	if s.IsMode(ModeCommand) {
		return fmt.Sprintf(":%s", string(s.pendingCommand))
	}
	mode := "VIEW"
	if s.IsMode(ModeInsert) {
		mode = "INSERT"
	}
	return fmt.Sprintf("-- %s --", mode)
}

func (s *State) GetName() string {
	return "State"
}

func (s *State) SetSize(size layout.Size) {
	s.windowSize = size
}

func (s *State) GetSize() layout.Size {
	return layout.Size{
		Width:  s.windowSize.Width,
		Height: 1,
	}
}

func (s *State) Render(screen tcell.Screen, point layout.Point) {
	layout.DrawText(screen, point, point.AddSize(s.GetSize()), s.GetInfoLine())
}
