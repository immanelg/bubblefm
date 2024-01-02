package main

type set[T comparable] map[T]struct{}

func (self set[T]) Add(v T) {
	self[v] = struct{}{}
}

func (self set[T]) Toggle(v T) {
	if self.Contains(v) {
		self.Remove(v)
	} else {
		self.Add(v)
	}
}

func (self set[T]) Remove(v T) {
	delete(self, v)
}

func (self set[T]) Contains(v T) bool {
	if _, ok := self[v]; ok {
		return true
	}
	return false
}

// non-strict superset
func (self set[T]) Superset(other []T) bool {
	for _, v := range other {
		if !self.Contains(v) {
			return false
		}
	}
	return true
}

func (self set[T]) Clear() {
	clear(self)
}
