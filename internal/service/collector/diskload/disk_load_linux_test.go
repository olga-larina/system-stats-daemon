package diskload

import (
	"testing"

	"github.com/olga-larina/system-stats-daemon/internal/model"
	"github.com/olga-larina/system-stats-daemon/internal/service/collector/mocks"
	"github.com/stretchr/testify/require"
)

/*
Генерация моков:

	go install github.com/vektra/mockery/v2@v2.43.2

	mockery --all --case underscore --keeptree --dir internal/service/collector \
		--output internal/service/collector/mocks --with-expecter --log-level warn.
*/
func TestDiskLoadLinux(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewCollector(mockedExecutor)

	commandOutput := `Linux 6.10.0-linuxkit (6ece98118879) 	09/18/24 	_aarch64_	(4 CPU)

Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
vda               1.01         1.74        29.16       305.06    1853029   30977604  324058376
vdb               0.02         0.21         0.10         0.00     218380         10        100

`
	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)
	require.Equal(t, &model.DisksLoadStats{
		Disks: map[string]*model.DiskLoad{
			"vda": {Tps: 1.01, Kbs: 30.9},
			"vdb": {Tps: 0.02, Kbs: 0.31},
		},
	}, result)
}
