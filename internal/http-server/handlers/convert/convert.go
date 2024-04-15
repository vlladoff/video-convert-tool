package convert

import (
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"video_convert_tool/internal/lib/helper"
	videoconvert "video_convert_tool/internal/lib/video-convert"
)

type Request struct {
	File   *multipart.File `form:"file"`
	Width  int             `form:"width"`
	Height int             `form:"height"`
	Ext    string          `form:"output_ext" validate:"required"`
}

func ConvertVideo(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//check form
		err := r.ParseMultipartForm(1000 << 20)
		if err != nil {
			log.Error("failed to parse multipart form", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusBadRequest)

			return
		}

		//get request
		var req Request
		reader := strings.NewReader(r.Form.Encode())
		if err := render.DecodeForm(reader, &req); err != nil {
			log.Error("failed to decode form data", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusBadRequest)

			return
		}

		//get file
		file, fileHeader, err := r.FormFile(helper.GetStructTagByField(req, "File", "form"))
		req.File = &file
		if err != nil {
			log.Error("failed to retrieve file from form data", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusInternalServerError)

			return
		}
		defer file.Close()

		//create temp file
		tempFile, err := os.CreateTemp("", "vctfile-")
		if err != nil {
			log.Error("failed to create temporary file", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusInternalServerError)

			return
		}
		defer tempFile.Close()

		//copy data from form file to temp file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			log.Error("failed to write to temporary file", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusInternalServerError)

			return
		}

		//generate output file name
		outputFileName := "converted_" + strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename)) + "." + req.Ext

		//convert file
		err = videoconvert.Convert(tempFile.Name(), outputFileName, req.Width, req.Height)
		if err != nil {
			log.Error("failed to convert the file", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusInternalServerError)

			return
		}

		//remove file after send
		defer func() {
			os.Remove(tempFile.Name())
			os.Remove(outputFileName)
		}()

		//read file
		outputFile, err := os.ReadFile(outputFileName)
		if err != nil {
			log.Error("failed to read converted file", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusInternalServerError)

			return
		}

		//set headers
		w.Header().Set("Content-Disposition", "attachment; filename="+outputFileName)
		w.Header().Set("Content-Type", "text/plain")

		//send file
		_, err = w.Write(outputFile)
		if err != nil {
			log.Error("failed to write a file to the response", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			render.Status(r, http.StatusInternalServerError)

			return
		}
	}
}
