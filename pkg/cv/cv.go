package cv

import (
	"fmt"

	"gocv.io/x/gocv"
)

type VideoInterface interface {
	Release() error
}

type VIDEO_IN struct {
	VideoSource *gocv.VideoCapture

	ImageMatrix gocv.Mat
}

func GetVideoInput() (*VIDEO_IN, error) {

	var vi VIDEO_IN

	v_source, err := gocv.OpenVideoCapture(0)

	vi.VideoSource = v_source

	if err != nil {

		return &vi, fmt.Errorf("failed to open video device: %s", err.Error())

	}

	vi.VideoSource = v_source

	im := gocv.NewMat()

	vi.ImageMatrix = im

	return &vi, nil
}

func ReleaseVideoInput(vi *VIDEO_IN) error {

	vi.ImageMatrix.Close()

	vi.VideoSource.Close()

	return nil
}

func Release(vinf VideoInterface) error {

	err := vinf.Release()

	if err != nil {

		return fmt.Errorf("video interface: %s", err.Error())

	}

	return nil

}
