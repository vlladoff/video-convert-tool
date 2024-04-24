FROM --platform=linux/amd64 golang:1.21-alpine AS builder

RUN apk add --no-progress --no-cache gcc musl-dev

WORKDIR /app/video_convert_tool

COPY . .

COPY cmd/video_convert_tool/ .

RUN go build -tags musl -ldflags '-extldflags "-static"' -o video_convert_tool main.go

FROM alpine:latest as video_convert_tool

COPY --from=builder /app/video_convert_tool/.env .env

COPY --from=builder /app/video_convert_tool/video_convert_tool /video_convert_tool

RUN apk add --no-cache ffmpeg

CMD ["./video_convert_tool"]