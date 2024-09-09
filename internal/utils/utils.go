package utils

import "strconv"

func ItoaSlice(s []int) []string {
	var r []string
	for i := range s {
		number := s[i]
		text := strconv.Itoa(number)
		r = append(r, text)
	}
	return r
}
