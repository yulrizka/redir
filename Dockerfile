# build stage
FROM golang:1.14.1-buster as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o redir

# final stage
FROM alpine:3.11.5
RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app
COPY --from=builder /app/redir /app/
EXPOSE 5545
ENTRYPOINT ["./redir"]
