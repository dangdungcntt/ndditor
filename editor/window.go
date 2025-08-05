package editor

import (
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
	"github.com/samber/lo"
)

var _ layout.Element = (*Window)(nil)

type Window struct {
	size      layout.Size
	tabs      []*Tab
	activeTab int
}

func NewWindow() *Window {
	return &Window{}
}

func (s *Window) GetActiveTab() *Tab {
	return s.tabs[s.activeTab]
}

func (s *Window) SetActiveTab(index int) {
	s.activeTab = index
}

func (s *Window) AddTab(tab *Tab) {
	s.tabs = append(s.tabs, tab)
	s.SetActiveTab(len(s.tabs) - 1)
}

func (s *Window) PreviousTab() {
	if s.activeTab > 0 {
		s.SetActiveTab(s.activeTab - 1)
	}
}

func (s *Window) NextTab() {
	if s.activeTab < len(s.tabs)-1 {
		s.SetActiveTab(s.activeTab + 1)
	}
}

func (s *Window) CloseTab() {
	if len(s.tabs) > 1 {
		s.tabs = append(s.tabs[:s.activeTab], s.tabs[s.activeTab+1:]...)
		if s.activeTab >= len(s.tabs) {
			s.activeTab = len(s.tabs) - 1
		}
	} else {
		s.tabs = []*Tab{}
	}
}

func (s *Window) SetSize(size layout.Size) {
	s.size = size
}

func (s *Window) Render(screen tcell.Screen, point layout.Point) {
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
	column.SetSize(s.size)
	column.Render(screen, point)
}

func (s *Window) GetSize() layout.Size {
	return s.size
}

func (s *Window) GetName() string {
	return "Window"
}

func (s *Window) getTitleComponent() layout.Element {
	var titles []layout.Element
	tabCount := len(s.tabs)
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
