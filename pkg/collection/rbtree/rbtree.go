package rbtree

import (
	"github.com/jpfourny/papaya/pkg/optional"
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/stream"
)

// RedBlackTree represents a Red-Black Tree with a map-like interface.
type RedBlackTree[K any, V any] struct {
	compare func(K, K) int
	root    *node[K, V] // Root node of the tree.
	size    int         // Number of nodes in the tree.
}

// color is the color of a node in a Red-Black Tree.
// A node is either red (true) or black (false).
// red is true, black is false.
type color bool

// node is a node in a Red-Black Tree.
// A node has a key, a value, a color, and up to two children.
// The children are either both present or both nil.
// The children are ordered by their keys.
type node[K any, V any] struct {
	key    K
	value  V
	left   *node[K, V]
	right  *node[K, V]
	parent *node[K, V]
	color  color
}

// New returns a new empty Red-Black Tree.
func New[K comparable, V any](compare func(K, K) int) *RedBlackTree[K, V] {
	return &RedBlackTree[K, V]{
		compare: compare,
	} // Empty tree.
}

// Size returns the number of nodes in the tree.
func (t *RedBlackTree[K, V]) Size() int {
	return t.size
}

// IsEmpty returns true if the tree is empty; false otherwise.
func (t *RedBlackTree[K, V]) IsEmpty() bool {
	return t.size == 0
}

// Clear removes all nodes from the tree.
func (t *RedBlackTree[K, V]) Clear() {
	t.root = nil
	t.size = 0
}

// Contains returns true if the tree contains a node with the given key; false otherwise.
func (t *RedBlackTree[K, V]) Contains(key K) (ok bool) {
	_, ok = t.get(key)
	return
}

// Get returns the value associated with the given key, or nil if the key is not present in the tree.
func (t *RedBlackTree[K, V]) Get(key K) optional.Optional[V] {
	if n, ok := t.get(key); ok {
		return optional.Of(n.value)
	}
	return optional.Empty[V]()
}

// Put inserts the given key-value pair into the tree.
// If the key is already present in the tree, then the value is updated.
// Returns true if the key was not already present in the tree; false otherwise.
func (t *RedBlackTree[K, V]) Put(key K, value V) (ok bool) {
	var n *node[K, V]
	n, ok = t.put(key, value)
	if ok {
		t.fixAfterInsert(n)
	}
	return
}

// Delete removes the node with the given key from the tree.
// Returns true if the key was present in the tree; false otherwise.
func (t *RedBlackTree[K, V]) Delete(key K) (ok bool) {
	var n *node[K, V]
	n, ok = t.get(key)
	if ok {
		t.delete(n)
	}
	return
}

// Keys returns a slice containing all keys in the tree.
// The keys are in ascending order.
func (t *RedBlackTree[K, V]) Keys() []K {
	keys := make([]K, 0, t.size)
	t.forEach(func(n *node[K, V]) bool {
		keys = append(keys, n.key)
		return true
	})
	return keys
}

// Values returns a slice containing all values in the tree.
// The values are in ascending order of their corresponding keys.
func (t *RedBlackTree[K, V]) Values() []V {
	values := make([]V, 0, t.size)
	t.forEach(func(n *node[K, V]) bool {
		values = append(values, n.value)
		return true
	})
	return values
}

// ForEach calls the given function for each node in the tree.
// The nodes are visited in ascending order of their keys.
// If the given function returns false, then the iteration is stopped.
func (t *RedBlackTree[K, V]) ForEach(f func(K, V) bool) {
	t.forEach(func(n *node[K, V]) bool {
		return f(n.key, n.value)
	})
}

// Stream returns a stream containing all key-value pairs in the tree.
// The pairs are in ascending order of their keys.
func (t *RedBlackTree[K, V]) Stream() stream.Stream[pair.Pair[K, V]] {
	return func(yield stream.Consumer[pair.Pair[K, V]]) bool {
		return t.forEach(func(n *node[K, V]) bool {
			return yield(pair.Of(n.key, n.value))
		})
	}
}

// forEach calls the given function for each node in the tree.
// The nodes are visited in ascending order of their keys.
// If the given function returns false, then the iteration is stopped.
// Returns false if the iteration was stopped before all nodes were visited; true otherwise.
func (t *RedBlackTree[K, V]) forEach(f func(*node[K, V]) bool) bool {
	n := t.minimum(t.root)
	for n != nil {
		if !f(n) {
			return false // Consumer saw enough.
		}
		n = t.successor(n)
	}
	return true
}

// get returns the node associated with the given key, or nil if the key is not present in the tree.
func (t *RedBlackTree[K, V]) get(key K) (*node[K, V], bool) {
	n := t.root
	for n != nil {
		c := t.compare(key, n.key)
		if c < 0 {
			n = n.left
		} else if c > 0 {
			n = n.right
		} else {
			return n, true
		}
	}
	return nil, false
}

// put inserts the given key-value pair into the tree.
// If the key is already present in the tree, then the value is updated.
// Returns the node that was inserted or updated and true if the key was not already present in the tree; false otherwise.
func (t *RedBlackTree[K, V]) put(key K, value V) (*node[K, V], bool) {
	var n *node[K, V]

	if t.root == nil {
		n = &node[K, V]{
			key:   key,
			value: value,
			color: false,
		}
		t.root = n
		t.size++
		return n, true
	}

	n = t.root
	for {
		c := t.compare(key, n.key)
		if c < 0 {
			if n.left == nil {
				n.left = &node[K, V]{
					key:    key,
					value:  value,
					parent: n,
					color:  true,
				}
				t.size++
				return n.left, true
			}
			n = n.left
		} else if c > 0 {
			if n.right == nil {
				n.right = &node[K, V]{
					key:    key,
					value:  value,
					parent: n,
					color:  true,
				}
				t.size++
				return n.right, true
			}
			n = n.right
		} else {
			n.value = value
			return n, false
		}
	}
}

