package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Node represents a node in the Merkle tree
type Node struct {
	Hash  string
	Left  *Node
	Right *Node
}

// ComputeHash generates a SHA-256 hash for a given input
func ComputeHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// BuildMerkleTree constructs a full Merkle Tree from a list of file hashes
func BuildMerkleTree(hashes []string) *Node {
	if len(hashes) == 0 {
		return nil
	}

	if len(hashes) == 1 {
		return &Node{Hash: hashes[0]}
	}

	var nodes []*Node

	for i := 0; i < len(hashes); i += 2 {
		leftNode := &Node{Hash: hashes[i]}
		var rightNode *Node

		if i+1 < len(hashes) {
			rightNode = &Node{Hash: hashes[i+1]}
		} else {
			rightNode = leftNode
		}

		parentHash := ComputeHash(leftNode.Hash + rightNode.Hash)
		parentNode := &Node{
			Hash:  parentHash,
			Left:  leftNode,
			Right: rightNode,
		}

		nodes = append(nodes, parentNode)
	}

	return buildTreeFromNodes(nodes)
}

// Helper function to build a tree from parent nodes
func buildTreeFromNodes(nodes []*Node) *Node {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var parentNodes []*Node
	for i := 0; i < len(nodes); i += 2 {
		leftNode := nodes[i]
		var rightNode *Node

		if i+1 < len(nodes) {
			rightNode = nodes[i+1]
		} else {
			rightNode = leftNode
		}

		parentHash := ComputeHash(leftNode.Hash + rightNode.Hash)
		parentNodes = append(parentNodes, &Node{
			Hash:  parentHash,
			Left:  leftNode,
			Right: rightNode,
		})
	}

	return buildTreeFromNodes(parentNodes)
}

func PrintTree(node *Node, level int) {
	if node == nil {
		return
	}
	fmt.Printf("%s%s\n", fmt.Sprint(' '+level*2), node.Hash)
	PrintTree(node.Left, level+1)
	PrintTree(node.Right, level+1)
}

// GenerateProof generates a Merkle path for a given file hash
func GenerateProof(root *Node, fileHash string) ([][2]string, bool) {
	var path [][2]string // Each entry is (siblingHash, "L" or "R")
	found := findProof(root, fileHash, &path)
	return path, found
}

// Helper function to find the proof path in the Merkle tree
func findProof(node *Node, fileHash string, path *[][2]string) bool {
	if node == nil {
		return false
	}

	if node.Left == nil && node.Right == nil {
		return node.Hash == fileHash
	}

	// Search in left subtree
	if node.Left != nil && findProof(node.Left, fileHash, path) {
		if node.Right != nil {
			*path = append(*path, [2]string{node.Right.Hash, "R"}) // Right sibling
		}
		return true
	}

	// Search in right subtree
	if node.Right != nil && findProof(node.Right, fileHash, path) {
		if node.Left != nil {
			*path = append(*path, [2]string{node.Left.Hash, "L"}) // Left sibling
		}
		return true
	}

	return false
}

// VerifyProof checks if a given file hash is part of the Merkle tree
func VerifyProof(fileHash string, proof [][2]string, rootHash string) bool {
	computedHash := fileHash

	for _, entry := range proof {
		siblingHash, position := entry[0], entry[1]

		if position == "L" {
			computedHash = ComputeHash(siblingHash + computedHash) 
		} else {
			computedHash = ComputeHash(computedHash + siblingHash) 
		}
	}

	return computedHash == rootHash
}
