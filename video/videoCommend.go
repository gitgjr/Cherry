package video

import "main/utils"

func runFFmpegCommend(commend string, args []string) error {
	err := utils.RunCommend("ffmpeg", args)
	if err != nil {
		return err
	}
	return nil
}
