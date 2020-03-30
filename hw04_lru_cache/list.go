package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый Item
	Back() *listItem                   // последний Item
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	len   int
	front *listItem
	back  *listItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	newItem := &listItem{
		Value: v,
	}
	if l.Len() == 0 {
		l.back = newItem
	} else {
		newItem.Prev = l.front
		l.front.Next = newItem
	}
	l.front = newItem
	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *listItem {
	newItem := &listItem{
		Value: v,
	}
	if l.Len() == 0 {
		l.front = newItem
	} else {
		newItem.Next = l.back
		l.back.Prev = newItem
	}
	l.back = newItem
	l.len++
	return l.front
}

func (l *list) Remove(i *listItem) {
	if l.len == 0 {
		return
	}
	if i.Prev == nil {
		l.back = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.front = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	if i.Next == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.back = i.Next
	}
	i.Next.Prev = i.Prev

	i.Prev = l.front
	i.Next = nil
	l.front.Next = i
	l.front = i
}

func NewList() List {
	return &list{}
}
