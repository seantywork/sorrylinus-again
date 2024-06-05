package controller

import (
	"fmt"

	pkgcv "github.com/lineworld-lab/go-tv/pkg/cv"
)

type TV_CTL struct {
	TVMode TV_MODE
}

type TV_MODE struct {
	RAW_Window    bool
	YOLO_Window   bool
	YOLO_Std      bool
	YOLO_Endpoint bool
	STREAM_File   bool
	STREAM_Peer   bool
}

func (tvctl *TV_CTL) Start() error {

	if err := tvctl.Verify(); err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	if tvctl.TVMode.RAW_Window {

		raww := pkgcv.OpenWindow()

		vi, err := pkgcv.GetVideoInput()

		if err != nil {

			return fmt.Errorf("start: %s", err.Error())

		}

		for {

			vi.VideoSource.Read(&vi.ImageMatrix)

			raww.Render(vi.ImageMatrix)

		}

		defer raww.Close()

		defer pkgcv.ReleaseVideoInput(vi)

	} else if tvctl.TVMode.YOLO_Window {

		raww := pkgcv.OpenWindow()

		vi, err := pkgcv.CreateYolo()

		if err != nil {

			return fmt.Errorf("start: %s", err.Error())

		}

		for {

			ods, err := vi.GetDetections()

			if err != nil {

				fmt.Println(err.Error())

				continue
			}

			vi.DrawDetections(ods)

			raww.Render(vi.VI.ImageMatrix)

		}

		defer raww.Close()

		defer pkgcv.Release(vi)

	} else if tvctl.TVMode.YOLO_Std {

		vi, err := pkgcv.CreateYolo()

		if err != nil {

			return fmt.Errorf("start: %s", err.Error())

		}

		for {

			ods, err := vi.GetDetections()

			if err != nil {

				fmt.Println(err.Error())

				continue
			}

			to_string := vi.FormatDetectionsToString(ods)

			fmt.Println(to_string)

		}

		defer pkgcv.Release(vi)

	} else if tvctl.TVMode.YOLO_Endpoint {

		vi, err := pkgcv.CreateYolo()

		if err != nil {

			return fmt.Errorf("start: %s", err.Error())

		}

		for {

			ods, err := vi.GetDetections()

			if err != nil {

				fmt.Println(err.Error())

				continue
			}

			to_json, err := vi.FormatDetectionsToJSON(ods)

			if err != nil {

				return fmt.Errorf("start: %s", err.Error())
			}

			fmt.Println(string(to_json))

		}

		defer pkgcv.Release(vi)

	} else if tvctl.TVMode.STREAM_File {

		streamctl := STREAM_CTL{}

		if err := streamctl.StartFiles(); err != nil {

			return fmt.Errorf("start: %s", err.Error())

		}

	} else if tvctl.TVMode.STREAM_Peer {

		streamctl := STREAM_CTL{}

		if err := streamctl.StartPeers(); err != nil {

			return fmt.Errorf("start: %s", err.Error())

		}

	}

	return nil
}

func (tvctl *TV_CTL) Verify() error {

	bool_count := 0

	if tvctl.TVMode.RAW_Window {
		bool_count += 1
	}

	if tvctl.TVMode.YOLO_Window {
		bool_count += 1
	}

	if tvctl.TVMode.YOLO_Std {

		bool_count += 1
	}

	if tvctl.TVMode.YOLO_Endpoint {

		bool_count += 1
	}

	if tvctl.TVMode.STREAM_File {

		bool_count += 1
	}

	if tvctl.TVMode.STREAM_Peer {

		bool_count += 1
	}

	if bool_count > 1 {

		return fmt.Errorf("more than one output mode selected")

	}

	if bool_count < 1 {

		return fmt.Errorf("no output mode selected")

	}

	return nil
}
