package layout

import (
	"github.com/gdamore/tcell/v2"
)

var _ Element = (*Column)(nil)

// Column represents a column
type Column struct {
	Children []Element
}

// GetName returns the name of the element
func (c *Column) GetName() string {
	return "Column"
}

// SetRenderSize sets the render size, and calculates the render size of all children
func (c *Column) SetRenderSize(size Size) {
	sizeByIndexes := make([]Size, 0, len(c.Children))
	var missingHeightIndexes []int
	remainingHeight := size.Height
	for index, child := range c.Children {
		childSize := child.GetPreferredSize()
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
		child.SetRenderSize(sizeByIndexes[index])
	}
}

// GetPreferredSize returns the preferred size of the column
func (c *Column) GetPreferredSize() Size {
	size := Size{}
	isUnknownHeight := false

	for _, child := range c.Children {
		childSize := child.GetPreferredSize()
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

// Render renders the column
func (c *Column) Render(screen tcell.Screen, mountPoint Point) Size {
	renderSize := Size{}
	for _, child := range c.Children {
		childRenderSize := child.Render(screen, mountPoint)
		renderSize.Width = max(renderSize.Width, childRenderSize.Width)
		renderSize.Height += childRenderSize.Height
		mountPoint.Y += renderSize.Height
	}
	return renderSize
}
