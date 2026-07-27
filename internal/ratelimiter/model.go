package ratelimiter

import (
	"time"
)
type OtpRequestPerEmail struct {
	NumberOfRequests int
}

var RedisKeyPrefix string = "ratelimiter:otp:peremail"
var TTL = 10*time.Minute
var MaxReqPerTTL = 5