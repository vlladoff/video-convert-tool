FROM golang:1.21-alpine AS video_convert_tool
WORKDIR /app/video_convert_tool
COPY . .
RUN go mod download
COPY cmd/video_convert_tool/ .
RUN go build -o video_convert_tool main.go
RUN apk add --no-cache ffmpeg

#FROM scratch AS video_convert_tool
#COPY --from=builder /app/video_convert_tool/.env .env
#COPY --from=builder /app/video_convert_tool/video_convert_tool /video_convert_tool

CMD ["./video_convert_tool"]