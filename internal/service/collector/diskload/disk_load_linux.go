//go:build linux
// +build linux

package diskload

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

func (c *Collector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("iostat", "-d")
}

func (c *Collector) ParseCommandOutput(output string) (any, error) {
	disksLoad := make(map[string]*model.DiskLoad)

	scanner := bufio.NewScanner(strings.NewReader(output))
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)

		if lineCount >= 3 && len(fields) < 8 && len(fields) > 0 {
			return nil, model.ErrStatsNotValid
		}

		if lineCount >= 3 && len(fields) > 0 {
			disk := fields[0]
			tps, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			kbReadPS, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, err
			}
			kbWrtnPS, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, err
			}
			disksLoad[disk] = &model.DiskLoad{
				Tps: tps,
				Kbs: kbReadPS + kbWrtnPS,
			}
		}

		lineCount++
	}

	return &model.DisksLoadStats{
		Disks: disksLoad,
	}, nil
}
