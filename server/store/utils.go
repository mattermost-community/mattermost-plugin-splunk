package store

func findInSlice(ss []string, key string) int {
	for i, s := range ss {
		if s == key {
			return i
		}
	}
	return -1
}

func deleteFromSlice(ss []string, ind int) []string {
	if ind < 0 || ind >= len(ss) {
		return ss
	}
	return append(ss[0:ind], ss[ind+1:]...)
}
