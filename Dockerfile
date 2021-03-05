FROM golang:1.15

WORKDIR /src
COPY . .

RUN go env -w GO111MODULE=on
RUN go mod download
RUN go build

EXPOSE 8080
ENTRYPOINT ["./open-balance"]
