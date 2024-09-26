//go:build linux
// +build linux

package la

import (
	"regexp"
	"strconv"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

// Пример строки с Load Average: "top - 12:03:01 up 10 days, 18:56,  0 user,  load average: 0.20, 0.18, 0.18".
var re = regexp.MustCompile(`load average: ([\d\.]+), ([\d\.]+), ([\d\.]+)`)

func (c *Collector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("top", "-b", "-n", "1")
}

func (c *Collector) ParseCommandOutput(output string) (any, error) {
	var err error

	matches := re.FindStringSubmatch(output)

	if len(matches) != 4 {
		return nil, model.ErrStatsNotValid
	}

	var loadAvg1, loadAvg5, loadAvg15 float64

	loadAvg1, err = parseLoadAvg(matches[1])
	if err != nil {
		return nil, err
	}

	loadAvg5, err = parseLoadAvg(matches[2])
	if err != nil {
		return nil, err
	}

	loadAvg15, err = parseLoadAvg(matches[3])
	if err != nil {
		return nil, err
	}

	return &model.LoadAvgStats{
		LoadAvg1:  loadAvg1,
		LoadAvg5:  loadAvg5,
		LoadAvg15: loadAvg15,
	}, nil
}

func parseLoadAvg(loadAvgStr string) (float64, error) {
	return strconv.ParseFloat(loadAvgStr, 64)
}
