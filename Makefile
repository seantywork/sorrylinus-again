

all:

	@echo "go tv"


build:

	go build -o soliagain.out .


.PHONY: test
test:

	go run test/test.go


clean:

	rm -r *.out