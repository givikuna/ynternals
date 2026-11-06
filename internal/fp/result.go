package fp

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](val T) Result[T] {
	return Result[T]{value: val}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func MapResult[A, B any](r Result[A], fn func(A) B) Result[B] {
	if r.err != nil {
		return Err[B](r.err)
	}
	return Ok(fn(r.value))
}

func FlatMap[A, B any](r Result[A], fn func(A) Result[B]) Result[B] {
	if r.err != nil {
		return Err[B](r.err)
	}
	return fn(r.value)
}
