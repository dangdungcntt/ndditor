package editor

import (
	"fmt"
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
	"github.com/samber/lo"
)

var _ layout.Element = (*Window)(nil)

// Window represents a list of tabs
type Window struct {
	layout.BaseElement
	tabs      []*Tab
	activeTab int
}

// NewWindow creates a new window with an empty tab and registers event listeners
func NewWindow() *Window {
	w := &Window{}
	w.initEventListeners()
	return w
}

// GetActiveTab returns the active tab
func (s *Window) GetActiveTab() *Tab {
	return s.tabs[s.activeTab]
}

// SetActiveTab sets the active tab
func (s *Window) SetActiveTab(index int) {
	s.activeTab = index
}

// AddTab adds a new tab
func (s *Window) AddTab(tab *Tab) {
	s.tabs = append(s.tabs, tab)
	s.SetActiveTab(len(s.tabs) - 1)
}

// PreviousTab moves to the previous tab
func (s *Window) PreviousTab() {
	if s.activeTab > 0 {
		s.SetActiveTab(s.activeTab - 1)
	}
}

// NextTab moves to the next tab
func (s *Window) NextTab() {
	if s.activeTab < len(s.tabs)-1 {
		s.SetActiveTab(s.activeTab + 1)
	}
}

// CloseTab closes the active tab
func (s *Window) CloseTab() {
	if len(s.tabs) > 1 {
		s.tabs = append(s.tabs[:s.activeTab], s.tabs[s.activeTab+1:]...)
		if s.activeTab >= len(s.tabs) {
			s.activeTab = len(s.tabs) - 1
		}
	} else {
		s.tabs = []*Tab{}
		s.AddTab(NewTab("new tab", NewEmptyLine(64)))
	}
}

// GetPreferredSize returns the preferred size of the window
func (s *Window) GetPreferredSize() layout.Size {
	return layout.Size{} // Auto size
}

// GetName returns the name of the window
func (s *Window) GetName() string {
	return "Window"
}

// Render renders the window
func (s *Window) Render(screen tcell.Screen, point layout.Point) layout.Size {
	// render list tab names
	column := layout.Column{
		Children: []layout.Element{
			s.getTitleComponent(),
			&layout.SizedBox{
				Border: layout.Border{
					Left:   true,
					Right:  true,
					Bottom: true,
				},
				Child: s.tabs[s.activeTab],
			},
		},
	}
	column.SetRenderSize(s.GetRenderSize())
	return column.Render(screen, point)
}

func (s *Window) getTitleComponent() layout.Element {
	tabCount := len(s.tabs)
	titles := make([]layout.Element, 0, tabCount+1)
	for i, tab := range s.tabs {
		isFirst := i == 0
		var content string
		isLast := i == tabCount-1
		if i == s.activeTab {
			content = " > " + tab.name + " "
		} else {
			content = "   " + tab.name + " "
		}
		titles = append(titles, &layout.SizedBox{
			Border: layout.Border{
				Top:            true,
				Right:          true,
				Bottom:         true,
				Left:           isFirst,
				TopRightTee:    lo.Ternary(isLast, tcell.RuneURCorner, tcell.RuneTTee),
				BottomRightTee: tcell.RuneBTee,
				BottomLeftTee:  lo.Ternary(isFirst, tcell.RuneLTee, 0),
			},
			Size:    layout.Size{Width: len(content) + lo.Ternary(isFirst, 2, 1), Height: 3},
			Content: content,
		})
	}
	// add last border
	titles = append(titles, &layout.SizedBox{
		Border: layout.Border{
			Bottom:         true,
			BottomRightTee: tcell.RuneURCorner,
		},
	})

	return &layout.Row{
		Children: titles,
	}
}

// MoveCursor moves the cursor in the active tab
func (s *Window) MoveCursor(dx, dy int) {
	s.GetActiveTab().MoveCursor(dx, dy)
}

func (s *Window) initEventListeners() {
	OnEvent(func(e KeyEvent) {
		if e.Target != s {
			return
		}

		activeTab := s.GetActiveTab()
		switch e.Ev.Key() {
		case tcell.KeyCtrlQ:
			s.PreviousTab()
		case tcell.KeyCtrlW:
			s.CloseTab()
		case tcell.KeyCtrlE:
			s.NextTab()
		case tcell.KeyCtrlT:
			s.AddTab(NewTab("new tab", NewEmptyLine(64)))
		case tcell.KeyCtrlS:
			filePath := activeTab.GetPath()
			if filePath == "" {
				GlobalState.SetMode(ModeCommand)
				GlobalState.WriteToCommand("path ")
				break
			}
			err := activeTab.Save()
			if err != nil {
				GlobalState.ToastMessage(fmt.Sprintf("err: %v", err))
			}
		case tcell.KeyEnter:
			if GlobalState.IsMode(ModeInsert) {
				activeTab.InsertNewline()
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if GlobalState.IsMode(ModeInsert) {
				activeTab.Backspace()
			}
		case tcell.KeyDelete:
			if GlobalState.IsMode(ModeInsert) {
				activeTab.Delete()
			}
		default:
			if e.Ev.Rune() == 0 {
				return
			}
			if GlobalState.IsMode(ModeInsert) && s.IsFocused() {
				activeTab.InsertRune(e.Ev.Rune())
				return
			}
		}
	})
}
