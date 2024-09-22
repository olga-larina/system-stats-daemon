//nolint:dupl
package fs

import (
	"testing"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/mocks"
	"github.com/stretchr/testify/require"
)

const epsilonInode = 1e-6

/*
Генерация моков:

	go install github.com/vektra/mockery/v2@v2.43.2

	mockery --all --case underscore --keeptree --dir internal/service/collector \
		--output internal/service/collector/mocks --with-expecter --log-level warn.
*/
func TestFsInodeDarwin(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewInodeCollector(mockedExecutor)

	commandOutput := `
Filesystem     512-blocks      Used  Available Capacity iused      ifree %iused  Mounted on
/dev/disk3s1s1 1942700360  26041680 1425546976     2%  404167 4293779928    0%   /
map auto_home           0         0          0   100%       0          0     -   /System/Volumes/Data/home
/dev/disk4s1      1324976    847128     477848    64%     611 4294966668    0%   /Volumes/MongoDB Compass
`
	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)

	actual := result.(*model.FilesystemsInodeStats)

	expected := map[model.Filesystem]*model.FilesystemStats{
		{Name: "/dev/disk3s1s1", MountedOn: "/"}:                        {Used: 404167, Total: 4294184095},
		{Name: "map auto_home", MountedOn: "/System/Volumes/Data/home"}: {Used: 0, Total: 0},
		{Name: "/dev/disk4s1", MountedOn: "/Volumes/MongoDB Compass"}:   {Used: 611, Total: 4294967279},
	}
	for fs, stats := range expected {
		actuaFs := actual.Fs[fs]
		require.NotNil(t, actuaFs)
		require.InDelta(t, stats.Used, actuaFs.Used, epsilonInode)
		require.InDelta(t, stats.Total, actuaFs.Total, epsilonInode)
	}
}
