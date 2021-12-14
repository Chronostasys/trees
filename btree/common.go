package btree

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
)

type (
	Int int
	KV  struct {
		K, V string
	}

	node struct {
		childs []*node
		father *node
		vals   []Item
		right  *node
		fn     int
		f      *os.File
		en     *gob.Encoder
		buf    *bytes.Buffer
	}

	Item interface {
		Less(than Item) bool
		Key() Item
		EQ(b Item) bool
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
		Vals   []Item
	}
)

func (i Int) Less(than Item) bool {
	return i < than.(Int)
}
func (i Int) EQ(than Item) bool {
	return i == than.(Int)
}
func (i Int) Key() Item {
	return i
}
func (i Int) Int() int {
	return int(i)
}
func (i KV) Less(than Item) bool {
	return i.K < than.(KV).K
}
func (i KV) EQ(than Item) bool {
	return i.K == than.(KV).K
}
func (i KV) Key() Item {
	return KV{K: i.K}
}
