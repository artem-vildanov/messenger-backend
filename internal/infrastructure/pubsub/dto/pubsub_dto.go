package dto

type PubsubDto[T any] struct {
	Payload T
	Error   error
}

func NewPubsubDto[T any](dto T) *PubsubDto[T] {
	return &PubsubDto[T]{
		Payload: dto,
		Error:   nil,
	}
}

func NewPubsubError[T any](err error) *PubsubDto[T] {
	return &PubsubDto[T]{
		Payload: *new(T),
		Error:   err,
	}
}
