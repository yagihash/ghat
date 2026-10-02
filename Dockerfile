FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS builder
WORKDIR /app

RUN apk add --no-cache upx

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -trimpath -o /ghat ./cmd/ghat
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -trimpath -o /post ./cmd/post
RUN upx --best /ghat
RUN upx --best /post

FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11
COPY --from=builder /ghat /ghat
COPY --from=builder /post /post
