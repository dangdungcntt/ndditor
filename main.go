package main

import (
	"github.com/dangdungcntt/ndditor/editor"
	"log"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("failed to create screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		log.Fatalf("failed to init screen: %v", err)
	}
	defer screen.Fini()

	screen.SetCursorStyle(tcell.CursorStyleSteadyBlock)

	editor.NewEditor(screen).Run()
}
