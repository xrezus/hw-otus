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
	length    int
	firstItem *ListItem
	lastItem  *ListItem
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}

	if l.length == 0 {
		l.lastItem = item
		l.firstItem = item
		l.length++

		return item
	}

	item.Prev = nil
	item.Next = l.firstItem
	l.firstItem.Prev = item
	l.firstItem = item
	l.length++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}

	if l.length == 0 {
		l.firstItem = item
		l.lastItem = item
		l.length++

		return item
	}

	item.Prev = l.lastItem
	item.Next = nil
	l.lastItem.Next = item
	l.lastItem = item
	l.length++

	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.firstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.lastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}
