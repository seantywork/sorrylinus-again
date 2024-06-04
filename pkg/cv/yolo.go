package cv

import (
	"encoding/json"
	"fmt"

	"github.com/wimspaargaren/yolov3"
)

var (
	yolov3WeightsPath = "vendor/yolo/yolo.weights"
	yolov3ConfigPath  = "vendor/yolo/yolo.cfg"
	cocoNamesPath     = "vendor/yolo/coco.names"
)

type VI_YOLO struct {
	YoloNet yolov3.Net

	VI *VIDEO_IN
}

type YOLO_OD struct {
	Detected   string  `json:"detected"`
	Confidence float32 `json:"confidence"`
}

func CreateYolo() (*VI_YOLO, error) {

	var vi_yolo VI_YOLO

	vi, err := GetVideoInput()

	if err != nil {

		return &vi_yolo, fmt.Errorf("failed to create yolo: %s", err.Error())

	}

	vi_yolo.VI = vi

	yolonet, err := yolov3.NewNet(yolov3WeightsPath, yolov3ConfigPath, cocoNamesPath)

	if err != nil {

		return &vi_yolo, fmt.Errorf("failed to get new yolo net: %s", err.Error())

	}

	vi_yolo.YoloNet = yolonet

	return &vi_yolo, nil

}

func (vi_yolo *VI_YOLO) GetDetections() ([]yolov3.ObjectDetection, error) {

	if ok := vi_yolo.VI.VideoSource.Read(&vi_yolo.VI.ImageMatrix); !ok {

		return nil, fmt.Errorf("unable to read from stream")

	}

	if vi_yolo.VI.ImageMatrix.Empty() {

		return nil, fmt.Errorf("empty frame")

	}

	detections, err := vi_yolo.YoloNet.GetDetections(vi_yolo.VI.ImageMatrix)

	if err != nil {

		return nil, fmt.Errorf("unable to retrieve detections: %s", err.Error())
	}

	return detections, nil

}

func (vi_yolo *VI_YOLO) DrawDetections(detections []yolov3.ObjectDetection) {

	yolov3.DrawDetections(&vi_yolo.VI.ImageMatrix, detections)

}

func (vi_yolo *VI_YOLO) FormatDetectionsToString(detections []yolov3.ObjectDetection) string {

	ret_str := ""

	detections_len := len(detections)

	for i := 0; i < detections_len; i++ {

		ret_str += fmt.Sprintf("detected: %s, confidence: %f\n", detections[i].ClassName, detections[i].Confidence)

	}

	return ret_str

}

func (vi_yolo *VI_YOLO) FormatDetectionsToJSON(detections []yolov3.ObjectDetection) ([]byte, error) {

	detections_json := make([]YOLO_OD, 0)

	detections_len := len(detections)

	for i := 0; i < detections_len; i++ {

		detections_json = append(detections_json, YOLO_OD{
			Detected:   detections[i].ClassName,
			Confidence: detections[i].Confidence,
		})

	}

	ret_byte, err := json.Marshal(detections_json)

	if err != nil {

		return nil, fmt.Errorf("failed to format detections: %s", err.Error())

	}

	return ret_byte, nil

}

func (vi_yolo *VI_YOLO) Release() error {

	err := vi_yolo.YoloNet.Close()

	if err != nil {

		return fmt.Errorf("failed to close vi yolo: %s", err.Error())
	}

	err = ReleaseVideoInput(vi_yolo.VI)

	if err != nil {

		return fmt.Errorf("failed to release video input: %s", err.Error())

	}

	return nil

}
