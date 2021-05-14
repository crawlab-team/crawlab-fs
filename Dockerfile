FROM golang:1.15

ADD ./go.mod .
ADD ./go.sum .
RUN go mod download

CMD ["sh", "./bin/test.sh"]
