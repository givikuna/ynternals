package fp

func Ternary[T any](cond bool, onTrue, onFalse T) T {
	if cond {
		return onTrue
	}
	return onFalse
}

func Map[A, B any](items []A, fn func(A) B) []B {
	if len(items) == 0 {
		return nil
	}
	head := fn(items[0])
	tail := Map(items[1:], fn)
	return append([]B{head}, tail...)
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	if len(items) == 0 {
		return nil
	}
	rest := Filter(items[1:], predicate)
	return Ternary(predicate(items[0]), append([]T{items[0]}, rest...), rest)
}

func Fold[T, Acc any](items []T, initial Acc, folder func(Acc, T) Acc) Acc {
	if len(items) == 0 {
		return initial
	}
	return Fold(items[1:], folder(initial, items[0]), folder)
}

func Every[T any](items []T, predicate func(T) bool) bool {
	if len(items) == 0 {
		return true
	}
	return predicate(items[0]) && Every(items[1:], predicate)
}

func ForEachErr[T any](items []T, fn func(T) error) error {
	if len(items) == 0 {
		return nil
	}
	if err := fn(items[0]); err != nil {
		return err
	}
	return ForEachErr(items[1:], fn)
}
