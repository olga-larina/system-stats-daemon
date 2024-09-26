//go:build darwin
// +build darwin

package fs

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

// Пример строки в таблице: "/dev/disk3s6 971350180 11534356 712765240 2% 11 7127652400 0% /System/Volumes/VM".
// В первом и последнем поле могут быть пробелы. Вместо чисел могут быть прочерки.

//nolint:lll
var reInode = regexp.MustCompile(`^(.+?)\s+(\d+|-)\s+(\d+|-)\s+(\d+|-)\s+(\d+%|-)\s+(\d+|-)\s+(\d+|-)\s+(\d+%|-)\s+(.+)$`)

func (c *InodeCollector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("df", "-i")
}

func (c *InodeCollector) ParseCommandOutput(output string) (any, error) {
	var err error

	filesystems := make(map[model.Filesystem]*model.FilesystemStats)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Filesystem") || len(line) == 0 {
			continue
		}

		matches := reInode.FindStringSubmatch(line)
		if len(matches) != 10 {
			return nil, model.ErrStatsNotValid
		}

		var used, available float64

		used, err = strconv.ParseFloat(matches[6], 64)
		if err != nil {
			return nil, err
		}

		available, err = strconv.ParseFloat(matches[7], 64)
		if err != nil {
			return nil, err
		}

		fs := model.Filesystem{Name: matches[1], MountedOn: matches[9]}

		filesystems[fs] = &model.FilesystemStats{
			Used:  used,
			Total: used + available,
		}
	}
	return &model.FilesystemsInodeStats{Fs: filesystems}, nil
}
