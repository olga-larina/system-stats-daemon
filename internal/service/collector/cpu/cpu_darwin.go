//go:build darwin
// +build darwin

package cpu

import (
	"regexp"
	"strconv"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

// Пример строки с CPU: "CPU usage: 2.21% user, 7.47% sys, 90.30% idle".
var re = regexp.MustCompile(`CPU usage: ([\d\.]+)% user, ([\d\.]+)% sys, ([\d\.]+)% idle`)

func (c *Collector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("bash", "-c", "top -l 1 | head -n 5")
}

func (c *Collector) ParseCommandOutput(output string) (any, error) {
	var err error

	matches := re.FindStringSubmatch(output)

	if len(matches) != 4 {
		return nil, model.ErrStatsNotValid
	}

	var userMode, systemMode, idle float64

	userMode, err = parseCPU(matches[1])
	if err != nil {
		return nil, err
	}

	systemMode, err = parseCPU(matches[2])
	if err != nil {
		return nil, err
	}

	idle, err = parseCPU(matches[3])
	if err != nil {
		return nil, err
	}

	return &model.CPUStats{
		UserMode:   userMode,
		SystemMode: systemMode,
		Idle:       idle,
	}, nil
}

func parseCPU(cpuStr string) (float64, error) {
	return strconv.ParseFloat(cpuStr, 64)
}
