package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.head = &ListItem{v, l.head, nil}
	if l.tail == nil {
		l.tail = l.head
	}
	l.len++
	if l.head.Next != nil {
		l.head.Next.Prev = l.head
	}

	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.tail = &ListItem{v, nil, l.tail}
	if l.head == nil {
		l.head = l.tail
	}
	l.len++
	if l.tail.Prev != nil {
		l.tail.Prev.Next = l.tail
	}
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.head == l.tail:
		l.head = nil
		l.tail = nil
	case i.Prev == nil:
		i.Next.Prev, l.head = i.Prev, l.head.Next
	case i.Next == nil:
		i.Prev.Next, l.tail = i.Next, l.tail.Prev
	default:
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev != nil {
		if i.Next == nil {
			l.tail = l.tail.Prev
			i.Prev.Next = nil
		} else {
			i.Prev.Next, i.Next.Prev = i.Next, i.Prev
		}
		i.Next, l.head, i.Prev, l.head.Prev = l.head, i, nil, i
	}
}

func NewList() List {
	return new(list)
}
