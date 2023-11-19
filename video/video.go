package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertToHLS(inputFile string, outputDirectory string) error {
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
		"-c:v", "libx264",
		"-hls_time", "5",
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

//func main() {
//	inputFile := "input.mp4"        // Replace with your input MP4 file
//	outputDirectory := "output_hls" // Replace with your desired output directory
//
//	err := convertToHLS(inputFile, outputDirectory)
//	if err != nil {
//		fmt.Println("Error:", err)
//	}
//}
