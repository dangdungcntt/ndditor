package layout

import (
	"github.com/gdamore/tcell/v2"
)

var _ Element = (*Row)(nil)

// Row represents a row
type Row struct {
	Children []Element
}

// GetName returns the name of the row
func (c *Row) GetName() string {
	return "Row"
}

// SetRenderSize sets the render size
func (c *Row) SetRenderSize(size Size) {
	sizeByIndexes := make([]Size, 0, len(c.Children))
	var missingWidthIndexes []int
	remainingWidth := size.Width
	for index, child := range c.Children {
		childSize := child.GetPreferredSize()
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
		child.SetRenderSize(sizeByIndexes[index])
	}
}

// GetPreferredSize returns the preferred size of the row which is calculated by
// getting the max preferred height of all children and the sum of preferred width
func (c *Row) GetPreferredSize() Size {
	size := Size{}
	isUnknownWidth := false

	for _, child := range c.Children {
		childSize := child.GetPreferredSize()
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

// Render renders the row
func (c *Row) Render(screen tcell.Screen, mountPoint Point) Size {
	renderSize := Size{}
	for _, child := range c.Children {
		childRenderSize := child.Render(screen, mountPoint)
		renderSize.Height = max(renderSize.Height, childRenderSize.Height)
		renderSize.Width += childRenderSize.Width
		mountPoint.X += childRenderSize.Width
	}
	return renderSize
}
