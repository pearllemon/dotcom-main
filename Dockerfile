FROM golang:1.15 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /server ./server

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /server ./
COPY public ./public
COPY views ./views
RUN chmod +x ./server
ENTRYPOINT ["./server"]
EXPOSE 8080 8443