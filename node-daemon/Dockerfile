FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /node-daemon ./node-daemon
FROM alpine:3.20
COPY --from=build /node-daemon /node-daemon
ENTRYPOINT ["/node-daemon"]
