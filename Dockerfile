FROM golang:1.19

WORKDIR /service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY errors.go .
COPY main.go .

ADD datahandling datahandling

RUN go build -o categoryservice

EXPOSE 8080

ENTRYPOINT ["./categoryservice"]