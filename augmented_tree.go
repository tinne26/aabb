package aabb

import "strings"
import "strconv"

// TODO: implement "Remove" and full Space interface

// An augmented tree only sorts by x, not y. This makes it decent for
// side scrollers and similar platformers that have wide levels that
// aren't insanely vertical. Among the positive properties, the boxes
// are not duplicated and the tree can be modified on the fly (deleting
// and adding boxes are O(log n) operations). Among the negatives, it's
// less general than other methods and rather slow. A transposed
// version of this would also work well for games with very vertical
// levels (e.g. Downwell).
type AugmentedTree struct {
	root *augTreeNode
}

func NewAugmentedTree() *AugmentedTree {
	return &AugmentedTree{}
}

func (self *AugmentedTree) Clear() {
	self.root = nil
}

func (self *AugmentedTree) Add(box Box) {
	if self.root == nil {
		self.root = newAugTreeNode(box)
	} else {
		self.root = augmentedTreeRecursiveAdd(self.root, box)
	}
}

// Returns the new top node, in case it's modified due to tree
// rebalancing in the process of addition.
func augmentedTreeRecursiveAdd(node *augTreeNode, box Box) *augTreeNode {
	// see if we have to continue left or right
	if box.XMin() <= node.Box.XMin() {
		if node.Left == nil { // add to leaf
			node.NewLeftChild(box)
		} else { // add recursively
			node.Left = augmentedTreeRecursiveAdd(node.Left, box)
			node.Refresh()
		}
	} else { // symmetrical case
		if node.Right == nil { // add to leaf
			node.NewRightChild(box)
		} else { // add recursively
			node.Right = augmentedTreeRecursiveAdd(node.Right, box)
			node.Refresh()
		}
	}

	return node.Rebalance()
}

// Returns a colliding Box in the Space, or nil if there are no
// collisions. Collisions where `givenBox == otherBox` are not
// reported. If multiple collisions are possible, any of them
// might be returned. See EachCollision if you need more control.
func (self *AugmentedTree) Collision(box Box) Box {
	return self.recursiveCollision(self.root, box)
}

func (self *AugmentedTree) recursiveCollision(node *augTreeNode, box Box) Box {
	if box.XMin() > node.GetMaxX() { return nil } // no possible collision in this sub-branch
	if BoxesCollide(box, node.Box) && box != node.Box { return node.Box }

	// check on the left, then on the right
	if node.Left != nil {
		collidingBox := self.recursiveCollision(node.Left, box)
		if collidingBox != nil { return collidingBox }
	}
	if node.Right != nil && box.XMax() >= node.Box.XMin() {
		collidingBox := self.recursiveCollision(node.Right, box)
		if collidingBox != nil { return collidingBox }
	}
	return nil
}

func (self *AugmentedTree) MutateBox(box MutableBox, xMin, xMax, yMin, yMax int) {
	// Note: this could be done more optimally in many cases, but it's a pain, honestly.
	if !self.Remove(box) { panic("box to be updated not found") }
	box.MutateBox(xMin, xMax, yMin, yMax)
	self.Add(box)
}

func (self *AugmentedTree) Remove(box Box) bool {
	var removed bool
	self.root, removed = self.recursiveRemove(self.root, box)
	return removed
}

func (self *AugmentedTree) recursiveRemove(node *augTreeNode, box Box) (*augTreeNode, bool) {
	if box.XMin() > node.GetMaxX() { return node, false } // no possible collision in this sub-branch
	if box == node.Box {
		// easy cases: node was leaf or had only one children
		if node.Left == nil {
			if node.Right == nil {
				return nil, true
			} else {
				return node.Right, true
			}
		} else if node.Right == nil {
			return node.Left, true
		}

		// hard case: node has two children. to solve this, we need to find
		// the first child that's bigger than the current node, so we can use
		// it to replace the node being deleted. this descendant will be the
		// leftmost node of the right subtree of node, aka the inorder successor.
		// this can also be done symmetrically with the largest child of the left
		// subtree. get some paper and draw it, hard to understand otherwise.
		if node.Right.Left == nil { // special case
			node.Box   = node.Right.Box // replace content
			node.Right = node.Right.Right // remove right (min of node's subtree)
		} else { // general case
			var inorderNode *augTreeNode
			node.Right, inorderNode = self.extractInorderNode(node.Right, node.Right.Left)
			node.Box = inorderNode.Box
		}
		node.Refresh()
		return node.Rebalance(), true
	}

	// check on the left, then on the right
	var removed bool
	if node.Left != nil {
		node.Left, removed = self.recursiveRemove(node.Left, box)
		if removed {
			node.Refresh()
			return node.Rebalance(), true
		}
	}
	if node.Right != nil && box.XMax() >= node.Box.XMin() {
		node.Right, removed = self.recursiveRemove(node.Right, box)
		if removed {
			node.Refresh()
			return node.Rebalance(), true
		}
	}
	return node, false
}

