FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/wallet ./wallet
FROM alpine:3.20
COPY --from=build /bin/wallet /wallet
CMD ["/wallet"]
