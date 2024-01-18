FROM golang:1.20-alpine

WORKDIR /app/testProject


COPY go.mod .
COPY go.sum .


RUN go mod download


COPY . .

RUN go build -o my-golang-app ./cmd/main.go


EXPOSE 8081


CMD ["./my-golang-app"]