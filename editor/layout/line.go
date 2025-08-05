package layout

import (
	"github.com/gdamore/tcell/v2"
)

var _ Element = (*VLine)(nil)
var _ Element = (*HLine)(nil)

type VLine struct {
	Height int
}

func (v *VLine) GetName() string {
	return "VLine"
}

func (v *VLine) SetSize(size Size) {
	v.Height = size.Height
}

func (v *VLine) GetSize() Size {
	return Size{
		Width:  1,
		Height: v.Height,
	}
}

func (v *VLine) Render(screen tcell.Screen, mountPoint Point) {
	for y := 0; y < v.Height; y++ {
		screen.SetContent(mountPoint.X, mountPoint.Y+y, tcell.RuneVLine, nil, tcell.StyleDefault)
	}
}

type HLine struct {
	Width int
}

func (h *HLine) GetName() string {
	return "HLine"
}

func (h *HLine) SetSize(size Size) {
	h.Width = size.Width
}

func (h *HLine) GetSize() Size {
	return Size{
		Width:  h.Width,
		Height: 1,
	}
}

func (h *HLine) Render(screen tcell.Screen, mountPoint Point) {
	for x := 0; x < h.Width; x++ {
		screen.SetContent(mountPoint.X+x, mountPoint.Y, tcell.RuneHLine, nil, tcell.StyleDefault)
	}
}
