package app

import (
	"github.com/AkvicorEdwards/util"
	"strconv"
	"time"
)

const randomStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenerateSecret() string {
	return util.RandomString(32, randomStr) + strconv.FormatInt(time.Now().UnixNano(), 36)
}
