package layout

import (
	"github.com/gdamore/tcell/v2"
)

var _ Element = (*SizedBox)(nil)

type SizedBox struct {
	Border  Border
	Content string
	Child   Element
	Size    Size
}

func (b *SizedBox) GetName() string {
	return "SizedBox"
}

func (b *SizedBox) SetSize(size Size) {
	b.Size = size
	if b.Child != nil {
		sub := Size{Width: 0, Height: 0}
		if b.Border.Left {
			sub.Width++
		}
		if b.Border.Top {
			sub.Height++
		}
		if b.Border.Right {
			sub.Width++
		}
		if b.Border.Bottom {
			sub.Height++
		}
		b.Child.SetSize(size.Subtract(sub))
	}
}

func (b *SizedBox) GetSize() Size {
	return b.Size
}

func (b *SizedBox) Render(screen tcell.Screen, mountPoint Point) {
	if b.Border.IsFull() {
		p2 := mountPoint.AddSize(Size{
			Width:  b.Size.Width - 1,
			Height: b.Size.Height - 1,
		})
		DrawBox(screen, mountPoint, p2)
		if b.Content != "" {
			DrawText(screen, mountPoint.AddSize(Size{Width: 1, Height: 1}), p2, b.Content)
		} else if b.Child != nil {
			b.Child.Render(screen, mountPoint.AddSize(Size{Width: 1, Height: 1}))
		}
		return
	}

	if b.Border.Top {
		DrawHLine(screen, mountPoint.Y, mountPoint.X, mountPoint.X+b.Size.Width-2)
	}
	if b.Border.Bottom {
		DrawHLine(screen, mountPoint.Y+b.Size.Height-1, mountPoint.X, mountPoint.X+b.Size.Width-2)
	}
	if b.Border.Left {
		DrawVLine(screen, mountPoint.X, mountPoint.Y, mountPoint.Y+b.Size.Height-1)
	}
	if b.Border.Right {
		DrawVLine(screen, mountPoint.X+b.Size.Width-1, mountPoint.Y, mountPoint.Y+b.Size.Height-1)
	}

	if r := b.Border.GetTopLeftCorner(); r != 0 {
		screen.SetContent(mountPoint.X, mountPoint.Y, r, nil, tcell.StyleDefault)
	}

	if r := b.Border.GetTopRightCorner(); r != 0 {
		screen.SetContent(mountPoint.X+b.Size.Width-1, mountPoint.Y, r, nil, tcell.StyleDefault)
	}

	if r := b.Border.GetBottomLeftCorner(); r != 0 {
		screen.SetContent(mountPoint.X, mountPoint.Y+b.Size.Height-1, r, nil, tcell.StyleDefault)
	}

	if r := b.Border.GetBottomRightCorner(); r != 0 {
		screen.SetContent(mountPoint.X+b.Size.Width-1, mountPoint.Y+b.Size.Height-1, r, nil, tcell.StyleDefault)
	}

	p1Delta := Size{Width: 0, Height: 0}
	if b.Border.Left {
		p1Delta.Width++
	}
	if b.Border.Top {
		p1Delta.Height++
	}

	if b.Content != "" {
		p2Delta := Size{Width: b.Size.Width - p1Delta.Width, Height: b.Size.Height - p1Delta.Height}
		DrawText(screen, mountPoint.AddSize(p1Delta), mountPoint.AddSize(p2Delta), b.Content)
	} else if b.Child != nil {
		b.Child.Render(screen, mountPoint.AddSize(p1Delta))
	}
}
