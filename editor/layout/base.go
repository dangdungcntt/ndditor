// Package layout provides a layout system for the editor
package layout

// BaseElement provides the base functionality for all elements
type BaseElement struct {
	renderSize Size
	isFocused  bool
}

// SetRenderSize sets the render size
func (s *BaseElement) SetRenderSize(size Size) {
	s.renderSize = size
}

// GetRenderSize returns the render size
func (s *BaseElement) GetRenderSize() Size {
	return s.renderSize
}

// Focus focuses the element
func (s *BaseElement) Focus() {
	s.isFocused = true
}

// Blur blurs the element
func (s *BaseElement) Blur() {
	s.isFocused = false
}

// IsFocused returns true if the element is focused
func (s *BaseElement) IsFocused() bool {
	return s.isFocused
}
