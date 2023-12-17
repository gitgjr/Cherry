package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func ConvertToHLS(inputFile string, outputDirectory string, duration string) error {
	//ffmpeg -i filename.mp4 -codec: copy -start_number 0 -hls_time 5 -hls_list_size 0 -f hls filename.m3u8
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		return err
	}

	// Generate the output HLS playlist file
	outputPlaylist := filepath.Join(outputDirectory, "output.m3u8")

	// Run ffmpeg command
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-codec:", "copy",
		"-start_number", "0",
		"-hls_time", duration,
		"-hls_list_size", "0",
		"-hls_segment_filename", filepath.Join(outputDirectory, "segment%03d.ts"),
		outputPlaylist,
	)

	// Redirect command output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running ffmpeg: %v", err)
	}

	fmt.Println("Conversion to HLS completed successfully.")
	return nil
}

// Merge n mp4 videos into one, not tested
func MergeMP4s(inputMP4s []string, outputDirectory string) error {
	//ffmpeg -i left.mp4 -i right.mp4 -filter_complex hstack output.mp4
	n := len(inputMP4s)
	//generate commands
	args := []string{}
	for i := 0; i < n; i++ {
		args = append(args, "-i")
		args = append(args, inputMP4s[i])
	}
	args = append(args, "-filter_complex")

	if n > 2 {
		args = append(args, "hstack=inputs="+strconv.Itoa(n))
	} else {
		args = append(args, "hstack")
	}

	args = append(args, outputDirectory)
	fmt.Println(args)
	cmd := exec.Command("ffmpeg", args...)

	// Redirect command output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running ffmpeg: %v", err)
	}
	fmt.Println("Merge successfully.")
	return nil
}

func SplitMP4() {}

//func main() {
//	inputFile := "input.mp4"        // Replace with your input MP4 file
//	outputDirectory := "output_hls" // Replace with your desired output directory
//
//	err := convertToHLS(inputFile, outputDirectory)
//	if err != nil {
//		fmt.Println("Error:", err)
//	}
//}
