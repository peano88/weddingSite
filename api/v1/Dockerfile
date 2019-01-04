FROM golang:latest
WORKDIR /go/src/api
COPY . .
RUN go get github.com/gorilla/mux
RUN go get github.com/rs/cors
RUN go get gopkg.in/mgo.v2
RUN go get github.com/thedevsaddam/renderer
RUN go install ./...
RUN go build -o main
CMD ["./main"]
