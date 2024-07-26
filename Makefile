

all:

	@echo "sorrylinus-again"


build:

	go build -o soliagain.out .


.PHONY: test
test:

	go run test/test.go

clean-data:

	rm -rf data/media/*.json
	rm -rf data/media/article/*.json
	rm -rf data/media/image/*.json data/media/image/*.jpg data/media/image/*.jpeg data/media/image/*.png
	rm -rf data/media/video/*.json data/media/video/*.mp4
	rm -rf data/session/*.json 
	rm -rf data/user/*.json
	rm -rf data/log/*.txt
	rm -rf data/log/*.json


clean:

	rm -rf *.out