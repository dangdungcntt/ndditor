package layout

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type Element interface {
	GetName() string
	SetSize(size Size)
	GetSize() Size
	Render(screen tcell.Screen, mountPoint Point)
}

type Size struct {
	Width  int
	Height int
}

func (s Size) String() string {
	return fmt.Sprintf("Size(%d, %d)", s.Width, s.Height)
}

func (s Size) Subtract(size Size) Size {
	return Size{
		Width:  s.Width - size.Width,
		Height: s.Height - size.Height,
	}
}

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
}

func (p Point) AddSize(size Size) Point {
	return Point{
		X: p.X + size.Width,
		Y: p.Y + size.Height,
	}
}
