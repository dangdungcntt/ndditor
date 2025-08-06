package editor

import (
	"github.com/btvoidx/mint"
	"github.com/dangdungcntt/ndditor/editor/layout"
	"github.com/gdamore/tcell/v2"
)

// CursorEventListener is an interface for elements that can move the cursor inside itself
type CursorEventListener interface {
	layout.Element
	MoveCursor(dx, dy int)
	Focus()
	Blur()
}

var eventEmitter *mint.Emitter

// InitEventEmitter initializes the event emitter
func InitEventEmitter() {
	eventEmitter = new(mint.Emitter)
}

// OnEvent registers a function to be called when an event is emitted
func OnEvent[T any](fn func(e T)) (off func() <-chan struct{}) {
	if eventEmitter == nil {
		InitEventEmitter()
	}
	return mint.On(eventEmitter, fn)
}

// EmitEvent emits an event
func EmitEvent[T any](e T) {
	mint.Emit(eventEmitter, e)
}

// ModeChangedEvent emits when the mode changes
type ModeChangedEvent struct {
	Mode int
}

// StateChangedEvent emits when the state changes
type StateChangedEvent struct {
}

// KeyEvent emits when a key is pressed
// Not including the move cursor event (Left, Right, Up, Down)
type KeyEvent struct {
	Target layout.Element
	Ev     *tcell.EventKey
}

// SubmittedCommandEvent emits when a command is submitted
type SubmittedCommandEvent struct {
	Command string
}
