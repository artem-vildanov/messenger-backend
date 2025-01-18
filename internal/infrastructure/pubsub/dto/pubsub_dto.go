package dto

import "messenger/internal/infrastructure/errors"

type PubsubDto[T any] struct {
	Payload T
	Error   *errors.Error
}

func NewPubsubDto[T any](dto T) *PubsubDto[T] {
	return &PubsubDto[T]{
		Payload: dto,
		Error:   nil,
	}
}

func NewPubsubError[T any](err *errors.Error) *PubsubDto[T] {
	return &PubsubDto[T]{
		Payload: *new(T),
		Error:   err,
	}
}
