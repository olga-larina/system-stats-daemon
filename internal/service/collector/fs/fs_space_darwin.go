//go:build darwin
// +build darwin

package fs

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

// Пример строки в таблице: "/dev/disk3s6 971350180 11534356 712760844 2% /System/Volumes/VM".
// В первом и последнем поле могут быть пробелы. Вместо чисел могут быть прочерки.
var reSpace = regexp.MustCompile(`^(.+?)\s+(\d+|-)\s+(\d+|-)\s+(\d+|-)\s+(\d+%|-)\s+(.+)$`)

func (c *SpaceCollector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("df", "-kI")
}

func (c *SpaceCollector) ParseCommandOutput(output string) (any, error) {
	var err error

	filesystems := make(map[model.Filesystem]*model.FilesystemStats)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Filesystem") || len(line) == 0 {
			continue
		}

		matches := reSpace.FindStringSubmatch(line)
		if len(matches) != 7 {
			return nil, model.ErrStatsNotValid
		}

		var used, available float64

		used, err = strconv.ParseFloat(matches[3], 64)
		if err != nil {
			return nil, err
		}

		available, err = strconv.ParseFloat(matches[4], 64)
		if err != nil {
			return nil, err
		}

		fs := model.Filesystem{Name: matches[1], MountedOn: matches[6]}

		filesystems[fs] = &model.FilesystemStats{
			Used:  used / 1024,
			Total: (used + available) / 1024,
		}
	}
	return &model.FilesystemsMbStats{Fs: filesystems}, nil
}
