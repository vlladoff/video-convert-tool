package video_convert

import (
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Convert(inputFile, outputFile string, width, height int) error {
	params := fmt.Sprintf("scale=w=%v:h=%v", width, height)
	err := ffmpeg.Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{"vf": params}).
		OverWriteOutput().ErrorToStdOut().Run()

	return err
}
