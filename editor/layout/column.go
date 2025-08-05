package layout

import (
	"github.com/gdamore/tcell/v2"
)

var _ Element = (*Column)(nil)

type Column struct {
	Children []Element
}

func (c *Column) GetName() string {
	return "Column"
}

func (c *Column) SetSize(size Size) {
	var sizeByIndexes []Size
	var missingHeightIndexes []int
	remainingHeight := size.Height
	for index, child := range c.Children {
		childSize := child.GetSize()
		if childSize.Height == 0 {
			missingHeightIndexes = append(missingHeightIndexes, index)
		} else {
			remainingHeight -= childSize.Height
		}
		if childSize.Width == 0 {
			childSize.Width = size.Width
		}

		sizeByIndexes = append(sizeByIndexes, childSize)
	}

	if remainingHeight > 0 && len(missingHeightIndexes) > 0 {
		flexedHeight := remainingHeight / len(missingHeightIndexes)
		for i, index := range missingHeightIndexes {
			if i == len(missingHeightIndexes)-1 {
				sizeByIndexes[index].Height = remainingHeight
			} else {
				sizeByIndexes[index].Height = flexedHeight
				remainingHeight -= flexedHeight
			}
		}
	}

	for index, child := range c.Children {
		child.SetSize(sizeByIndexes[index])
	}
}

func (c *Column) GetSize() Size {
	size := Size{}
	isUnknownHeight := false

	for _, child := range c.Children {
		childSize := child.GetSize()
		if childSize.Height == 0 {
			isUnknownHeight = true
		}
		size.Width = max(size.Width, childSize.Width)
		size.Height += childSize.Height
	}
	if isUnknownHeight {
		size.Height = 0
	}

	return size
}

func (c *Column) Render(screen tcell.Screen, mountPoint Point) {
	for _, child := range c.Children {
		child.Render(screen, mountPoint)
		childSize := child.GetSize()
		mountPoint.Y += childSize.Height
	}
}
