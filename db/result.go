package db

type result struct {
	err   error
	value interface{}
}

func newResult(v interface{}, err error) *result {
	return &result{
		value: v,
		err:   err,
	}
}
