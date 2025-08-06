package layout

import (
	"github.com/gdamore/tcell/v2"
)

// DrawText draws text from (x1, y1) to (x2, y2)
func DrawText(screen tcell.Screen, p1 Point, p2 Point, content string, color ...tcell.Color) {
	style := tcell.StyleDefault
	if len(color) > 0 {
		style = style.Foreground(color[0])
	}
	drawText(screen, p1.X, p1.Y, p2.X, p2.Y, style, content)
}

// DrawBox draws a box from (x1, y1) to (x2, y2)
func DrawBox(screen tcell.Screen, p1 Point, p2 Point) {
	drawBox(screen, p1.X, p1.Y, p2.X, p2.Y, tcell.StyleDefault)
}

// DrawVLine draws a vertical line from (x, y1) to (x, y2)
func DrawVLine(screen tcell.Screen, x, y1, y2 int) {
	for row := y1; row <= y2; row++ {
		screen.SetContent(x, row, tcell.RuneVLine, nil, tcell.StyleDefault)
	}
}

// DrawHLine draws a horizontal line from (x1, y) to (x2, y)
func DrawHLine(screen tcell.Screen, y, x1, x2 int, color ...tcell.Color) {
	style := tcell.StyleDefault
	if len(color) > 0 {
		style = style.Foreground(color[0])
	}
	for col := x1; col <= x2; col++ {
		screen.SetContent(col, y, tcell.RuneHLine, nil, style)
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
