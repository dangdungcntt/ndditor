package layout

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

// Element is a interface for all layout elements
type Element interface {
	GetName() string
	GetPreferredSize() Size
	SetRenderSize(size Size)
	Render(screen tcell.Screen, mountPoint Point) Size
}

// Size is a 2D size
type Size struct {
	Width  int
	Height int
}

// String returns a string representation of the size
func (s Size) String() string {
	return fmt.Sprintf("Size(%d, %d)", s.Width, s.Height)
}

// Subtract subtracts the size from the size
func (s Size) Subtract(size Size) Size {
	return Size{
		Width:  s.Width - size.Width,
		Height: s.Height - size.Height,
	}
}

// Point is a 2D point
type Point struct {
	X int
	Y int
}

// String returns a string representation of the point
func (p Point) String() string {
	return fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
}

// AddSize adds the size to the point
func (p Point) AddSize(size Size) Point {
	return Point{
		X: p.X + size.Width,
		Y: p.Y + size.Height,
	}
}
