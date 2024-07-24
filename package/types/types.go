package types

type RequestType[T any] struct {
	data T
}

func NewRequestType[T any](data T) RequestType[T] {
	return RequestType[T]{
		data: data,
	}
}

func (r RequestType[T]) Data() T {
	return r.data
}

type ResponseType[T any] struct {
	data T
}

func NewResponseType[T any](data T) ResponseType[T] {
	return ResponseType[T]{
		data: data,
	}
}

func (r ResponseType[T]) Data() T {
	return r.data
}
