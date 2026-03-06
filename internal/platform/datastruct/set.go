package datastruct

type Set[T comparable] interface {
	Add(T)
	Remove(T)
	Contains(T) bool
	Elements() []T
	Size() int
}

type HashSet[T comparable] struct {
	internalMap map[T]bool
}

func NewHashSet[T comparable](size ...int) *HashSet[T] {
	if len(size) > 0 {
		return &HashSet[T]{make(map[T]bool, size[0])}
	}

	return &HashSet[T]{make(map[T]bool)}
}

func (hs *HashSet[T]) Add(value T) {
	hs.internalMap[value] = true
}

func (hs *HashSet[T]) Remove(value T) {
	delete(hs.internalMap, value)
}

func (hs HashSet[T]) Contains(value T) (exists bool) {
	_, exists = hs.internalMap[value]
	return
}

func (hs HashSet[T]) Elements() []T {
	elements := make([]T, 0, len(hs.internalMap))

	for key := range hs.internalMap {
		elements = append(elements, key)
	}

	return elements
}

func (hs HashSet[T]) Size() int {
	return len(hs.internalMap)
}
