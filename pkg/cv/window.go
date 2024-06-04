package cv

import (
	"fmt"

	"gocv.io/x/gocv"
)

type RAW_WINDOW struct {
	Window *gocv.Window
}

func OpenWindow() *RAW_WINDOW {

	var raw_window RAW_WINDOW

	ww := gocv.NewWindow("window")

	raw_window.Window = ww

	return &raw_window

}

func (raw_window *RAW_WINDOW) Close() error {

	err := raw_window.Window.Close()

	if err != nil {

		return fmt.Errorf("failed to close raw window: %s", err.Error())

	}

	return nil
}

func (raw_window *RAW_WINDOW) Render(frame gocv.Mat) {

	raw_window.Window.IMShow(frame)
	raw_window.Window.WaitKey(1)

}
