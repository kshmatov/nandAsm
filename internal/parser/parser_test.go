package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMemLabel(t *testing.T) {
	s := "@2"
	i, err := addMemLabel(s)
	require.Nil(t, err)
	require.EqualValues(t, 2, i)
	i, err = getMemLabel(s)
	require.Nil(t, err)
	require.EqualValues(t, 2, i)
	v := formatAOp(i)
	require.EqualValues(t, 2, v)

	s = "@last"
	i, err = addMemLabel(s)
	require.EqualValues(t, 16, i)
	require.Nil(t, err)
	i, err = getMemLabel(s)
	require.Nil(t, err)
	require.EqualValues(t, 16, i)

	s = "@lastM"
	i, err = getMemLabel(s)
	require.Nil(t, err)
	require.EqualValues(t, 17, i)
}

func TestJumpLabel(t *testing.T) {
	s := []string{"(L)", "@L"}
	firstPass(s)
	require.EqualValues(t, 0, symbolTable["L"])
}

func TestOpAndSrc(t *testing.T) {
	sr, err := getSrc('M')
	require.Nil(t, err)
	require.EqualValuesf(t, 0b0001000000000000, sr, "%016b", sr)
	s := "M-1"
	op, sr, err := extractOpAndSrc(s)
	require.Nil(t, err)
	require.EqualValues(t, 0b0001000000000000, sr)
	require.EqualValues(t, 0b110010000000, op)
}
