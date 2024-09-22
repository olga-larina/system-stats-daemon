package settings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSettingsService(t *testing.T) {
	t.Parallel()

	service := NewService()

	max, exists := service.GetMax()
	require.Equal(t, uint32(0), max)
	require.False(t, exists)

	insertValues := []uint32{2, 12, 10, 3, 1, 53, 1, 10}
	for _, val := range insertValues {
		service.Add(val)
	}

	max, exists = service.GetMax()
	require.Equal(t, uint32(53), max)
	require.True(t, exists)

	deleted := service.Remove(53)
	require.True(t, deleted)

	max, exists = service.GetMax()
	require.Equal(t, uint32(12), max)
	require.True(t, exists)

	deleted = service.Remove(12)
	require.True(t, deleted)
	deleted = service.Remove(10)
	require.True(t, deleted)

	max, exists = service.GetMax()
	require.Equal(t, uint32(10), max)
	require.True(t, exists)

	deleted = service.Remove(24)
	require.False(t, deleted)
}
