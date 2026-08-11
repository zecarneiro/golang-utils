package models

type Response[T any] struct {
	Data  T
	Error error
}

func NewCustomResponse[T any](data T, err error) *Response[T] {
	return &Response[T]{
		Data:  data,
		Error: err,
	}
}

func NewResponse[T any]() *Response[T] {
	data := *new(T)
	return NewCustomResponse(data, nil)
}

func (r Response[T]) HasData() bool {
	data := &r.Data
	return data != nil
}

func (r Response[T]) HasError() bool {
	return r.Error != nil
}
