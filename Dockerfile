FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .

RUN apk update
RUN apk add upx

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /ghat ./cmd/ghat/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /post ./cmd/post/main.go
RUN upx --best /ghat
RUN upx --best /post

FROM alpine:3.23.2
COPY --from=builder /ghat /ghat
COPY --from=builder /post /post