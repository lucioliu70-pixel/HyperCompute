FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/control-plane ./control-plane
FROM alpine:3.20
COPY --from=build /bin/control-plane /control-plane
CMD ["/control-plane"]