// fixAfterInsert restores the Red-Black Tree invariants after inserting a node.
// The given node must be non-nil and have a red parent.
func (t *RedBlackTree[K, V]) fixAfterInsert(n *node[K, V]) {
	n.color = true
	for n != nil && n != t.root && n.parent.color {
		if n.parent == n.parent.parent.left {
			y := n.parent.parent.right
			if y != nil && y.color {
				n.parent.color = false
				y.color = false
				n.parent.parent.color = true
				n = n.parent.parent
			} else {
				if n == n.parent.right {
					n = n.parent
					t.rotateLeft(n)
				}
				n.parent.color = false
				n.parent.parent.color = true
				t.rotateRight(n.parent.parent)
			}
		} else {
			y := n.parent.parent.left
			if y != nil && y.color {
				n.parent.color = false
				y.color = false
				n.parent.parent.color = true
				n = n.parent.parent
			} else {
				if n == n.parent.left {
					n = n.parent
					t.rotateRight(n)
				}
				n.parent.color = false
				n.parent.parent.color = true
				t.rotateLeft(n.parent.parent)
			}
		}
	}
	t.root.color = false
}

// rotateLeft performs a left rotation on the given node.
// The given node must have a non-nil right child and not be the root of the tree.
// The right child becomes the new parent of the given node.
// The given node becomes the left child of the new parent.
// The left child of the new parent becomes the right child of the given node.
// The new parent is returned.
func (t *RedBlackTree[K, V]) rotateLeft(n *node[K, V]) {
	r := n.right
	n.right = r.left
	if r.left != nil {
		r.left.parent = n
	}
	r.parent = n.parent
	if n.parent == nil {
		t.root = r
	} else if n == n.parent.left {
		n.parent.left = r
	} else {
		n.parent.right = r
	}
	r.left = n
	n.parent = r
}

// rotateRight performs a right rotation on the given node.
// The given node must have a non-nil left child and not be the root of the tree.
// The left child becomes the new parent of the given node.
// The given node becomes the right child of the new parent.
// The right child of the new parent becomes the left child of the given node.
// The new parent is returned.
func (t *RedBlackTree[K, V]) rotateRight(n *node[K, V]) {
	l := n.left
	n.left = l.right
	if l.right != nil {
		l.right.parent = n
	}
	l.parent = n.parent
	if n.parent == nil {
		t.root = l
	} else if n == n.parent.right {
		n.parent.right = l
	} else {
		n.parent.left = l
	}
	l.right = n
	n.parent = l
}

// delete removes the given node from the tree.
// The given node must be non-nil and present in the tree.
func (t *RedBlackTree[K, V]) delete(n *node[K, V]) {
	t.size--
	if n.left != nil && n.right != nil {
		// n has two children.
		// Replace n with its successor.
		s := t.successor(n)
		n.key = s.key
		n.value = s.value
		n = s
	}
	// n has at most one child.
	// Replace n with its child.
	var child *node[K, V]
	if n.left != nil {
		child = n.left
	} else {
		child = n.right
	}

	if child != nil {
		// n has one child.
		// Replace n with its child.
		child.parent = n.parent
		if n.parent == nil {
			t.root = child
		} else if n == n.parent.left {
			n.parent.left = child
		} else {
			n.parent.right = child
		}
		if !n.color {
			t.fixAfterDelete(child)
		}
	}

	if n.parent == nil {
		// n was the root node.
		t.root = nil
	} else {
		// n was not the root node.
		if !n.color {
			t.fixAfterDelete(n)
		}
		if n.parent.left == n {
			n.parent.left = nil
		} else {
			n.parent.right = nil
		}
		n.parent = nil
	}
}

// fixAfterDelete restores the Red-Black Tree invariants after deleting a node.
func (t *RedBlackTree[K, V]) fixAfterDelete(child *node[K, V]) {
	for child != t.root && !child.color {
		if child == child.parent.left {
			s := child.parent.right
			if s.color {
				s.color = false
				child.parent.color = true
				t.rotateLeft(child.parent)
				s = child.parent.right
			}
			if !s.left.color && !s.right.color {
				s.color = true
				child = child.parent
			} else {
				if !s.right.color {
					s.left.color = false
					s.color = true
					t.rotateRight(s)
					s = child.parent.right
				}
				s.color = child.parent.color
				child.parent.color = false
				s.right.color = false
				t.rotateLeft(child.parent)
				child = t.root
			}
		} else {
			s := child.parent.left
			if s.color {
				s.color = false
				child.parent.color = true
				t.rotateRight(child.parent)
				s = child.parent.left
			}
			if !s.right.color && !s.left.color {
				s.color = true
				child = child.parent
			} else {
				if !s.left.color {
					s.right.color = false
					s.color = true
					t.rotateLeft(s)
					s = child.parent.left
				}
				s.color = child.parent.color
				child.parent.color = false
				s.left.color = false
				t.rotateRight(child.parent)
				child = t.root
			}
		}
	}
	child.color = false
}

// successor returns the successor of the given node.
// The given node must be non-nil.
func (t *RedBlackTree[K, V]) successor(n *node[K, V]) *node[K, V] {
	if n.right != nil {
		return t.minimum(n.right)
	}
	p := n.parent
	for p != nil && n == p.right {
		n = p
		p = p.parent
	}
	return p
}

// minimum returns the minimum node in the subtree rooted at the given node.
func (t *RedBlackTree[K, V]) minimum(n *node[K, V]) *node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}
