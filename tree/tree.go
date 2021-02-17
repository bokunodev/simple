package tree

import (
	"errors"
	"path"
	"strings"
	"sync"
)

var (
	ErrNotFound  = errors.New("ErrNotFound")
	ErrDuplicate = errors.New("ErrDuplicate")
)

type Root struct {
	node
	rwmu sync.RWMutex
}

func (r *Root) Put(s string, v interface{}) error {
	r.rwmu.Lock()
	defer r.rwmu.Unlock()

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
	r.rwmu.RLock()
	defer r.rwmu.RUnlock()

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
