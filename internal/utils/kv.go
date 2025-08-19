package utils

type KV struct {
	Key string
	Val string
}

type KVSlice []KV

func (s KVSlice) Len() int           { return len(s) }
func (s KVSlice) Less(i, j int) bool { return s[i].Key < s[j].Key }
func (s KVSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s KVSlice) InsertionSort() {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1

		for j >= 0 && s[j].Key > key.Key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}
