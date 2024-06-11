package main

import (
	"fmt"

	tvctl "github.com/lineworld-lab/go-tv/controller"
)

func main() {

	err := tvctl.LoadConfig()

	if err != nil {

		fmt.Println(err.Error())

		return
	}

	tvctlrunner := tvctl.TV_CTL{
		TVMode: tvctl.TV_MODE{
			RAW_Window:    false,
			YOLO_Window:   false,
			YOLO_Std:      false,
			YOLO_Endpoint: false,
			STREAM_File:   false,
			STREAM_Peer:   false,
			STREAM_CCTV:   false,
			STREAM_Room:   true,
		},
	}

	if err := tvctlrunner.Start(); err != nil {

		fmt.Println(err.Error())

	}

}
