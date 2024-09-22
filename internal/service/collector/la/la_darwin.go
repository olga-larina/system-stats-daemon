//go:build darwin
// +build darwin

package la

import (
	"strconv"
	"strings"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

func (c *Collector) ExecuteCommand() ([]byte, error) {
	// Пример строки: "15:28  up 51 days,  8:03, 3 users, load averages: 1,64 1,39 1,48"
	return c.executor.Execute("uptime")
}

func (c *Collector) ParseCommandOutput(output string) (any, error) {
	var err error

	fields := strings.Fields(output)

	if len(fields) < 3 {
		return nil, model.ErrStatsNotValid
	}

	var loadAvg1, loadAvg5, loadAvg15 float64

	loadAvg1, err = parseLoadAvg(fields[len(fields)-3])
	if err != nil {
		return nil, err
	}

	loadAvg5, err = parseLoadAvg(fields[len(fields)-2])
	if err != nil {
		return nil, err
	}

	loadAvg15, err = parseLoadAvg(fields[len(fields)-1])
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
	return strconv.ParseFloat(strings.ReplaceAll(loadAvgStr, ",", "."), 64)
}
