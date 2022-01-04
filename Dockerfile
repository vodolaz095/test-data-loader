FROM golang:1.13.3

RUN mkdir /opt/test-data-loader
ADD input/ /opt/test-data-loader/input
ADD main.go /opt/test-data-loader/main.go

WORKDIR /opt/test-data-loader
RUN	go run main.go --source-dir=./input/ --ignore-duplicates --out-file=output.json
RUN	cat output.json

