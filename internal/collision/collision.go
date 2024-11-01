package collision

import (
	"math/rand"
	"time"
)

func RandomlyAddCollision(data []byte) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	collisionalData := append(data, []byte{'a', 'b', 'c'}...)

	if r.Float64() < 0.6 {
		return collisionalData
	}

	return data
}
