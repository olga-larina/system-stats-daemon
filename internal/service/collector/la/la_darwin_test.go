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
func TestLADarwin(t *testing.T) {
	t.Parallel()

	mockedExecutor := mocks.NewCommandExecutor(t)
	collector := NewCollector(mockedExecutor)

	commandOutput := `
16:53  up 53 days,  9:28, 3 users, load averages: 2,40 2,91 2,80`

	result, err := collector.ParseCommandOutput(commandOutput)
	require.NoError(t, err)

	require.Equal(t, &model.LoadAvgStats{
		LoadAvg1:  2.4,
		LoadAvg5:  2.91,
		LoadAvg15: 2.8,
	}, result)
}
