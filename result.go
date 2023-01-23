package main

type Result[T any] struct {
	v T
	e error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		v: value,
		e: nil,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		e: err,
	}
}

func (self Result[T]) Unwrap() T {
	if self.e != nil {
		panic(self.e)
	}
	return self.v
}

func (self Result[T]) Expect(msg string) T {
	if self.e != nil {
		panic(msg)
	}
	return self.v
}

func (self Result[T]) IsErr() bool {
	return self.e != nil
}

func (self Result[T]) IsOk() bool {
	return self.e == nil
}

func (self Result[T]) Extract() (T, error) {
	return self.v, self.e
}

func (self Result[T]) IsOkAnd(f func(v T) bool) bool {
	isOk := self.IsOk()
	if isOk {
		return true
	}
	f(self.v)
	return false
}

func Map[T, U any](self Result[T], op func(v T) U) Result[U] {
	if self.IsErr() {
		return Err[U](self.e)
	}
	u := op(self.Unwrap())
	return Ok(u)
}
