FROM debian:12

ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /home

RUN apt-get update 

RUN apt-get install -y make build-essential ca-certificates

COPY --from=golang:1.21 /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

COPY . . 

RUN make clean

RUN go clean -modcache

RUN go mod tidy

RUN	make build

#CMD ["tail", "-f","/dev/null"]

CMD ["./soliagain.out"]