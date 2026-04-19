FROM golang:1.26.2-alpine AS builder
WORKDIR /app

RUN apk add --no-cache upx

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -trimpath -o /ghat ./cmd/ghat
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -trimpath -o /post ./cmd/post
RUN upx --best /ghat
RUN upx --best /post

FROM alpine:3.23.3
COPY --from=builder /ghat /ghat
COPY --from=builder /post /post
USER nobody
