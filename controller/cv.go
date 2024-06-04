package controller

import (
	"fmt"

	pkgcv "github.com/lineworld-lab/go-tv/pkg/cv"
)

type CV_CTL struct {
	OutMode CV_OUT_MODE
}

type CV_OUT_MODE struct {
	RAW_Window    bool
	YOLO_Window   bool
	YOLO_Std      bool
	YOLO_Endpoint bool
}

func (cvctl *CV_CTL) Start() error {

	if err := cvctl.Verify(); err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	if cvctl.OutMode.RAW_Window {

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

	} else if cvctl.OutMode.YOLO_Window {

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

	} else if cvctl.OutMode.YOLO_Std {

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

	} else if cvctl.OutMode.YOLO_Endpoint {

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

	}

	return nil
}

func (cvctl *CV_CTL) Verify() error {

	bool_count := 0

	if cvctl.OutMode.RAW_Window {
		bool_count += 1
	}

	if cvctl.OutMode.YOLO_Window {
		bool_count += 1
	}

	if cvctl.OutMode.YOLO_Std {

		bool_count += 1
	}

	if cvctl.OutMode.YOLO_Endpoint {

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
