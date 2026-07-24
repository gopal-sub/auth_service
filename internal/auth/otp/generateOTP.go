package otp

import (
	"strconv"
	"math/rand/v2"
)

func GenerateOPT() string{
	random_number := rand.IntN(999999) + 10000
	return strconv.Itoa(random_number)
}