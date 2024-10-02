package utils

import (
	"fmt"
	"strconv"
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
		return 0, fmt.Errorf("строка пуста")
	}

	lastChar := s[len(s)-1]
	num, err := strconv.Atoi(string(lastChar))

	if err != nil {
		return 0, fmt.Errorf("последний символ '%c' не является числом", lastChar)
	}

	return num, nil
}
