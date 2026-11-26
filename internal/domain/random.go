package domain

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

var taskType = [3]string{"once", "interval", "cron"}

func RandomTaskType() string {
	return taskType[rand.Intn(len(taskType))]
}

var MethodType = [4]HTTPMethod{
	MethodGet,
	MethodPost,
	MethodPut,
	MethodDelete,
}

func RandomMethod() HTTPMethod {
	return MethodType[rand.Intn(len(MethodType))]
}
