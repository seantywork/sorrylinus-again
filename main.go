package main

import (
	"fmt"

	tvctl "github.com/lineworld-lab/go-tv/controller"
)

func main() {

	cvctl := tvctl.CV_CTL{
		OutMode: tvctl.CV_OUT_MODE{
			RAW_Window:    false,
			YOLO_Window:   false,
			YOLO_Std:      false,
			YOLO_Endpoint: true,
		},
	}

	if err := cvctl.Start(); err != nil {

		fmt.Println(err.Error())

	}

}
