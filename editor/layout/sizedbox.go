package layout

import (
	"github.com/gdamore/tcell/v2"
)

var _ Element = (*SizedBox)(nil)

// SizedBox represents a sized box
type SizedBox struct {
	BaseElement
	Border  Border
	Content string
	Child   Element
	// Size original size
	Size Size
}

// GetName returns the name of the sized box
func (b *SizedBox) GetName() string {
	return "SizedBox"
}

// SetRenderSize sets the render size
func (b *SizedBox) SetRenderSize(size Size) {
	b.BaseElement.SetRenderSize(size)
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
		b.Child.SetRenderSize(size.Subtract(sub))
	}
}

// GetPreferredSize returns the preferred size of the sized box
func (b *SizedBox) GetPreferredSize() Size {
	return b.Size
}

// Render renders the sized box
func (b *SizedBox) Render(screen tcell.Screen, mountPoint Point) Size {
	renderSize := b.GetRenderSize()
	if b.Border.IsFull() {
		p2 := mountPoint.AddSize(Size{
			Width:  renderSize.Width - 1,
			Height: renderSize.Height - 1,
		})
		DrawBox(screen, mountPoint, p2)
		if b.Content != "" {
			DrawText(screen, mountPoint.AddSize(Size{Width: 1, Height: 1}), p2, b.Content)
		} else if b.Child != nil {
			b.Child.Render(screen, mountPoint.AddSize(Size{Width: 1, Height: 1}))
		}
		return renderSize
	}

	if b.Border.Top {
		DrawHLine(screen, mountPoint.Y, mountPoint.X, mountPoint.X+renderSize.Width-2)
	}
	if b.Border.Bottom {
		DrawHLine(screen, mountPoint.Y+renderSize.Height-1, mountPoint.X, mountPoint.X+renderSize.Width-2)
	}
	if b.Border.Left {
		DrawVLine(screen, mountPoint.X, mountPoint.Y, mountPoint.Y+renderSize.Height-1)
	}
	if b.Border.Right {
		DrawVLine(screen, mountPoint.X+renderSize.Width-1, mountPoint.Y, mountPoint.Y+renderSize.Height-1)
	}

	if r := b.Border.GetTopLeftCorner(); r != 0 {
		screen.SetContent(mountPoint.X, mountPoint.Y, r, nil, tcell.StyleDefault)
	}

	if r := b.Border.GetTopRightCorner(); r != 0 {
		screen.SetContent(mountPoint.X+renderSize.Width-1, mountPoint.Y, r, nil, tcell.StyleDefault)
	}

	if r := b.Border.GetBottomLeftCorner(); r != 0 {
		screen.SetContent(mountPoint.X, mountPoint.Y+renderSize.Height-1, r, nil, tcell.StyleDefault)
	}

	if r := b.Border.GetBottomRightCorner(); r != 0 {
		screen.SetContent(mountPoint.X+renderSize.Width-1, mountPoint.Y+renderSize.Height-1, r, nil, tcell.StyleDefault)
	}

	p1Delta := Size{Width: 0, Height: 0}
	if b.Border.Left {
		p1Delta.Width++
	}
	if b.Border.Top {
		p1Delta.Height++
	}

	if b.Content != "" {
		p2Delta := Size{Width: renderSize.Width - p1Delta.Width, Height: renderSize.Height - p1Delta.Height}
		DrawText(screen, mountPoint.AddSize(p1Delta), mountPoint.AddSize(p2Delta), b.Content)
	} else if b.Child != nil {
		b.Child.Render(screen, mountPoint.AddSize(p1Delta))
	}

	return renderSize
}
