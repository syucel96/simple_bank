package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min, max := int64(10), int64(1000)
	randInt := RandomInt(min, max)
	require.NotEmpty(t, randInt)
	require.IsType(t, max, randInt)
	require.GreaterOrEqual(t, randInt, min)
	require.LessOrEqual(t, randInt, max)
}

func TestRandomLength(t *testing.T) {
	randLen := randomLength()
	require.NotEmpty(t, randLen)
	require.IsType(t, int(1), randLen)
	require.GreaterOrEqual(t, randLen, 0)
	require.LessOrEqual(t, randLen, 18)

	randLen = randomLength(100)
	require.IsType(t, int(1), randLen)
	require.NotEmpty(t, randLen)
	require.GreaterOrEqual(t, randLen, 0)
	require.LessOrEqual(t, randLen, 100)
}

func TestRandomString(t *testing.T) {
	randLen := randomLength()
	randStr := randomString(randLen)
	require.NotEmpty(t, randStr)
	require.IsType(t, "", randStr)
	require.Len(t, randStr, randLen)
}

func TestRandomFloat(t *testing.T) {
	min, max := int64(10), int64(1000)
	randFloat := RandomFloat(min, max)
	require.NotEmpty(t, randFloat)
	require.IsType(t, float64(1), randFloat)
	require.GreaterOrEqual(t, randFloat, float64(min))
	require.LessOrEqual(t, randFloat, float64(max))
}

func TestRandomOwner(t *testing.T) {
	randStr := RandomOwner()
	require.NotEmpty(t, randStr)
	require.IsType(t, "", randStr)
	require.Len(t, randStr, 6)
}

func TestRandBool(t *testing.T) {
	prevCache := cache
	rem := remaining
	randomBool := randBool()
	require.IsType(t, true, randomBool)
	if rem > 1 {
		require.Equal(t, rem-1, remaining)
		prevCache >>= 1
		require.Equal(t, prevCache, cache)
	} else {
		require.Equal(t, 63, remaining)
	}
}

func TestRandomMoney(t *testing.T) {
	randMoney := RandomMoney(true)
	require.True(t, IsDecimal(randMoney))
	require.IsType(t, "", randMoney)

	randMoney = RandomMoney(false)
	require.True(t, IsDecimal(randMoney))
	require.IsType(t, "", randMoney)
	require.GreaterOrEqual(t, ParseFloat(randMoney), float64(0))

	min, max := int64(10), int64(1000)
	randMoney = RandomMoney(false, min, max)
	require.True(t, IsDecimal(randMoney))
	require.IsType(t, "", randMoney)
	require.GreaterOrEqual(t, ParseFloat(randMoney), float64(min))
	require.LessOrEqual(t, ParseFloat(randMoney), float64(max))

	randMoney = RandomMoney(false, min)
	require.True(t, IsDecimal(randMoney))
	require.IsType(t, "", randMoney)
	require.LessOrEqual(t, ParseFloat(randMoney), float64(max))
}

func TestRandomCurrency(t *testing.T) {
	randCurr := RandomCurrency()
	require.NotEmpty(t, randCurr)
	require.IsType(t, "", randCurr)
	require.True(t, IsSupportedCurrency(randCurr))
	require.False(t, IsSupportedCurrency("ABX"))
}

func TestRandomFullName(t *testing.T) {
	randName := RandomFullName()
	require.NotEmpty(t, randName)
	require.IsType(t, "", randName)
	require.Contains(t, randName, " ")
}

func TestRandomEmail(t *testing.T) {
	randEmail := RandomEmail()
	require.NotEmpty(t, randEmail)
	require.IsType(t, "", randEmail)
	require.Contains(t, randEmail, "@")
	require.Contains(t, randEmail, ".com")
}
