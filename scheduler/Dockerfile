FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/scheduler ./scheduler
FROM alpine:3.20
COPY --from=build /bin/scheduler /scheduler
CMD ["/scheduler"]
