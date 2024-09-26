//nolint:dupl
package fs

import (
	"testing"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/mocks"
	"github.com/stretchr/testify/require"
)

const epsilonSpace = 1e-6

/*
Генерация моков:

	go install github.com/vektra/mockery/v2@v2.43.2

	mockery --all --case underscore --keeptree --dir internal/service/collector \
		--output internal/service/collector/mocks --with-expecter --log-level warn.
*/
func TestFsSpaceDarwin(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewSpaceCollector(mockedExecutor)

	commandOutput := `
Filesystem     1024-blocks      Used Available Capacity  Mounted on
/dev/disk3s1s1   971350180  13020840 712760844     2%    /
map auto_home            0         0         0   100%    /System/Volumes/Data/home
/dev/disk4s1        662488    423564    238924    64%    /Volumes/MongoDB Compass
`
	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)

	actual := result.(*model.FilesystemsMbStats)

	expected := map[model.Filesystem]*model.FilesystemStats{
		{Name: "/dev/disk3s1s1", MountedOn: "/"}:                        {Used: 12715.6640625, Total: 708771.17578125},
		{Name: "map auto_home", MountedOn: "/System/Volumes/Data/home"}: {Used: 0, Total: 0},
		{Name: "/dev/disk4s1", MountedOn: "/Volumes/MongoDB Compass"}:   {Used: 413.63671875, Total: 646.9609375},
	}
	for fs, stats := range expected {
		actuaFs := actual.Fs[fs]
		require.NotNil(t, actuaFs)
		require.InDelta(t, stats.Used, actuaFs.Used, epsilonSpace)
		require.InDelta(t, stats.Total, actuaFs.Total, epsilonSpace)
	}
}
