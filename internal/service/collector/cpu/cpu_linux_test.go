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
func TestCPULinux(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewCollector(mockedExecutor)

	commandOutput := `
top - 13:16:29 up 11 days, 12:28,  0 user,  load average: 0.14, 0.13, 0.10
Tasks:   2 total,   1 running,   1 sleeping,   0 stopped,   0 zombie
%Cpu(s):  2.1 us,  3.2 sy,  0.0 ni, 94.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :   7839.6 total,    828.1 free,    662.7 used,   6557.2 buff/cache
MiB Swap:   1024.0 total,   1024.0 free,      0.0 used.   7176.8 avail Mem

    PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
      1 root      20   0    4296   3548   3036 S   0.0   0.0   0:00.06 bash
      9 root      20   0    8488   4536   2616 R   0.0   0.1   0:00.04 top`

	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)
	require.Equal(t, &model.CPUStats{
		UserMode:   2.1,
		SystemMode: 3.2,
		Idle:       94.7,
	}, result)
}
