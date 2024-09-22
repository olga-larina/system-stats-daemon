package la

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
func TestLALinux(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewCollector(mockedExecutor)

	commandOutput := `
top - 13:57:46 up 11 days, 13:09,  0 user,  load average: 0.11, 0.10, 0.09
Tasks:   2 total,   1 running,   1 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.0 us,  2.1 sy,  0.0 ni, 97.9 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :   7839.6 total,    822.9 free,    667.9 used,   6557.3 buff/cache
MiB Swap:   1024.0 total,   1024.0 free,      0.0 used.   7171.7 avail Mem

    PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
      1 root      20   0    4296   3548   3036 S   0.0   0.0   0:00.07 bash
     12 root      20   0    8488   4544   2624 R   0.0   0.1   0:00.01 top`

	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)

	require.Equal(t, &model.LoadAvgStats{
		LoadAvg1:  0.11,
		LoadAvg5:  0.1,
		LoadAvg15: 0.09,
	}, result)
}
