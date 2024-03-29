package main

import (
	"fmt"
	"math/rand"

	"github.com/thanhtranna/ring"
)

func main() {
	nodes := []string{"a", "b", "c"}
	hashRing := ring.NewRing(nodes, 300)

	keyCount := 1000000
	distribution := make(map[string]int)
	key := make([]byte, 4)
	for i := 0; i < keyCount; i++ {
		rand.Read(key)
		node, err := hashRing.Get(string(key))
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}

		distribution[node]++
	}
	for node, count := range distribution {
		fmt.Printf("node: %s, key count: %d\n", node, count)
	}
}
