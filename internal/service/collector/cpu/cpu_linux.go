//go:build linux
// +build linux

package cpu

import (
	"regexp"
	"strconv"

	"github.com/olga-larina/system-stats-daemon/internal/model"
)

// Пример строки с CPU: "%Cpu(s):  0.0 us,  2.4 sy,  0.0 ni, 97.6 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st".
var re = regexp.MustCompile(`%Cpu\(s\):\s*([\d\.]+)\s*us,\s*([\d\.]+)\s*sy,\s*([\d\.]+)\s*ni,\s*([\d\.]+)\s*id`)

func (c *Collector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("top", "-b", "-n", "1")
}

func (c *Collector) ParseCommandOutput(output string) (any, error) {
	var err error

	matches := re.FindStringSubmatch(output)

	if len(matches) != 5 {
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

	idle, err = parseCPU(matches[4])
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
