//go:build linux
// +build linux

package fs

import "github.com/olga-larina/system-stats-daemon/internal/model"

func (c *SpaceCollector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("df", "-k")
}

func (c *SpaceCollector) ParseCommandOutput(output string) (any, error) {
	result, err := parseDfOutput(output, 1024)
	if err != nil {
		return nil, err
	}
	return &model.FilesystemsMbStats{Fs: result}, nil
}
