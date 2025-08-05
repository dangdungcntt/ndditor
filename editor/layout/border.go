package layout

import "github.com/gdamore/tcell/v2"

var FullBorder = Border{
	Top:    true,
	Bottom: true,
	Left:   true,
	Right:  true,
}

type Border struct {
	Top, Bottom, Left, Right                               bool
	TopRightTee, TopLeftTee, BottomRightTee, BottomLeftTee rune
}

func (b Border) IsFull() bool {
	return b.Top && b.Bottom && b.Left && b.Right && (b.TopRightTee+b.TopLeftTee+b.BottomRightTee+b.BottomLeftTee == 0)
}

func (b Border) GetTopLeftCorner() rune {
	if (b.Top || b.Left) && b.TopLeftTee != 0 {
		return b.TopLeftTee
	}
	if b.Top && b.Left {
		return tcell.RuneULCorner
	}
	return 0
}

func (b Border) GetTopRightCorner() rune {
	if (b.Top || b.Right) && b.TopRightTee != 0 {
		return b.TopRightTee
	}
	if b.Top && b.Right {
		return tcell.RuneURCorner
	}
	return 0
}

func (b Border) GetBottomLeftCorner() rune {
	if (b.Bottom || b.Left) && b.BottomLeftTee != 0 {
		return b.BottomLeftTee
	}
	if b.Bottom && b.Left {
		return tcell.RuneLLCorner
	}
	return 0
}

func (b Border) GetBottomRightCorner() rune {
	if (b.Bottom || b.Right) && b.BottomRightTee != 0 {
		return b.BottomRightTee
	}
	if b.Bottom && b.Right {
		return tcell.RuneLRCorner
	}
	return 0
}
