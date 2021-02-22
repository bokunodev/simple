package trie

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("ErrNotFound")
	ErrDuplicate = errors.New("ErrDuplicate")
)

type Root struct{ *node }

func New() Root { return Root{&node{branch: make(branch, 0)}} }

func (r *Root) Put(s string, v interface{}) error {
	if r.node.branch == nil {
		r.node.branch = branch{&node{path: s}}
		r.node.leaf = true
		r.node.value = v
		return nil
	}
	return r.node.put(s, v)
}

type node struct {
	branch branch
	value  interface{}
	path   string
	leaf   bool
}

func (n *node) put(s string, v interface{}) error {
	return nil
}

type branch []*node

func compare(s1, s2 string) int {
	if len(s1) > len(s2) {
		tmp := s1
		s1 = s2
		s2 = tmp
	}
	ls := len(s1)
	for i := 0; i < ls; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}
	return ls
}
