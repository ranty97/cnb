package utils

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func ItoaSlice(s []int) []string {
	var r []string
	for i := range s {
		number := s[i]
		text := strconv.Itoa(number)
		r = append(r, text)
	}
	return r
}

func LastCharacterAsNumber(s string) (int, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("empty string")
	}

	lastChar := s[len(s)-1]
	num, err := strconv.Atoi(string(lastChar))

	if err != nil {
		return 0, fmt.Errorf("last character '%c' not a number", lastChar)
	}

	return num, nil
}

func InvertRandomBitWithProbability(data []byte, probability float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if rand.Float64() < probability {
		randomByte := r.Intn(len(data))
		randomBit := r.Intn(8)

		data[randomByte] ^= 1 << randomBit

		log.Printf("Inverted bit %d in byte %d\n", randomBit, randomByte)
	}
}
