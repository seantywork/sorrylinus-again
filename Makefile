

all:

	@echo "go tv"


build:

	go build -o soliagain.out .


.PHONY: test
test:

	go run test/test.go

clean-data:

	rm -f data/media/*.json
	rm -f data/meda_video

clean:

	rm -r *.out