// NDDitor
// Copyright (c) 2025, Dung Nguyen Dang <dangdungcntt@gmail.com>
// https://github.com/dangdungcntt/ndditor
package main

import (
	"github.com/dangdungcntt/ndditor/editor"
	"log"
	"os"

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

	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	editor.NewEditor(screen).Run(args)
}