func (self *AugmentedTree) extractInorderNode(parent, node *augTreeNode) (*augTreeNode, *augTreeNode) {
	if node.Left == nil { // base case
		parent.Left = node.Right
	} else { // recursive case
		parent.Left, node = self.extractInorderNode(node, node.Left)
	}
	parent.Refresh()
	return parent.Rebalance(), node
}

func (self *AugmentedTree) Stabilize() {
	if self.root == nil { return }
	left, right := self.root.Left, self.root.Right
	self.root.Left, self.root.Right = nil, nil
	self.root.Height = 0
	self.root.RefreshMaxX()
	if left != nil {
		self.root = augmentedTreeRecursiveDfsAdd(self.root, left)
	}
	if right != nil {
		self.root = augmentedTreeRecursiveDfsAdd(self.root, right)
	}
}

func augmentedTreeRecursiveDfsAdd(root, node *augTreeNode) *augTreeNode {
	root = augmentedTreeRecursiveAdd(root, node.Box)
	left, right := node.Left, node.Right
	if left != nil {
		root = augmentedTreeRecursiveDfsAdd(root, left)
	}
	if right != nil {
		root = augmentedTreeRecursiveDfsAdd(root, right)
	}
	return root
}

// A more complete version of Collision() that calls the given
// function for each collision (instead of stopping at one).
func (self *AugmentedTree) EachCollision(box Box, fn func(Box) SearchControl) {
	_ = self.recursiveEachCollision(self.root, box, fn)
}

func (self *AugmentedTree) recursiveEachCollision(node *augTreeNode, box Box, fn func(Box) SearchControl) SearchControl {
	if box.XMin() > node.GetMaxX() { return SearchContinue } // no possible collision in this sub-branch
	if BoxesCollide(box, node.Box) && box != node.Box {
		if fn(node.Box) == SearchStop { return SearchStop }
	}

	// check recursively on left and right branches
	if node.Left != nil {
		control := self.recursiveEachCollision(node.Left, box, fn)
		if control == SearchStop { return SearchStop }
	}
	if node.Right != nil && box.XMax() >= node.Box.XMin() {
		control := self.recursiveEachCollision(node.Right, box, fn)
		if control == SearchStop { return SearchStop }
	}
	return SearchContinue
}

func (self *AugmentedTree) EachInXRange(minX, maxX int, fn func(Box) SearchControl) {
	panic("EachInXRange() unimplemented")
}

// --- debug and testing methods ---

func (self *AugmentedTree) dfsHeight() int {
	return self.recDfsHeight(self.root)
}

func (self *AugmentedTree) recDfsHeight(node *augTreeNode) int {
	if node == nil { return 0 }
	return 1 + max(
		self.recDfsHeight(node.Left),
		self.recDfsHeight(node.Right),
	)
}

// Mostly for debug.
func (self *AugmentedTree) GetRootBox() Box {
	return self.root.Box
}

// Mostly for debug.
func (self *AugmentedTree) StringRepresentation() string {
	var strBuilder strings.Builder
	self.recStringRepresentation(&strBuilder, self.root)
	return strBuilder.String()
}

func (self *AugmentedTree) recStringRepresentation(strBuilder *strings.Builder, node *augTreeNode) {
	if node == nil { strBuilder.Write([]byte{'(', ')'}) ; return }
	strBuilder.Write([]byte{'(', 'H'})
	strBuilder.WriteString(strconv.Itoa(node.Height))
	strBuilder.Write([]byte{' ', '|', ' ', 'M', 'a', 'x', 'X'})
	strBuilder.WriteString(strconv.Itoa(node.MaxX))
	strBuilder.Write([]byte{' '})
	strBuilder.WriteString(BoxString(node.Box))
	strBuilder.Write([]byte{' ', '-', '>', ' ', 'L'})
	self.recStringRepresentation(strBuilder, node.Left)
	strBuilder.Write([]byte{',', ' ', 'R'})
	self.recStringRepresentation(strBuilder, node.Right)
	strBuilder.Write([]byte{')'})
}
