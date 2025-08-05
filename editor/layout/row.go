package layout

import (
	"github.com/gdamore/tcell/v2"
)

var _ Element = (*Row)(nil)

type Row struct {
	Children []Element
}

func (c *Row) GetName() string {
	return "Row"
}

func (c *Row) SetSize(size Size) {
	var sizeByIndexes []Size
	var missingWidthIndexes []int
	remainingWidth := size.Width
	for index, child := range c.Children {
		childSize := child.GetSize()
		if childSize.Width == 0 {
			missingWidthIndexes = append(missingWidthIndexes, index)
		} else {
			remainingWidth -= childSize.Width
		}
		if childSize.Height == 0 {
			childSize.Height = size.Height
		}
		sizeByIndexes = append(sizeByIndexes, childSize)
	}

	if remainingWidth > 0 && len(missingWidthIndexes) > 0 {
		flexedWidth := remainingWidth / len(missingWidthIndexes)
		for i, index := range missingWidthIndexes {
			if i == len(missingWidthIndexes)-1 {
				sizeByIndexes[index].Width = remainingWidth
			} else {
				sizeByIndexes[index].Width = flexedWidth
				remainingWidth -= flexedWidth
			}
		}
	}

	for index, child := range c.Children {
		child.SetSize(sizeByIndexes[index])
	}
}

func (c *Row) GetSize() Size {
	size := Size{}
	isUnknownWidth := false

	for _, child := range c.Children {
		childSize := child.GetSize()
		if childSize.Width == 0 {
			isUnknownWidth = true
		}
		size.Height = max(size.Height, childSize.Height)
		size.Width += childSize.Width
	}

	if isUnknownWidth {
		size.Width = 0
	}

	return size
}

func (c *Row) Render(screen tcell.Screen, mountPoint Point) {
	for _, child := range c.Children {
		child.Render(screen, mountPoint)
		childSize := child.GetSize()
		mountPoint.X += childSize.Width
	}
}
