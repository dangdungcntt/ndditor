// Package editor provides a text editor
package editor

import (
	"fmt"
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/dangdungcntt/ndditor/editor/logger"
	"github.com/gdamore/tcell/v2"
	"log"
	"regexp"
	"strings"
)

// GlobalState is the global state of the editor
var GlobalState *State

// Editor is the main editor
type Editor struct {
	events         chan tcell.Event
	screen         tcell.Screen
	root           layout.Element
	window         *Window
	focusedElement CursorEventListener
}

// NewEditor creates a new editor
func NewEditor(screen tcell.Screen) *Editor {
	GlobalState = NewState()
	return &Editor{
		screen: screen,
		events: make(chan tcell.Event),
	}
}

// Run starts the editor
func (s *Editor) Run(args []string) {
	s.initEventListeners()
	go s.eventLoop()

	s.window = s.initWindow(args)
	s.focusedElement = s.window
	s.root = &layout.Column{
		Children: []layout.Element{
			s.window,
			GlobalState,
		},
	}

	s.eventConsumer()
}

func (s *Editor) initWindow(args []string) *Window {
	window := NewWindow()

	if len(args) == 0 {
		tab := NewTab("new tab", NewEmptyLine(64))
		window.AddTab(tab)
		return window
	}
	filePath := args[0]

	tab, err := NewTabFromPath(filePath)
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	tab.SetPath(filePath)
	window.AddTab(tab)

	return window
}

func (s *Editor) eventLoop() {
	for {
		ev := s.screen.PollEvent()
		if GlobalState.IsFinished() {
			return
		}

		s.events <- ev
	}
}

func (s *Editor) eventConsumer() {
	s.render()

	for tEvent := range s.events {
		switch ev := tEvent.(type) {
		case *tcell.EventResize:
			s.screen.Sync()
			s.render()
		case *tcell.EventKey:
			logger.WriteLog(ev.Modifiers(), ev.Name(), ev.Key(), ev.Rune())
			switch ev.Key() {
			case tcell.KeyCtrlC:
				return
			case tcell.KeyLeft:
				s.moveCursor(-1, 0)
			case tcell.KeyRight:
				s.moveCursor(1, 0)
			case tcell.KeyUp:
				s.moveCursor(0, -1)
			case tcell.KeyDown:
				s.moveCursor(0, 1)
			default:
				if GlobalState.IsMode(ModeView) {
					EmitEvent(KeyEvent{
						Ev: ev,
					})
				} else {
					EmitEvent(KeyEvent{
						Target: s.focusedElement,
						Ev:     ev,
					})
				}
			}
			s.render()
		}
		if GlobalState.IsFinished() {
			close(s.events)
			return
		}
	}
}

func (s *Editor) render() {
	// TODO: can I only redraw the changed lines?
	s.screen.Clear()
	s.screen.HideCursor()
	screenW, screenH := s.screen.Size()
	s.root.SetRenderSize(layout.Size{
		Width:  screenW,
		Height: screenH,
	})
	s.root.Render(s.screen, layout.Point{
		X: 0,
		Y: 0,
	})

	s.screen.Show()
}

func (s *Editor) moveCursor(dx, dy int) {
	s.focusedElement.MoveCursor(dx, dy)
}

var cmdRegex = regexp.MustCompile("[wq]")

func (s *Editor) executeCommand(cmd string) {
	switch {
	case strings.HasPrefix(cmd, "path"):
		s.getActiveTab().SetPath(cmd[5:])
	case strings.HasPrefix(cmd, "open"):
		filePath := cmd[5:]
		tab, err := NewTabFromPath(filePath)
		if err != nil {
			GlobalState.ToastMessage(fmt.Sprintf("err: %v", err))
			return
		}
		s.window.AddTab(tab)
	case cmdRegex.MatchString(cmd):
		if strings.Contains(cmd, "w") {
			err := s.getActiveTab().Save()
			if err != nil {
				GlobalState.ToastMessage(fmt.Sprintf("err: %v", err))
				return
			}
		}

		if strings.Contains(cmd, "q") {
			GlobalState.SetFinished()
		}
	default:
		GlobalState.ToastMessage(fmt.Sprintf("unknown command: %s", cmd))
	}
}

func (s *Editor) getActiveTab() *Tab {
	return s.window.GetActiveTab()
}

func (s *Editor) initEventListeners() {
	OnEvent(func(e ModeChangedEvent) {
		if s.focusedElement != nil {
			s.focusedElement.Blur()
		}
		if e.Mode == ModeCommand {
			s.focusedElement = GlobalState
		} else {
			s.focusedElement = s.window
		}
		s.focusedElement.Focus()
	})
	OnEvent(func(_ StateChangedEvent) {
		s.render()
	})
	OnEvent(func(e SubmittedCommandEvent) {
		s.executeCommand(e.Command)
	})
}
