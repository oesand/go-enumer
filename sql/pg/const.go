package pg

import "hash/fnv"

func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func HashString(s string) uint32 {
	return hashString(s)
}
