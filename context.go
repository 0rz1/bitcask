package bitcask

type context struct {
	path          string
	max_filesize  int
	max_keysize   int
	max_valuesize int
}
