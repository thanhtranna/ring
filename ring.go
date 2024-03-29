package ring

import (
	"errors"
	"fmt"
	"sync"

	hash "github.com/OneOfOne/xxhash"
	"github.com/thanhtranna/rbt"
)

var (
	ErrEmptyRing   = errors.New("empty string")
	ErrKeyNotFound = errors.New("key not found")
)

type hasher interface {
	hash(string) int64
}

func newXXHash() hasher {
	return xxHash{}
}

// https://cyan4973.github.io/xxHash/
type xxHash struct{}

func (x xxHash) hash(data string) int64 {
	h := hash.New32()
	h.Write([]byte(data))
	r := h.Sum32()
	h.Reset()

	return int64(r)
}

type Ring struct {
	store        *rbt.Tree
	nodeMap      map[string]bool
	virtualNodes int
	hashfn       hasher

	mu sync.RWMutex
}

func New() *Ring {
	return &Ring{
		store:   rbt.NewTree(),
		nodeMap: make(map[string]bool),
		hashfn:  newXXHash(),
	}
}

func NewRing(nodes []string, virtualNodes int) *Ring {
	r := &Ring{
		store:        rbt.NewTree(),
		nodeMap:      make(map[string]bool),
		virtualNodes: virtualNodes,
		hashfn:       newXXHash(),
	}

	for _, node := range nodes {
		r.Add(node)
	}

	return r
}

func (r *Ring) hash(val string) int64 {
	return r.hashfn.hash(val)
}

func (r *Ring) Add(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.nodeMap[node]; ok {
		return
	}
	r.nodeMap[node] = true
	hashKey := r.hash(node)
	r.store.Insert(hashKey, node)

	for i := 0; i < r.virtualNodes; i++ {
		vNodeKey := fmt.Sprintf("%s-%d", node, i)
		r.nodeMap[vNodeKey] = true
		hashKey := r.hash(vNodeKey)
		r.store.Insert(hashKey, node)
	}
}

func (r *Ring) Remove(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.nodeMap[node]; !ok {
		return
	}

	hashKey := r.hash(node)
	r.store.Delete(hashKey)

	for i := 0; i < r.virtualNodes; i++ {
		vNodeKey := fmt.Sprintf("%s-%d", node, i)
		hashKey := r.hash(vNodeKey)
		r.store.Delete(hashKey)
		delete(r.nodeMap, vNodeKey)
	}

	delete(r.nodeMap, node)
}

func (r *Ring) Get(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.store.Size() == 0 {
		return "", ErrEmptyRing
	}

	var q *rbt.Node
	hashKey := r.hash(key)
	q = r.store.Nearest(hashKey)

	if hashKey > q.GetKey() {
		g := rbt.FindSuccessor(q)
		if g != nil {
			q = g
		} else {
			// If no successor found, return root(wrap around)
			q = r.store.Minimum()
		}
	}

	return q.GetValue(), nil
}
