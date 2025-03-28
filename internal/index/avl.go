// In the index tree created the key is the hash(domain name)
// and the value is corresponding blockchain transaction

package index

import (
	"crypto/sha256"
	"fmt"
	"github.com/bleasey/bdns/internal/blockchain"
	"math"
	"strconv"
)

type AVLTree struct {
	root *AVLNode
}

func (t *AVLTree) Add(key string, value *blockchain.Transaction) {
	t.root = t.root.add(key, value)
}

func (t *AVLTree) Remove(key string) {
	t.root = t.root.remove(key)
}

func (t *AVLTree) Update(key string, newValue *blockchain.Transaction) {
	t.root = t.root.remove(key)
	t.root = t.root.add(key, newValue)
}

func (t *AVLTree) Search(key string) (node *AVLNode) {
	return t.root.search(key)
}

func (t *AVLTree) DisplayInOrder() {
	t.root.displayNodesInOrder()
}

// AVLNode structure
type AVLNode struct {
	key    string
	value  *blockchain.Transaction
	height int
	left   *AVLNode
	right  *AVLNode
}

// Adds a new node
func (n *AVLNode) add(key string, value *blockchain.Transaction) *AVLNode {
	if n == nil {
		return &AVLNode{key, value, 1, nil, nil}
	}

	if key < n.key {
		n.left = n.left.add(key, value)
	} else if key > n.key {
		n.right = n.right.add(key, value)
	} else {
		// if same key exists update value
		n.value = value
	}
	return n.rebalanceTree()
}

// Removes a node
func (n *AVLNode) remove(key string) *AVLNode {
	if n == nil {
		return nil
	}
	if key < n.key {
		n.left = n.left.remove(key)
	} else if key > n.key {
		n.right = n.right.remove(key)
	} else {
		if n.left != nil && n.right != nil {
			// node to delete found with both children;
			// replace values with smallest node of the right sub-tree
			rightMinNode := n.right.findSmallest()
			n.key = rightMinNode.key
			n.value = rightMinNode.value
			// delete smallest node that we replaced
			n.right = n.right.remove(rightMinNode.key)
		} else if n.left != nil {
			// node only has left child
			n = n.left
		} else if n.right != nil {
			// node only has right child
			n = n.right
		} else {
			// node has no children
			n = nil
			return n
		}
	}
	return n.rebalanceTree()
}

// Searches for a node
func (n *AVLNode) search(key string) *AVLNode {
	if n == nil {
		return nil
	}
	if key < n.key {
		return n.left.search(key)
	} else if key > n.key {
		return n.right.search(key)
	}
	return n
}

// Displays nodes left-depth first (used for debugging)
func (n *AVLNode) displayNodesInOrder() {
	if n.left != nil {
		n.left.displayNodesInOrder()
	}
	fmt.Print(n.key, " ")
	if n.right != nil {
		n.right.displayNodesInOrder()
	}
}

func (n *AVLNode) getHeight() int {
	if n == nil {
		return 0
	}
	return n.height
}

func (n *AVLNode) recalculateHeight() {
	n.height = 1 + int(math.Max(float64(n.left.getHeight()), float64(n.right.getHeight())))
}

// Checks if node is balanced and rebalance
func (n *AVLNode) rebalanceTree() *AVLNode {
	if n == nil {
		return n
	}
	n.recalculateHeight()

	// check balance factor and rotateLeft if right-heavy and rotateRight if left-heavy
	balanceFactor := n.left.getHeight() - n.right.getHeight()
	if balanceFactor == -2 {
		// check if child is left-heavy and rotateRight first
		if n.right.left.getHeight() > n.right.right.getHeight() {
			n.right = n.right.rotateRight()
		}
		return n.rotateLeft()
	} else if balanceFactor == 2 {
		// check if child is right-heavy and rotateLeft first
		if n.left.right.getHeight() > n.left.left.getHeight() {
			n.left = n.left.rotateLeft()
		}
		return n.rotateRight()
	}
	return n
}

// Rotate nodes left to balance node
func (n *AVLNode) rotateLeft() *AVLNode {
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (n *AVLNode) rotateRight() *AVLNode {
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (n *AVLNode) findSmallest() *AVLNode {
	if n.left != nil {
		return n.left.findSmallest()
	}
	return n
}

// Compute hash for AVL tree node
func ComputeIndexNodeHash(node *AVLNode) []byte {
	if node == nil {
		return nil
	}

	leftHash := ComputeIndexNodeHash(node.left)
	rightHash := ComputeIndexNodeHash(node.right)
	data := append([]byte(node.key+strconv.Itoa(node.value.TID)), leftHash...)
	data = append(data, rightHash...)

	hash := sha256.Sum256(data)

	return hash[:]
}
