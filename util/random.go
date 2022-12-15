package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var cache int64
var remaining uint8

func init() {
	rand.Seed(time.Now().UnixNano())
	cache = rand.Int63()
	remaining = 63
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomFloat(min, max int64) float64 {
	return NewFloat(RandomInt(min, max), rand.Float64())
}

func RandomOwner() string {
	return randomString(6)
}

func randBool() bool {
	res := cache&0x01 == 1
	cache >>= 1
	remaining--
	if remaining == 0 {
		cache = rand.Int63()
		remaining = 63
	}
	return res
}

func RandomMoney(signed bool, minMax ...int64) string {
	var sign string = ""
	if signed && !randBool() {
		sign = "-"
	}

	var floor, ceiling int64 = 0, 1000000000
	if len(minMax) > 1 {
		floor = minMax[0]
		ceiling = minMax[1]
	} else if len(minMax) > 0 {
		ceiling = minMax[0]
	}
	return fmt.Sprintf("%s%.2f", sign, RandomFloat(floor, ceiling))
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, GBP, CAD, JPY, TRY}

	return currencies[rand.Intn(len(currencies))]
}
