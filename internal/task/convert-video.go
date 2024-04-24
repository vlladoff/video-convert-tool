package task

import (
	"fmt"
	videoConvert "video_convert_tool/internal/lib/video-convert"
	workerPool "video_convert_tool/internal/worker-pool"
)

type ConvertVideoTask struct {
	workerPool.TaskMeta
	ID         int    `json:"id"`
	Path       string `json:"path"`
	OutputPath string `json:"output_path"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Ext        string `json:"ext"`
}

func (cvt *ConvertVideoTask) Process(id int) {
	fmt.Printf("Workerk n: %d, task id: %d", id, cvt.Id)

	if cvt.Path != "" && cvt.OutputPath != "" && cvt.Width > 0 && cvt.Height > 0 {
		err := videoConvert.Convert(cvt.Path, cvt.OutputPath, cvt.Width, cvt.Height)
		if err != nil {
			return
		}
		//
	}
}
