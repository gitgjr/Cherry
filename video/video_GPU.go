package video

import (
	"fmt"
	"strconv"
)

func convertToHLS_GPU(inputFile string, outputDirectory string, duration int, workDir string) error {
	args := []string{}
	args = append(args, "-hwaccel", "cuda")

	args = append(args, "-i")
	args = append(args, inputFile)
	args = append(args, "-c:v")
	args = append(args, "h264_nvenc")

	args = append(args, "-codec:", "copy")
	args = append(args, "-start_number", "0")
	args = append(args, "-hls_time", strconv.Itoa(duration))
	args = append(args, "-hls_list_size", "0")
	args = append(args, "-f", "hls")
	args = append(args, "-hls_flags", "split_by_time", outputDirectory)

	err := runFFmpegCommend(args, workDir)
	if err != nil {
		return err
	}
	return nil
}

func addKeyFrame_GPU(inputFile string, outputDirectory string, duration int, workDir string) error {
	keyint := fmt.Sprintf("keyint=%d:min-keyint=%d", duration, duration)
	args := []string{
		"-hwaccel", "cuda",
		"-i", inputFile,
		"-codec:v:", "h264_cuvid",
		"-x264-params",
		keyint,
		"-codec:a", "copy",
		outputDirectory,
	}
	err := runFFmpegCommend(args, workDir)
	if err != nil {
		return err
	}
	return nil
}

// Mp4toHLS convert mp4 to hls after add key_frame,input is name without suffix ,call in creator
func Mp4toHLS_GPU(inputFileName string, duration int, workDir string) error {
	//xxx.mp4->added_xxx.mp4->xxxn.ts
	err := addKeyFrame_GPU(inputFileName+".mp4", "added_"+inputFileName+".mp4", duration, workDir)
	if err != nil {
		return err
	}
	err = convertToHLS_GPU("added_"+inputFileName+".mp4", inputFileName+".m3u8", duration, workDir)
	if err != nil {
		return err
	}
	return nil
}
