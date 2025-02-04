package ffprobe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

var binPath = "ffprobe"

// SetFFProbeBinPath sets the global path to find and execute the ffprobe program
func SetFFProbeBinPath(newBinPath string) {
	binPath = newBinPath
}

// ProbeURL is used to probe the given media file using ffprobe. The URL can be a local path, a HTTP URL or any other
// protocol supported by ffprobe, see here for a full list: https://ffmpeg.org/ffmpeg-protocols.html
// This function takes a context to allow killing the ffprobe process if it takes too long or in case of shutdown.
// Any additional ffprobe parameter can be supplied as well using extraFFProbeOptions.
func ProbeURL(fileURL string, extraFFProbeOptions ...string) (data *ProbeData, err error) {
	args := buildArgs(extraFFProbeOptions)

	// Add the file argument
	args = append(args, fileURL)

	cmd := exec.Command(binPath, args...)
	cmd.SysProcAttr = procAttributes()

	return runProbe(cmd)
}

// ProbeReader is used to probe a media file using an io.Reader. The reader is piped to the stdin of the ffprobe command
// and the data is returned.
// This function takes a context to allow killing the ffprobe process if it takes too long or in case of shutdown.
// Any additional ffprobe parameter can be supplied as well using extraFFProbeOptions.
func ProbeReader(reader io.Reader, extraFFProbeOptions ...string) (data *ProbeData, err error) {
	args := buildArgs(extraFFProbeOptions)

	// Add the file from stdin argument
	args = append(args, "-")

	cmd := exec.Command(binPath, args...)
	cmd.Stdin = reader
	cmd.SysProcAttr = procAttributes()

	return runProbe(cmd)
}

// runProbe takes the fully configured ffprobe command and executes it, returning the ffprobe data if everything went fine.
func runProbe(cmd *exec.Cmd) (*ProbeData, error) {
	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err := cmd.Run()

	probeErr := &rootError{}
	unmarshallingProbeError := json.Unmarshal(outputBuf.Bytes(), probeErr)
	if unmarshallingProbeError == nil && probeErr.Err.Message != "" {
		return nil, &probeErr.Err
	}

	if err != nil {
		return nil, &ProbeError{
			Message: fmt.Sprintf("error running %s [%s] %s", binPath, stdErr.String(), err.Error()),
		}
	}

	if stdErr.Len() > 0 {
		return nil, &ProbeError{
			Message: fmt.Sprintf("ffprobe error: %s", stdErr.String()),
		}
	}

	data := &ProbeData{}
	err = json.Unmarshal(outputBuf.Bytes(), data)
	if err != nil {
		return data, &ProbeError{
			Message: fmt.Sprintf("error parsing ffprobe output: %s", err.Error()),
		}
	}

	return data, nil
}

func buildArgs(extraFFProbeOptions []string) []string {
	args := append([]string{
		"-loglevel", "fatal",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_error",
	}, extraFFProbeOptions...)
	return args
}
