package store

func FindInSlice(ss []string, key string) int {
	for i, s := range ss {
		if s == key {
			return i
		}
	}
	return -1
}

func DeleteFromSlice(ss []string, ind int) []string {
	if ind < 0 || ind >= len(ss) {
		return ss
	}
	return append(ss[0:ind], ss[ind+1:]...)
}
