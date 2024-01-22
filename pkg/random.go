package pkg

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "1234567890"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
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

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomPhone() string {
	var sb strings.Builder
	k := len(numbers)
	sb.WriteByte('+')

	for i := 0; i < 11; i++ {
		c := numbers[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}