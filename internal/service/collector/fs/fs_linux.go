//go:build linux
// +build linux

package fs

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

// Пример строки в таблице: "tmpfs 1003464 1 1003463 1% /sys/firmware".
// В первом и последнем поле могут быть пробелы. Вместо чисел могут быть прочерки.
var re = regexp.MustCompile(`^(.+?)\s+(\d+|-)\s+(\d+|-)\s+(\d+|-)\s+(\d+%|-)\s+(.+)$`)

func parseDfOutput(output string, divider float64) (map[model.Filesystem]*model.FilesystemStats, error) {
	var err error

	filesystems := make(map[model.Filesystem]*model.FilesystemStats)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Filesystem") || len(line) == 0 {
			continue
		}

		matches := re.FindStringSubmatch(line)
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
			Used:  used / divider,
			Total: (used + available) / divider,
		}
	}

	return filesystems, nil
}
