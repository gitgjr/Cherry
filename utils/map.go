package utils

import "math/rand"

func RandMapKey(m map[string]int) string {
	r := rand.Intn(len(m))
	for k := range m {
		if r == 0 {
			return k
		}
		r--
	}
	panic("unreachable")
}
