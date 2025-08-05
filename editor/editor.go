package editor

import (
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
	"strings"
)

type Editor struct {
	screen tcell.Screen
	state  *State
	root   layout.Element
	window *Window
}

func NewEditor(screen tcell.Screen) *Editor {
	window := NewWindow()
	tab := NewTab("new tab", NewEmptyLine(64))
	window.AddTab(tab)
	state := &State{
		mode: ModeView,
	}
	return &Editor{
		screen: screen,
		state:  state,
		window: window,
		root: &layout.Column{
			Children: []layout.Element{
				window,
				state,
			},
		},
	}
}

func (s *Editor) Run() {
	s.refreshScreen()

	for {
		ev := s.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.screen.Sync()
			s.refreshScreen()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlQ:
				s.window.PreviousTab()
			case tcell.KeyCtrlW:
				s.window.CloseTab()
			case tcell.KeyCtrlE:
				s.window.NextTab()
			case tcell.KeyCtrlT:
				s.window.AddTab(NewTab("new tab", NewEmptyLine(64)))
			case tcell.KeyEscape:
				if !s.state.IsMode(ModeView) {
					s.state.mode = ModeView
					break
				}
			case tcell.KeyEnter:
				if s.state.IsMode(ModeCommand) {
					s.executeCommand()
					break
				}
				if s.state.IsMode(ModeInsert) {
					s.getActiveTab().InsertNewline()
					break
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if s.state.IsMode(ModeCommand) {
					s.state.DeleteLastRuneFromCommand()
					break
				}
				if s.state.IsMode(ModeInsert) {
					s.getActiveTab().Backspace()
					break
				}
			case tcell.KeyDelete:
				if s.state.IsMode(ModeInsert) {
					s.getActiveTab().Delete()
					break
				}
			case tcell.KeyLeft:
				s.getActiveTab().MoveCursor(-1, 0)
			case tcell.KeyRight:
				s.getActiveTab().MoveCursor(1, 0)
			case tcell.KeyUp:
				s.getActiveTab().MoveCursor(0, -1)
			case tcell.KeyDown:
				s.getActiveTab().MoveCursor(0, 1)
			default:
				if ev.Rune() == 0 {
					break
				}
				if s.state.IsMode(ModeView) {
					switch ev.Rune() {
					case 'i':
						s.state.SetMode(ModeInsert)
					case ':':
						s.state.SetMode(ModeCommand)
					}
					break
				}
				if s.state.IsMode(ModeCommand) {
					s.state.AppendToCommand(ev.Rune())
					break
				}
				if s.state.IsMode(ModeInsert) {
					s.getActiveTab().InsertRune(ev.Rune())
					break
				}
			}
			s.refreshScreen()
		}
		if s.state.IsFinished() {
			return
		}
	}
}

func (s *Editor) refreshScreen() {
	// TODO: can I only redraw the changed lines?
	s.screen.Clear()
	s.screen.HideCursor()
	screenW, screenH := s.screen.Size()
	s.root.SetSize(layout.Size{
		Width:  screenW,
		Height: screenH,
	})
	s.root.Render(s.screen, layout.Point{
		X: 0,
		Y: 0,
	})

	s.screen.Show()
}

func (s *Editor) executeCommand() {
	cmd := string(s.state.GetCommand())
	s.state.ClearCommand()
	s.state.SetMode(ModeView)
	if strings.Contains(cmd, "q") {
		s.state.SetFinished()
		return
	}

	// TODO : implement commands write
}

func (s *Editor) getActiveTab() *Tab {
	return s.window.GetActiveTab()
}
