package video

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

//MAP PHRASE:
//xxx.mp4
//1.Add key frame(added_xxx.mp4)
//2.Spilt(xxxn.ts)
//3.Distribute(xxxn.ts)
//4.Merge(merge_xxxn.ts)
//5.Reset timeline of .ts files(new_xxn.ts)

func convertToHLS(inputFile string, outputDirectory string, duration int, workDir string) error {
	//ffmpeg -i filename.mp4 -codec: copy -start_number 0 -hls_time 5 -hls_list_size 0 -f hls filename.m3u8
	// Create the output directory if it doesn't exist

	// if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
	// 	return err
	// }

	// Generate the output HLS playlist file,if use folder use this
	// outputPlaylist := filepath.Join(outputDirectory, "output.m3u8")

	// Run ffmpeg command
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-codec:", "copy",
		"-start_number", "0",
		"-hls_time", strconv.Itoa(duration),
		"-hls_list_size", "0",
		"-f", "hls",
		"-hls_flags", "split_by_time",
		// "-hls_segment_filename", filepath.Join(outputDirectory, "segment%03d.ts"),
		outputDirectory,
	)

	cmd.Dir = workDir
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

// mergeMP4 Merge n mp4 videos into one, not used
func mergeMP4(inputMP4s []string, outputDirectory string) error {
	//ffmpeg -i left.mp4 -i right.mp4 -filter_complex hstack output.mp4
	n := len(inputMP4s)

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

// StackChunks call hstak,vstack or grid filter in ffmpeg
func stackChunks(inputFiles []string, outputDirectory string, stackMethod string, workDir string) error {
	n := len(inputFiles)
	if n%2 != 0 {
		return errors.New("number of input is incorrect")
	}
	args := []string{}
	for i := 0; i < n; i++ {
		args = append(args, "-i")
		args = append(args, inputFiles[i])
	}
	args = append(args, "-filter_complex")
	switch stackMethod {
	case "hstak":
		{
			if n > 2 {
				args = append(args, "hstack=inputs="+strconv.Itoa(n))
			} else {
				args = append(args, "hstack")
			}
		}
	case "vstack":
		{
			if n > 2 {
				args = append(args, "vstack=inputs="+strconv.Itoa(n))
			} else {
				args = append(args, "vstack")
			}
		}
	case "grid": //up to 2*3,not test
		{
			if n != 4 || n != 6 {
				return errors.New("incorrect input for grid")
			} else {
				if n == 4 {
					args = append(args, `"[0:v][1:v]hstack=inputs=2[top]; [2:v][3:v]hstack=inputs=2[bottom]; [top][bottom]vstack=inputs=2[v]"`)

				}
				if n == 6 {
					args = append(args, `"[0:v][1:v][2:v]hstack=inputs=3[top]; [3:v][4:v][5:v]hstack=inputs=3[bottom]; [top][bottom]vstack=inputs=2[v]"`)
				}
				args = append(args, "-map")
				args = append(args, `"[v]"`)
			}
		}
	default:
		{
			return errors.New("unacceptable stack method")
		}
	}

	args = append(args, "-c:v")
	args = append(args, "libx264")
	args = append(args, "-c:a")
	args = append(args, "aac")
	args = append(args, outputDirectory)
	err := runFFmpegCommend(args, workDir)
	if err != nil {
		return err
	}
	return nil
}

func addKeyFrame(inputFile string, outputDirectory string, duration int, workDir string) error {
	keyint := fmt.Sprintf("keyint=%d:min-keyint=%d", duration, duration)
	args := []string{
		"-i", inputFile,
		"-codec:v:", "libx264",
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

func changeKeyFrameInterval(inputFile string, outputDirectory string, duration, FPS int, workDir string) error {
	keyint := fmt.Sprintf("keyint=%d:scenecut=0", duration*FPS)
	args := []string{
		"-i", inputFile,
		"-codec:v:", "libx264",
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

// getVideoStartTime extract the start_time of a certain video file
func getVideoStartTime(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=start_time", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("Error running ffprobe: %v", err)
	}

	startTimeStr := string(output)
	startTimeStr = regexp.MustCompile(`\r?\n`).ReplaceAllString(startTimeStr, "") // Remove newline characters

	startTime, err := strconv.ParseFloat(startTimeStr, 64)
	if err != nil {
		return 0, fmt.Errorf("Error parsing start_time: %v", err)
	}

	return startTime, nil
}

func getVideoFPS(filePath string) (int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("error running ffprobe: %v", err)
	}

	fpsStr := string(output)
	fpsStr = regexp.MustCompile(`\r?\n`).ReplaceAllString(fpsStr, "") // Remove newline characters

	// Parse the numerator and denominator of the frame rate
	numerator, denominator := 0, 0
	_, err = fmt.Sscanf(fpsStr, "%d/%d", &numerator, &denominator)
	if err != nil {
		return 0, fmt.Errorf("error parsing frame rate: %v", err)
	}

	// Calculate the fps as an integer (rounding to the nearest integer)
	fps := int(float64(numerator)/float64(denominator) + 0.5)

	return fps, nil
}

func resetTimeStamp(inputFile, outputDirectory string, index int, duration float64, startTime float64, workDir string) error {
	args := []string{
		"-i", inputFile,
		"-output_ts_offset", strconv.FormatFloat((duration * float64(index)), 'f', -1, 64),
		"-c", "copy",
		outputDirectory,
	}
	runFFmpegCommend(args, workDir)
	return nil
}

// Mp4toHLS convert mp4 to hls after add key_frame,input is name without suffix ,call in creator
func Mp4toHLS(inputFileName string, duration int, workDir string) error {
	//xxx.mp4->added_xxx.mp4->xxxn.ts
	err := addKeyFrame(inputFileName+".mp4", "added_"+inputFileName+".mp4", duration, workDir)
	if err != nil {
		return err
	}
	err = convertToHLS("added_"+inputFileName+".mp4", inputFileName+".m3u8", duration, workDir)
	if err != nil {
		return err
	}
	return nil
}

func Mp4toHLS_2(inputFileName string, duration, FPS int, workDir string) error {

	err := changeKeyFrameInterval(inputFileName+".mp4", "added_"+inputFileName+".mp4", duration, FPS, workDir)
	if err != nil {
		return err
	}
	err = convertToHLS("added_"+inputFileName+".mp4", inputFileName+".m3u8", duration, workDir)
	if err != nil {
		return err
	}
	return nil
}

// MergeTSFile merge ts file with same timestamp and reset start_time,input is names with suffix ,call in worker
func MergeTSFile(inputFileName []string, outputFileName string, index int, stackMethod string, duration int, workDir string) error {
	//xxxn.ts->merge_xxxn.ts->new_xxxn.ts

	err := stackChunks(inputFileName, "merged_"+outputFileName, stackMethod, workDir)
	if err != nil {
		return err
	}
	err = resetTimeStamp("merged_"+outputFileName, "new_"+outputFileName, index, float64(duration), 0, workDir)
	if err != nil {
		return err
	}
	return nil
}

func NewM3u8(m3u8FilePath, newM3u8FilePath string) error {
	inputFile, err := os.Open(m3u8FilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(newM3u8FilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, ".ts") {
			line = "new_" + line // Modify the line
		}
		_, err := outputFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
