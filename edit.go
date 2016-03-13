package main

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/nsf/termbox-go"
)

var ErrDontSave = errors.New("don't save")

func edit(data []byte) ([]byte, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	termbox.SetInputMode(termbox.InputAlt)

	ui := newUI(data)
	go ui.eventLoop()

	defer func() {
		if r := recover(); r != nil {
			f, err := ioutil.TempFile("", "e-save")
			if err != nil {
				log.Println("editor crashed, but cannot create temp file to save your work, sorry...", err)
			} else {
				defer f.Close()
				_, err := f.Write(toBytes(ui.lines))
				if err != nil {
					log.Println("editor crashed, but cannot write to temp file to save your work, sorry...", err)
				} else {
					log.Println("editor crashed, but saved your work in ", f.Name(), ", sorry...")
				}
			}
			panic(r) // break execution in usual way
		}
	}()
	defer termbox.Close() // must be executed before recover above to enable logging

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
