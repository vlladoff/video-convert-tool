package videoconvert

import (
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
)

func Convert(inputFile, outputFile string, width, height int) error {
	params := fmt.Sprintf("scale=w=%v:h=%v", width, height)

	// not working
	log.SetOutput(io.Discard)

	err := ffmpeg.Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{"vf": params}).
		OverWriteOutput().ErrorToStdOut().Run()

	return err
}
