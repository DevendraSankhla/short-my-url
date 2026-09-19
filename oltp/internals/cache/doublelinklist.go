package cache

import (
	"errors"

	"github.com/devendrasankhla/short-my-url/oltp/internals/database"
)

type Node struct {
	Url  *database.URL
	Prev *Node
	Next *Node
}

type DouobleLinkedList struct {
	Head *Node
	Tail *Node
}

// PushFront
func (dq *DouobleLinkedList) PushFront(url *database.URL) error {
	if url == nil {
		return errors.New("nil parameter")
	}

	newNode := &Node{Url: url, Prev: nil, Next: nil}
	if dq.Head == nil {
		dq.Head = newNode
		dq.Tail = newNode
	} else {
		temp := dq.Head.Next
		dq.Head = newNode
		newNode.Next = temp
		temp.Prev = newNode
	}

	return nil
}

// PushBack
func (dq *DouobleLinkedList) PushBack(url *database.URL) error {
	if url == nil {
		return errors.New("nil parameter")
	}

	newNode := &Node{Url: url, Prev: nil, Next: nil}
	if dq.Head == nil {
		dq.PushFront(url)
	} else {
		dq.Tail.Next = newNode
		newNode.Prev = dq.Tail
		dq.Tail = newNode
	}

	return nil
}

// Front
func (dq DouobleLinkedList) Front() *database.URL {
	return dq.Head.Url
}

// Back
func (dq DouobleLinkedList) Back() *database.URL {
	return dq.Tail.Url
}

// Delete
func (dq *DouobleLinkedList) Delete(urlID int64) (*database.URL, error) {
	if dq.Head == nil {
		return nil, errors.New("List Empty")
	}
	if urlID == dq.Tail.Url.ID {
		nodeToDelete := dq.Tail
		dq.Tail = dq.Tail.Prev
		dq.Tail.Next = nil
		return nodeToDelete.Url, nil
	} else if urlID == dq.Head.Url.ID {
		nodeToDelete := dq.Head
		dq.Head = dq.Head.Next
		dq.Head.Prev = nil
		return nodeToDelete.Url, nil
	}

	temp := dq.Head
	for temp != nil {
		if temp.Url.ID == urlID {
			temp.Next.Prev = temp.Prev
			temp.Prev.Next = temp.Next
			return temp.Url, nil
		}
	}
	if temp == nil {
		return nil, errors.New("Item not found")
	}
	return nil, errors.New("Internal issue.")
}
