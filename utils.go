package eval

import (
	"container/list"
)

//Stack is a structure of normal stack
type Stack struct {
	Val *list.List
}

func NewStack() *Stack {
	s := new(Stack)
	val := list.New()
	s.Val = val
	return s
}

func (s *Stack) Push(a interface{}) {
	s.Val.PushBack(a)
	return
}

func (s *Stack) Pop() interface{} {
	ele := s.Val.Back()
	if ele == nil {
		return nil
	}
	return s.Val.Remove(ele)
}

func (s *Stack) IsEmpty() bool {
	return s.Val.Len() <= 0
}

func (s *Stack) Top() interface{} {
	ele := s.Val.Back()
	if ele == nil {
		return nil
	}
	return ele.Value
}

func (s *Stack) StackLen() int {
	return s.Val.Len()
}
