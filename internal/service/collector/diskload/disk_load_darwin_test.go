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
func TestDiskLoadDarwin(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewCollector(mockedExecutor)

	commandOutput := `	              disk0               disk4
	    KB/t  tps  MB/s     KB/t  tps  MB/s
	   19.42   62  1.17    52.69    1  2.34`

	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)

	require.Equal(t, &model.DisksLoadStats{
		Disks: map[string]*model.DiskLoad{
			"disk0": {Tps: 62, Kbs: 1198.08},
			"disk4": {Tps: 1, Kbs: 2396.16},
		},
	}, result)
}
