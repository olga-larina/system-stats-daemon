//go:build darwin
// +build darwin

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

	var disks []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		switch {
		case lineCount == 0 && len(fields) > 0:
			disks = fields
		case lineCount == 2 && len(fields) > 0:
			for i := 0; i < len(disks); i++ {
				tpsIndex := 1 + i*3
				mbpsIndex := 2 + i*3

				if tpsIndex < len(fields) && mbpsIndex < len(fields) {
					tps, err := strconv.ParseFloat(fields[tpsIndex], 64)
					if err != nil {
						return nil, err
					}
					mbps, err := strconv.ParseFloat(fields[mbpsIndex], 64)
					if err != nil {
						return nil, err
					}

					disksLoad[disks[i]] = &model.DiskLoad{
						Tps: tps,
						Kbs: mbps * 1024,
					}
				} else {
					return nil, model.ErrStatsNotValid
				}
			}
		case len(fields) == 0 && (lineCount == 1 || lineCount == 3):
			return nil, model.ErrStatsNotValid
		}
		lineCount++
	}

	return &model.DisksLoadStats{
		Disks: disksLoad,
	}, nil
}
