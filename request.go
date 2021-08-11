package bitcask

type response struct {
	err error
	ret interface{}
}

type request struct {
	param interface{}
	ch    <-chan *response
}

func asyncResponse(param interface{}, fn func(*request)) *response {
	ch := make(<-chan *response)
	req := &request{
		param: param,
		ch:    ch,
	}
	go fn(req)
	return <-ch
}
