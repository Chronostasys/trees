package btree

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
)

type (
	Int int

	node struct {
		childs []*node
		father *node
		vals   []Hasher
		right  *node
		fn     int
		f      *os.File
		en     *gob.Encoder
		buf    *bytes.Buffer
	}

	Hasher interface {
		Hash() int
	}

	Tree struct {
		snmu         *sync.Mutex
		root         *node
		total        int
		m            int // 阶
		edge         int
		first        *node
		gfn          int
		fs           map[*node]struct{}
		f            *os.File
		buf          *bytes.Buffer
		en           *gob.Encoder
		takeSnapshot bool
		snapshot     map[int][]byte
		persist      bool
	}

	BinNode struct {
		Childs []int64
		Father int64
		Right  int64
		Vals   []Hasher
	}
)

func (i Int) Hash() int {
	return int(i)
}
