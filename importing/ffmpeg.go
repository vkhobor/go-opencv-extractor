package importing

import (
	"log/slog"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

var fpsRegexp = regexp.MustCompile(`\d{2}[.]?\d{0,2} fps`)

func extractMetadata(videoPath string) (fps float64, err error) {
	cmd := exec.Command("ffmpeg", "-i", videoPath)
	output, err := cmd.CombinedOutput()
	slog.Debug("ffmpeg output", "stdout", string(output))

	if err == nil {
		return 0, errors.Errorf("ffmpeg expected error but success: %v", string(output))
	}

	fpsString := fpsRegexp.FindString(string(output))
	fpsNumString := strings.Split(fpsString, " ")[0]
	fps, err = strconv.ParseFloat(fpsNumString, 64)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return fps, nil
}
