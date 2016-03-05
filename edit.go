package main

import (
	"errors"

	"github.com/nsf/termbox-go"
)

var ErrDontSave = errors.New("don't save")

func edit(data []byte) ([]byte, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)

	ui := newUI(data)
	go ui.eventLoop()

	for ui.running {
		err = ui.draw()
		if err != nil {
			return nil, err
		}
		ui.update()
	}
	if ui.dontSave {
		return nil, ErrDontSave
	}

	return toBytes(ui.lines), nil
}
