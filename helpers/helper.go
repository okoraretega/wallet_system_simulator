package helpers

import (
	"fmt"
	"math/rand"
)

func GenerateWalletNumber() string {
	num := rand.Intn(9_000_000_000) + 1_000_000_000 // always 10 digits
	return fmt.Sprintf("%d", num)
}
