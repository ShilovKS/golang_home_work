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
	front  *ListItem // первый элемент списка
	back   *ListItem // последний элемент списка
	length int       // длина списка
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.Front(),
	}

	if l.front == nil {
		l.back = newItem
	} else {
		l.front.Prev = newItem
	}

	l.front = newItem
	l.length++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  l.Back(),
	}

	if l.back == nil {
		l.front = newItem
	} else {
		l.back.Next = newItem
	}
	l.back = newItem
	l.length++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i == l.front {
		l.front = i.Next
	} else if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i == l.back {
		l.back = i.Prev
	} else if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}

	if i == l.back {
		l.back = i.Prev
		if l.back != nil {
			l.back.Next = nil
		}
	} else {
		if i.Prev != nil {
			i.Prev.Next = i.Next
		}
		if i.Next != nil {
			i.Next.Prev = i.Prev
		}
	}

	i.Prev = nil
	i.Next = l.front
	if l.front != nil {
		l.front.Prev = i
	}
	l.front = i

	if l.back == nil {
		l.back = i
	}
}

func NewList() List {
	return new(list)
}
