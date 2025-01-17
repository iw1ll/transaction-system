package utils

import (
	"fmt"
	"math/rand"
)

func GenerateRandomAddress() string {
	return fmt.Sprintf("%x", rand.Int63())
}
