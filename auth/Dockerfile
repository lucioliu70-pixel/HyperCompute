FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/auth ./auth
FROM alpine:3.20
COPY --from=build /bin/auth /auth
CMD ["/auth"]
