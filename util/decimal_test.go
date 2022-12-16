package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFloat(t *testing.T) {
	fl := NewFloat(int64(1100), float64(0.77777777))
	require.NotZero(t, fl)
	require.IsType(t, float64(1), fl)
	require.Equal(t, float64(1100.78), fl)
}

func TestIsDecimal(t *testing.T) {
	require.False(t, IsDecimal("a"))
	require.True(t, IsDecimal(RandomMoney(true)))
}

func TestParseFloat(t *testing.T) {
	require.Zero(t, ParseFloat("a"))

	randFloat := RandomFloat(int64(10), int64(100000))
	require.Equal(t, randFloat, ParseFloat(fmt.Sprintf("%v", randFloat)))
}

func TestChangeDecimalSign(t *testing.T) {
	randDec := RandomFloat(int64(10), int64(100000))
	require.Equal(t, -randDec, ParseFloat(ChangeDecimalSign(fmt.Sprintf("%v", randDec))))
	require.Equal(t, randDec, ParseFloat(ChangeDecimalSign(fmt.Sprintf("%v", -randDec))))
}

func TestAddDecimals(t *testing.T) {
	randDec1 := RandomFloat(int64(10), int64(100000))
	randDec2 := RandomFloat(int64(10), int64(100000))

	str1 := fmt.Sprintf("%v", randDec1)
	str2 := fmt.Sprintf("%v", randDec2)
	str3 := fmt.Sprintf("%v", -randDec2)

	sum := randDec1 + randDec2
	sub := randDec1 - randDec2

	added, err := AddDecimals(str1, str2)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%0.2f", sum), added)

	subbed, err := AddDecimals(str1, str3)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%0.2f", sub), subbed)

	added, err = AddDecimals("a", str1)
	require.Error(t, err)
	require.Equal(t, added, "")

	added, err = AddDecimals(str1, "")
	require.Error(t, err)
	require.Equal(t, added, "")
}

func TestSubtractDecimals(t *testing.T) {
	randDec1 := RandomFloat(int64(10), int64(100000))
	randDec2 := RandomFloat(int64(10), int64(100000))

	str1 := fmt.Sprintf("%v", randDec1)
	str2 := fmt.Sprintf("%v", randDec2)
	str3 := fmt.Sprintf("%v", -randDec2)

	sum := randDec1 + randDec2
	sub := randDec1 - randDec2

	subbed, err := SubtractDecimals(str1, str2)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%0.2f", sub), subbed)

	added, err := SubtractDecimals(str1, str3)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%0.2f", sum), added)

	subbed, err = SubtractDecimals("a", str1)
	require.Error(t, err)
	require.Equal(t, subbed, "")

	subbed, err = SubtractDecimals(str1, "")
	require.Error(t, err)
	require.Equal(t, subbed, "")
}
