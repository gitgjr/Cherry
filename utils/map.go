package utils

import (
	"math/rand"
)

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

type hashMap interface {
	MergeHashMaps()
	GetAllKeys()
}

//TODO: Design a func to return all keys whatever the dtype of key
// func GetAllKeys(m1 map[any]any) []any {
// 	//use reflect
// 	var dtype any
// 	for k, v := range m1 {
// 		dtype=reflect.TypeOf(k)
// 	}
// }
