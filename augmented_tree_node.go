package aabb

import "math"

type augTreeNode struct {
	Box Box
	Left *augTreeNode
	Right *augTreeNode
	MaxX int
	Height int
}

func newAugTreeNode(box Box) *augTreeNode {
	return &augTreeNode{
		Box: box,
		MaxX: box.XMax(),
	}
}

func (self *augTreeNode) GetHeight() int {
	if self == nil { return -1 }
	return self.Height
}

func (self *augTreeNode) GetMaxX() int {
	if self == nil { return math.MinInt }
	return self.MaxX
}

func (self *augTreeNode) GetBalance() int {
	return self.Right.GetHeight() - self.Left.GetHeight()
}

func (self *augTreeNode) RefreshHeight() {
	self.Height = max(self.Left.GetHeight(), self.Right.GetHeight()) + 1
}

func (self *augTreeNode) RefreshMaxX() {
	self.MaxX = max(self.Box.XMax(), self.Right.GetMaxX())
	self.MaxX = max(self.MaxX, self.Left.GetMaxX())
}

// Both RefreshHeight() and RefreshMaxX() at once.
func (self *augTreeNode) Refresh() {	
	self.MaxX   = self.Box.XMax()
	self.Height = 0
	if self.Left != nil {
		self.Height = self.Left.Height + 1
		self.MaxX   = max(self.MaxX, self.Left.MaxX)
	}
	if self.Right != nil {
		self.Height = self.Right.Height + 1
		self.MaxX   = max(self.MaxX, self.Right.MaxX)
	}
}

// ---- children addition ----

func (self *augTreeNode) NewLeftChild(box Box) {
	self.Left = newAugTreeNode(box)

	if self.Right == nil { self.Height += 1 }
	self.MaxX = max(box.XMax(), self.MaxX)
}

func (self *augTreeNode) NewRightChild(box Box) {
	self.Right = newAugTreeNode(box)

	if self.Left == nil { self.Height += 1 }
	self.MaxX = max(box.XMax(), self.MaxX)
}

// ---- rebalancing ----

// For better understanding of how rotations and rebalancing work,
// https://ksw2000.medium.com/implement-an-avl-tree-with-go-49e5952389d4
// is probably one of the best illustrated articles out there.

// Rebalance returns the new root node for the subtree that
// starts at this node.
func (self *augTreeNode) Rebalance() *augTreeNode {
	balance := self.GetBalance()
	if balance < -1 { // tree leaning left
		if self.Left.GetBalance() >= 0 {
			self.Left = self.Left.rotateLeft()
		}
		return self.rotateRight()
	} else if balance > 1 { // tree leaning right
		if self.Right.GetBalance() <= 0 {
			self.Right = self.Right.rotateRight()
		}	
		return self.rotateLeft()
	}
	return self
}

func (self *augTreeNode) rotateRight() *augTreeNode {
	originalLeftChild := self.Left
	self.Left = originalLeftChild.Right
	originalLeftChild.Right = self

	self.Refresh()
	originalLeftChild.Refresh()
	return originalLeftChild
}

func (self *augTreeNode) rotateLeft() *augTreeNode {
	originalRightChild := self.Right
	self.Right = originalRightChild.Left
	originalRightChild.Left = self

	self.Refresh()
	originalRightChild.Refresh()
	return originalRightChild
}
