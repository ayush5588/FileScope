FROM golang:1.20-alpine

WORKDIR $GOPATH/src/github.com/ayush5588/FileScope

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY token.env token.env

RUN go build -o filescope

EXPOSE 8080

CMD [ "./filescope" ]
