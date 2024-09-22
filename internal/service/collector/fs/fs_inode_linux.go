//go:build linux
// +build linux

package fs

import "github.com/olga-larina/system-stats-daemon/internal/model"

func (c *InodeCollector) ExecuteCommand() ([]byte, error) {
	return c.executor.Execute("df", "-i")
}

func (c *InodeCollector) ParseCommandOutput(output string) (any, error) {
	result, err := parseDfOutput(output, 1)
	if err != nil {
		return nil, err
	}
	return &model.FilesystemsInodeStats{Fs: result}, nil
}
