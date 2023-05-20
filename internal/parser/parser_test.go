package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMemLabel(t *testing.T) {
	s := "@2"
	err := addMemLabel(s)
	require.Nil(t, err)
	i, err := getMemLabel(s)
	require.Nil(t, err)
	require.EqualValues(t, 2, i)
	v := formatAOp(i)
	require.EqualValues(t, 2, v)
}
