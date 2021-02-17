package trie

import (
	"errors"
	"path"
	"strings"
)

var (
	ErrNotFound  = errors.New("ErrNotFound")
	ErrDuplicate = errors.New("ErrDuplicate")
)

type Root struct {
	node
}

func (r *Root) Put(s string, v interface{}) error {
	s = strings.TrimSpace(s)
	s = path.Clean(s)
	ss := strings.Split(s, "/")
	tmp := &r.node
	for _, each := range ss {
		next, ok := tmp.branch.get(each)
		if !ok {
			next = &node{path: each}
			tmp.branch = append(tmp.branch, next)
		}
		tmp = next
	}
	if tmp.leaf {
		return ErrDuplicate
	}
	tmp.leaf = true
	tmp.value = v
	return nil
}

func (r *Root) Get(s string) (interface{}, error) {
	s = strings.TrimSpace(s)
	s = path.Clean(s)
	ss := strings.Split(s, "/")
	tmp := &r.node
	for _, each := range ss {
		next, ok := tmp.branch.get(each)
		if !ok {
			next, ok = tmp.branch.get("*")
			if !ok {
				return nil, ErrNotFound
			}
		}
		tmp = next
	}
	if !tmp.leaf {
		return nil, ErrNotFound
	}
	return tmp.value, nil
}

type node struct {
	branch branch
	value  interface{}
	path   string
	leaf   bool
}

type branch []*node

func (b branch) get(s string) (*node, bool) {
	for _, each := range b {
		if each.path == s {
			return each, true
		}
	}
	return nil, false
}

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
