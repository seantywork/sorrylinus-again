


VENDOR_GOCV := vendor/gocv

VENDOR_YOLO := vendor/yolo

VENDOR_COTURN := vendor/coturn

VENDOR_FFMPEG := vendor/ffmpeg

all:

	@echo "go tv"


.PHONY: vendor
vendor: $(VENDOR_GOCV) $(VENDOR_YOLO) $(VENDOR_COTURN) $(VENDOR_FFMPEG)

	cd vendor/gocv && make install

$(VENDOR_GOCV):

	cd vendor && git clone https://github.com/hybridgroup/gocv.git


$(VENDOR_YOLO):

	mkdir -p vendor/yolo 

	wget https://pjreddie.com/media/files/yolov3.weights -O ./vendor/yolo/yolo.weights

	wget https://github.com/pjreddie/darknet/blob/master/cfg/yolov3.cfg?raw=true -O ./vendor/yolo/yolo.cfg

	wget https://github.com/pjreddie/darknet/blob/master/data/coco.names?raw=true -O ./vendor/yolo/coco.names

$(VENDOR_COTURN):

	sudo apt-get update -y 

	sudo apt-get install coturn -y

	mkdir -p vendor/coturn

$(VENDOR_FFMPEG):

	sudo apt-get update -y

	sudo apt-get install ffmpeg v4l-utils

	mkdir -p vendor/ffmpeg


build:

	go build -o soliagain.out .


.PHONY: test
test:

	go run test/test.go


clean:

	rm -r *.out