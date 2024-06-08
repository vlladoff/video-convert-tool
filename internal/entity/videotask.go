package entity

import (
	"fmt"
	"github.com/vlladoff/video-convert-tool/internal/lib/videoconvert"
	"os"
)

type VideoSaver interface {
	SaveVideo(path, s3Path string) error
}

type ConvertVideoTask struct {
	ID         int    `json:"id"`
	Path       string `json:"path"`
	OutputPath string `json:"output_path"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Ext        string `json:"ext"`
	videoSaver VideoSaver
}

type ConvertVideoTaskDone struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}

func (cvt *ConvertVideoTask) Process(id int) (bool, int) {
	fmt.Printf("Workerk n: %d, task id: %d", id, cvt.ID)

	//ext
	if cvt.Path != "" && cvt.OutputPath != "" && cvt.Width > 0 && cvt.Height > 0 {

		tmpFile, err := os.CreateTemp("", "output-*.mp4")
		if err != nil {
			//return fmt.Errorf("failed to create temporary file: %w", err)
		}

		defer func(tmpFile *os.File) {
			err := tmpFile.Close()
			if err != nil {

			}
		}(tmpFile)

		tmpFilePath := tmpFile.Name()

		err = videoconvert.Convert(cvt.Path, tmpFilePath, cvt.Width, cvt.Height)
		if err != nil {
			// todo log
			return false, id
		}

		err = cvt.videoSaver.SaveVideo(tmpFilePath, cvt.OutputPath)
		if err != nil {
			return false, 0
		}
	}

	return true, id
}
