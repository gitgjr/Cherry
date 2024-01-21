package video

import "main/utils"

func runFFmpegCommend(args []string, workDir string) error {
	err := utils.RunCommend("ffmpeg", args, workDir)
	if err != nil {
		return err
	}
	return nil
}
