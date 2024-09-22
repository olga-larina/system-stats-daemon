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
func TestFsInodeLinux(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewInodeCollector(mockedExecutor)

	commandOutput := `
Filesystem      Inodes  IUsed   IFree IUse% Mounted on
overlay 1       61202244 26840512  31220408  47% /
tmpfs              65568       32     65536   0% /dev
tmpfs            4013969      123   4013856   1% /sys/firmware info
`
	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)

	actual := result.(*model.FilesystemsInodeStats)

	expected := map[model.Filesystem]*model.FilesystemStats{
		{Name: "overlay 1", MountedOn: "/"}:              {Used: 26840512, Total: 58060920},
		{Name: "tmpfs", MountedOn: "/dev"}:               {Used: 32, Total: 65568},
		{Name: "tmpfs", MountedOn: "/sys/firmware info"}: {Used: 123, Total: 4013979},
	}
	for fs, stats := range expected {
		actuaFs := actual.Fs[fs]
		require.NotNil(t, actuaFs)
		require.InDelta(t, stats.Used, actuaFs.Used, epsilonInode)
		require.InDelta(t, stats.Total, actuaFs.Total, epsilonInode)
	}
}
