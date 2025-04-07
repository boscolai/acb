package acb

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateAcbCsv(t *testing.T) {
	const file = "test-data/multi-year.csv"
	require.NoError(t, GenerateAcbCsv(file, ';'))
}
