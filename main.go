package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type node struct {
	parent      *node
	left, right *node
	hash        string
}

type proofNode struct {
	left bool
	hash string
}
type merkleProof []proofNode

func (n node) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Left  *node `json:"left"`
		Right *node `json:"right"`
		//Parent *node `json:"parent"`
		Hash string
	}{
		n.left,
		n.right,
		//n.parent,
		n.hash,
	})
}

func (n *node) calculate() {
	n.hash = hash([]byte(n.left.hash + n.right.hash))
	n.left.parent = n
	n.right.parent = n
}

func (n *node) Sibling() *node {
	if n.parent == nil {
		return nil
	}
	if n.parent.left.hash == n.hash {
		return n.parent.right
	}
	return n.parent.left
}

func (n *node) Uncle() *node {
	if n.parent == nil {
		return nil
	}
	return n.parent.Sibling()
}

func hash(in []byte) string {
	h := sha256.New()
	h.Write(in)
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	chunkSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, chunkSize)
	info, err := os.Stat(os.Args[1])
	if err != nil {
		panic(err)
	}
	numChunks := math.Ceil(float64(info.Size()/int64(chunkSize))) + 1
	chunks := make([][]byte, int64(numChunks))
	file, err := os.Open(os.Args[1])
	nodes := make([][]*node, 0)
	if err != nil {
		panic(err)
	}
	nodes = append(nodes, make([]*node, 0))
	fmt.Println("Chunks:")
	i := 0
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		hash := hash(buf)

		fmt.Printf("#%03d %s\n", i, hash)
		chunks = append(chunks, buf)
		nodes[0] = append(nodes[0], &node{hash: hash})
		i++
	}
	var rootNode *node
	rootNode = nil
	height := 0
	for rootNode == nil {
		if len(nodes[height]) < 2 {
			rootNode = nodes[height][0]
		}
		nextHeight := make([]*node, 0)
		for i := 0; i < int(len(nodes[height])/2)*2; i += 2 {
			newNode := &node{left: nodes[height][i], right: nodes[height][i+1]}
			newNode.calculate()
			nextHeight = append(nextHeight, newNode)
		}
		if len(nodes[height])%2 != 0 {
			nextHeight = append(nextHeight, nodes[height][len(nodes[height])-1])
		}
		nodes = append(nodes, nextHeight)
		height++
	}
	fmt.Printf("Root %s\n\n", rootNode.hash)

	idx := 0
	fmt.Printf("node %s\n", nodes[0][idx].hash)

	proveNode := nodes[0][idx]
	proof := getProof(proveNode)
	fmt.Printf("Proof: %v\n", proof)

	fmt.Printf("Good? %v\n", proof.Prove(proveNode.hash))

	// jsonStr, err := json.Marshal(rootNode)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(jsonStr))
}

// getProof returns a list of hashes and whether to place your computed hash on the left
// The first item in the hash is for the sibling at height 0, then for the sibling of the computed hash
// The last item in the hash is the root hash and should be compared against the computed root hash.
func getProof(n *node) merkleProof {
	proof := merkleProof{}
	nextProof := n.Sibling()
	for {
		left := nextProof.parent.right.hash == nextProof.hash
		proof = append(proof, proofNode{left, nextProof.hash})
		if nextProof.Uncle() == nil {
			proof = append(proof, proofNode{left, nextProof.parent.hash})
			break
		}
		nextProof = nextProof.Uncle()
	}
	return proof
}

// Prove proves that the a hash is correct for the given proof
func (p merkleProof) Prove(ha string) bool {
	rootHash := p[len(p)-1].hash
	for i := 0; i < len(p)-1; i++ {
		if p[i].left {
			ha = hash([]byte(ha + p[i].hash))
		} else {
			ha = hash([]byte(p[i].hash + ha))
		}
	}

	return ha == rootHash
}
