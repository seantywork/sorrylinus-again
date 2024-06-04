


VENDOR_GOCV := vendor/gocv

VENDOR_YOLO := vendor/yolo

all:

	@echo "go tv"


.PHONY: vendor
vendor: $(VENDOR_GOCV) $(VENDOR_YOLO)

	cd vendor/gocv && make install

$(VENDOR_GOCV):

	cd vendor && git clone https://github.com/hybridgroup/gocv.git


$(VENDOR_YOLO):

	mkdir -p vendor/yolo 

	wget https://pjreddie.com/media/files/yolov3.weights -O ./vendor/yolo/yolo.weights

	wget https://github.com/pjreddie/darknet/blob/master/cfg/yolov3.cfg?raw=true -O ./vendor/yolo/yolo.cfg

	wget https://github.com/pjreddie/darknet/blob/master/data/coco.names?raw=true -O ./vendor/yolo/coco.names


build:

	go build -o tv.out .


.PHONY: test
test:

	go run test/test.go


clean:

	rm -r *.out