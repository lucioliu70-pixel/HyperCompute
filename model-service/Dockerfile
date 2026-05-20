FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/model-service ./model-service
FROM alpine:3.20
COPY --from=build /bin/model-service /model-service
CMD ["/model-service"]
