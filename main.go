package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
)

type node struct {
	parent      *node
	left, right *node
	hash        string
}

func (n *node) calculate() {
	n.hash = hash([]byte(n.left.hash + n.right.hash))
	n.left.parent = n
	n.right.parent = n
}

func hash(in []byte) string {
	h := sha256.New()
	h.Write(in)
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	chunkSize := int64(64)
	buf := make([]byte, chunkSize)
	info, err := os.Stat("arrival_in_nara.txt")
	if err != nil {
		panic(err)
	}
	numChunks := math.Ceil(float64(info.Size()/chunkSize)) + 1
	chunks := make([][]byte, int64(numChunks))
	file, err := os.Open("arrival_in_nara.txt")
	nodes := make([][]node, 0)
	if err != nil {
		panic(err)
	}
	nodes = append(nodes, make([]node, 0))
	fmt.Println("Chunks:")
	i := 0
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		hash := hash(buf)

		fmt.Printf("#%02d %s\n", i, hash)
		chunks = append(chunks, buf)
		nodes[0] = append(nodes[0], node{hash: hash})
		i++
	}
	var rootNode *node
	rootNode = nil
	height := 0
	for rootNode == nil {
		if len(nodes[height]) < 2 {
			rootNode = &nodes[height][0]
		}
		nextHeight := make([]node, 0)
		for i := 0; i < int(len(nodes[height])/2)*2; i += 2 {
			newNode := node{left: &nodes[height][i], right: &nodes[height][i+1]}
			newNode.calculate()
			nextHeight = append(nextHeight, newNode)
		}
		if len(nodes[height])%2 != 0 {
			nextHeight = append(nextHeight, nodes[height][len(nodes[height])-1])
		}
		nodes = append(nodes, nextHeight)
		height++
	}
	fmt.Printf("\nRoot %s\n\n", rootNode.hash)
}
