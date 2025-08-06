package layout

import "github.com/gdamore/tcell/v2"

// FullBorder is a full border
var FullBorder = Border{
	Top:    true,
	Bottom: true,
	Left:   true,
	Right:  true,
}

// Border represents a border
type Border struct {
	Top, Bottom, Left, Right                               bool
	TopRightTee, TopLeftTee, BottomRightTee, BottomLeftTee rune
}

// IsFull returns true if the border is full
func (b Border) IsFull() bool {
	return b.Top && b.Bottom && b.Left && b.Right && (b.TopRightTee+b.TopLeftTee+b.BottomRightTee+b.BottomLeftTee == 0)
}

// GetTopLeftCorner returns the top left corner
func (b Border) GetTopLeftCorner() rune {
	if (b.Top || b.Left) && b.TopLeftTee != 0 {
		return b.TopLeftTee
	}
	if b.Top && b.Left {
		return tcell.RuneULCorner
	}
	return 0
}

// GetTopRightCorner returns the top right corner
func (b Border) GetTopRightCorner() rune {
	if (b.Top || b.Right) && b.TopRightTee != 0 {
		return b.TopRightTee
	}
	if b.Top && b.Right {
		return tcell.RuneURCorner
	}
	return 0
}

// GetBottomLeftCorner returns the bottom left corner
func (b Border) GetBottomLeftCorner() rune {
	if (b.Bottom || b.Left) && b.BottomLeftTee != 0 {
		return b.BottomLeftTee
	}
	if b.Bottom && b.Left {
		return tcell.RuneLLCorner
	}
	return 0
}

// GetBottomRightCorner returns the bottom right corner
func (b Border) GetBottomRightCorner() rune {
	if (b.Bottom || b.Right) && b.BottomRightTee != 0 {
		return b.BottomRightTee
	}
	if b.Bottom && b.Right {
		return tcell.RuneLRCorner
	}
	return 0
}
