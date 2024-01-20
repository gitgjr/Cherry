package video

import "main/utils"

func runFFmpegCommend(commend string, args []string, workDir string) error {
	err := utils.RunCommend("ffmpeg", args, workDir)
	if err != nil {
		return err
	}
	return nil
}
