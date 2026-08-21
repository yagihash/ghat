FROM golang:1.27.0-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder
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
