package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый Item
	Back() *listItem                   // последний Item
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(item *listItem)             // удалить элемент
	MoveToFront(item *listItem)        // переместить элемент в начало
}

type listItem struct {
	Prev  *listItem
	Next  *listItem
	Value interface{}
}

type list struct {
	head *listItem
	tail *listItem

	length int
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *listItem {
	return l.head
}

func (l *list) Back() *listItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *listItem {
	node := &listItem{
		Value: v,
		Prev:  l.head,
		Next:  nil,
	}

	if !l.initRange(node) {
		l.head.Next = node
	}

	l.head = node
	l.length++

	return node
}

func (l *list) PushBack(v interface{}) *listItem {
	node := &listItem{
		Value: v,
		Prev:  nil,
		Next:  l.tail,
	}

	if !l.initRange(node) {
		l.tail.Prev = node
	}

	l.tail = node
	l.length++

	return node
}

func (l *list) Remove(item *listItem) {
	if l.length == 0 || item == nil {
		return
	}

	// is a head
	if item == l.head {
		l.head = l.head.Prev
		l.length--
		return
	}

	// is a tail
	if item == l.tail {
		l.tail = l.tail.Next
		l.length--
		return
	}

	// not equal head/tail of current list
	if item.Prev == nil || item.Next == nil {
		return
	}

	item.Prev.Next = item.Next
	item.Next.Prev = item.Prev
	l.length--
}

func (l *list) MoveToFront(item *listItem) {
	if item == nil {
		return
	}

	l.Remove(item)
	l.PushFront(item.Value)
}

func (l *list) initRange(item *listItem) bool {

	if l.length == 0 {
		l.head = item
		l.tail = item

		return true
	}

	return false
}
