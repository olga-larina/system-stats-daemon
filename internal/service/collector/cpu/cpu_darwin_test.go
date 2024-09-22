package cpu

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
func TestCPUDarwin(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewCollector(mockedExecutor)

	commandOutput := `
Processes: 437 total, 3 running, 434 sleeping, 3518 threads
2024/09/16 16:49:25
Load Avg: 3.40, 3.03, 2.78
CPU usage: 9.33% user, 10.88% sys, 79.77% idle
SharedLibs: 471M resident, 104M data, 36M linkedit.`

	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)
	require.Equal(t, &model.CPUStats{
		UserMode:   9.33,
		SystemMode: 10.88,
		Idle:       79.77,
	}, result)
}
