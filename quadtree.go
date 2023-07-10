package aabb

import "strings"
import "strconv"

// [Quadtree] implements a quadtree to check collisions between [Box]es.
// Optimizations take inspiration from stackoverflow.com/questions/41946007.
//
// Before considering quatrees you should fully understand regular [Grid]s.
// Quatrees, when fully expanded, are basically *slow* grids. They share all
// their same key defects, and in the pathological cases their performance
// will degrade even faster as they are less cache-friendly structures. The
// only advantage they have is when significant portions of the [Space] are
// unused. In those cases, a Quadtree might be created with a higher resolution
// than a Grid would allow, while still consuming less memory. That being said,
// this kind of scenario is not particularly common in 2D games, though, so don't
// rush to blindly use quadtrees just because the name sounds cool and popular.
type Quadtree struct {
	workingArea Box
	maxDepth    int // max number of cells = 4^maxDepth
	splitCutoff int // max number of boxes in leaf before split
	                // (unless we are at maxDepth and can't split any more)

	// main data
	cellIters cellIterList
	markedBoxes markedBoxList
	nodes []int // node values are indices to be interpreted as follows:
	            // - int > 0: index to next node group
	            // - int == 0: uninitialized leaf node
					// - int < 0: leaf node, use ^int to address cellIters

	// helper fields for internal operation
	nodesFreeIndex int

	// helper fields for certain operations that require
	// additional annotations or details
	helperIndex int // used on multiple operations, like on EachCollision()
	                // to detect if we have to abort, but also on Remove()
						 // to track the first removal index
}

// Max amount of cells = 4^maxDepth
const StdQuadtreeDepth       = 4 // common values between 4 and 8
const StdQuadtreeSplitCutoff = 8 // common values between 6 and 12
func NewQuadtree(workingArea Box, maxDepth int, splitCutoff int) *Quadtree {
	// assert parameters validity
	BoxPanicIfInvalid(workingArea)
	if splitCutoff < 1 { panic("splitCutoff < 1") }
	if maxDepth < 1 { panic("maxDepth < 1") } // 0 could be used too, but it's
	                                          // kind of a silly self-imposed
															// edge case, so... no, thanks
	if (1 << maxDepth) > BoxWidth(workingArea) {
		panic("2^maxDepth > BoxWidth(workingArea) [too many horizontal cells]")
	}
	if (1 << maxDepth) > BoxHeight(workingArea) {
		panic("2^maxDepth > BoxHeight(workingArea) [too many vertical cells]")
	}

	// create and return the quadtree
	startingNodes := make([]int, 4, 128)
	for i := 0; i < 4; i++ { startingNodes[i] = 0 } // uninitialized leaves/cells
   return &Quadtree {
		workingArea: workingArea,
		maxDepth: maxDepth,
		splitCutoff: splitCutoff,

		markedBoxes: newMarkedBoxList(32),
		cellIters: newCellIterList(64),
		nodes: startingNodes,
		nodesFreeIndex: -1,
	}
}

func (self *Quadtree) Clear() {
	self.markedBoxes.Clear()
	self.cellIters.Clear()
	self.nodesFreeIndex = -1
	self.nodes = self.nodes[0 : 4]
	for i := 0; i < 4; i++ { self.nodes[i] = 0 } // uninitialized leaves/cells
}

func (self *Quadtree) Add(box Box) {
	// we place cellIters in all nodes where boxes collide.
	xMin, xMax := self.workingArea.XMin(), self.workingArea.XMax()
	yMin, yMax := self.workingArea.YMin(), self.workingArea.YMax()
	spaceBoxChecksWith(box, xMin, xMax, yMin, yMax)

	boxIndex := self.markedBoxes.AddBox(box)
	nextQuadIdx := 0 // we use 0 because the root is omitted
	quadDepth   := 1 // 0 would be root
	self.quadAdd(box, boxIndex, nextQuadIdx, quadDepth, xMin, xMax, yMin, yMax)
}

// Try to add a Box to a given group of 4 nodes (and their potential children).
func (self *Quadtree) quadAdd(box Box, boxIndex, quadIdx, quadDepth, xMin, xMax, yMin, yMax int) {
	// compute node group space partition on the fly
	centerInnerX := xMin + ((xMax - xMin) >> 1)
	centerInnerY := yMin + ((yMax - yMin) >> 1)

	// try to add the given box on each node at the current level
	xMinChild, xMaxChild, yMinChild, yMaxChild := xMin, centerInnerX, yMin, centerInnerY
	self.nodeAdd(box, boxIndex, quadIdx + 0, quadDepth, xMinChild, xMaxChild, yMinChild, yMaxChild)
	xMinChild, xMaxChild = centerInnerX + 1, xMax
	self.nodeAdd(box, boxIndex, quadIdx + 1, quadDepth, xMinChild, xMaxChild, yMinChild, yMaxChild)
	yMinChild, yMaxChild = centerInnerY + 1, yMax
	self.nodeAdd(box, boxIndex, quadIdx + 2, quadDepth, xMinChild, xMaxChild, yMinChild, yMaxChild)
	xMinChild, xMaxChild = xMin, centerInnerX
	self.nodeAdd(box, boxIndex, quadIdx + 3, quadDepth, xMinChild, xMaxChild, yMinChild, yMaxChild)
}

// Try to add a Box to a specific node (and their potential children).
func (self *Quadtree) nodeAdd(box Box, boxIndex, nodeIdx, nodeDepth, xMin, xMax, yMin, yMax int) {
	// prune descent
	if !BoxCollidesWith(box, xMin, xMax, yMin, yMax) { return }

	nodeValue := self.nodes[nodeIdx]
	if nodeValue > 0 {
		// interpret nodeValue as index to next quad of nodes
		self.quadAdd(box, boxIndex, nodeValue, nodeDepth + 1, xMin, xMax, yMin, yMax)
	} else if nodeValue == 0 {
		// uninitialized leaf/cell node
		self.nodes[nodeIdx] = ^self.cellIters.AddIterTo(boxIndex, -1)
	} else {
		// negative value, leaf/cell reached
		if nodeDepth >= self.maxDepth || self.cellHasRoom(^nodeValue) {
			self.nodes[nodeIdx] = ^self.cellIters.AddIterTo(boxIndex, ^nodeValue)
		} else {
			// split node into quad
			quadIndex := self.registerQuad()
			self.nodes[nodeIdx] = quadIndex

			// add box to the new quad
			self.quadAdd(box, boxIndex, quadIndex, nodeDepth + 1, xMin, xMax, yMin, yMax)

			// reflow current leaf elements into the new quad
			cellIterIndex := ^nodeValue
			for cellIterIndex != -1 {
				boxIndex, nextIndex := self.cellIters.Next(cellIterIndex)
				self.cellIters.CutIter(cellIterIndex, -1)
				reflownBox := self.markedBoxes.GetBoxAt(boxIndex)
				self.quadAdd(reflownBox, boxIndex, quadIndex, nodeDepth + 1, xMin, xMax, yMin, yMax)
				cellIterIndex = nextIndex
			}
		}
	}
}

func (self *Quadtree) MutateBox(box MutableBox, xMin, xMax, yMin, yMax int) {
	matchedBox := self.removeInternal(box)
	if matchedBox == nil {
		panic("can't update box, it can't be found in the quadtree")
	}
	matchedBox.(MutableBox).MutateBox(xMin, xMax, yMin, yMax)
	spaceBoxChecks(matchedBox, self.workingArea)

	xMinArea, xMaxArea := self.workingArea.XMin(), self.workingArea.XMax()
	yMinArea, yMaxArea := self.workingArea.YMin(), self.workingArea.YMax()
	newBox := NewBox(xMin, xMax, yMin, yMax)
	spaceBoxChecks(newBox, self.workingArea)
	self.quadAdd(matchedBox, self.helperIndex, 0, 1, xMinArea, xMaxArea, yMinArea, yMaxArea)
}

func (self *Quadtree) Stabilize() {
	// clear nodes and cellIters, but only pack boxes
	self.cellIters.Clear()
	self.nodesFreeIndex = -1
	self.nodes = self.nodes[0 : 4]
	for i := 0; i < 4; i++ { self.nodes[i] = 0 }
	self.markedBoxes.Pack()

	// re-add each box
	xMin, xMax := self.workingArea.XMin(), self.workingArea.XMax()
	yMin, yMax := self.workingArea.YMin(), self.workingArea.YMax()
	for i, markedBox := range self.markedBoxes.list {
		spaceBoxChecks(markedBox.box, self.workingArea)
		self.quadAdd(markedBox.box, i, 0, 1, xMin, xMax, yMin, yMax)
	}
}

func (self *Quadtree) Collision(box Box) Box {
	xMin, xMax := self.workingArea.XMin(), self.workingArea.XMax()
	yMin, yMax := self.workingArea.YMin(), self.workingArea.YMax()
	return self.quadCollision(box, 0, xMin, xMax, yMin, yMax)
}

func (self *Quadtree) quadCollision(box Box, quadIdx, xMin, xMax, yMin, yMax int) Box {
	centerInnerX := xMin + ((xMax - xMin) >> 1)
	centerInnerY := yMin + ((yMax - yMin) >> 1)
	xMinChild, xMaxChild, yMinChild, yMaxChild := xMin, centerInnerX, yMin, centerInnerY
	clbx := self.nodeCollision(box, quadIdx + 0, xMinChild, xMaxChild, yMinChild, yMaxChild)
	if clbx != nil { return clbx }
	xMinChild, xMaxChild = centerInnerX + 1, xMax
	clbx = self.nodeCollision(box, quadIdx + 1, xMinChild, xMaxChild, yMinChild, yMaxChild)
	if clbx != nil { return clbx }
	yMinChild, yMaxChild = centerInnerY + 1, yMax
	clbx = self.nodeCollision(box, quadIdx + 2, xMinChild, xMaxChild, yMinChild, yMaxChild)
	if clbx != nil { return clbx }
	xMinChild, xMaxChild = xMin, centerInnerX
	return self.nodeCollision(box, quadIdx + 3, xMinChild, xMaxChild, yMinChild, yMaxChild)
}

func (self *Quadtree) nodeCollision(box Box, nodeIdx, xMin, xMax, yMin, yMax int) Box {
	if !BoxCollidesWith(box, xMin, xMax, yMin, yMax) { return nil }
	nodeValue := self.nodes[nodeIdx]
	if nodeValue == 0 { return nil } // uninitialized leaf/cell node
	if nodeValue > 0 {
		return self.quadCollision(box, nodeValue, xMin, xMax, yMin, yMax)
	} else {
		var boxIndex int
		iterIndex := ^nodeValue
		for iterIndex != -1 {
			boxIndex, iterIndex = self.cellIters.Next(iterIndex)
			cbox := self.markedBoxes.CollisionAt(box, boxIndex)
			if cbox != nil { return cbox }
		}
		return nil
	}
}

func (self *Quadtree) EachCollision(box Box, eachFunc func(Box) SearchControl) {
	self.markedBoxes.IncNoDupIndex()
	self.helperIndex = 0 // when 1, we treat it as abort the search
	xMin, xMax := self.workingArea.XMin(), self.workingArea.XMax()
	yMin, yMax := self.workingArea.YMin(), self.workingArea.YMax()
	self.quadEachCollision(box, eachFunc, 0, xMin, xMax, yMin, yMax)
}

func (self *Quadtree) quadEachCollision(box Box, eachFunc func(Box) SearchControl, quadIdx, xMin, xMax, yMin, yMax int) {
	centerInnerX := xMin + ((xMax - xMin) >> 1)
	centerInnerY := yMin + ((yMax - yMin) >> 1)
	xMinChild, xMaxChild, yMinChild, yMaxChild := xMin, centerInnerX, yMin, centerInnerY
	self.nodeEachCollision(box, eachFunc, quadIdx + 0, xMinChild, xMaxChild, yMinChild, yMaxChild)
	xMinChild, xMaxChild = centerInnerX + 1, xMax
	self.nodeEachCollision(box, eachFunc, quadIdx + 1, xMinChild, xMaxChild, yMinChild, yMaxChild)
	yMinChild, yMaxChild = centerInnerY + 1, yMax
	self.nodeEachCollision(box, eachFunc, quadIdx + 2, xMinChild, xMaxChild, yMinChild, yMaxChild)
	xMinChild, xMaxChild = xMin, centerInnerX
	self.nodeEachCollision(box, eachFunc, quadIdx + 3, xMinChild, xMaxChild, yMinChild, yMaxChild)
}

func (self *Quadtree) nodeEachCollision(box Box, eachFunc func(Box) SearchControl, nodeIdx, xMin, xMax, yMin, yMax int) {
	if self.helperIndex != 0 { return } // search aborted
	if !BoxCollidesWith(box, xMin, xMax, yMin, yMax) { return }

	nodeValue := self.nodes[nodeIdx]
	if nodeValue == 0 { return } // uninitialized leaf/cell node
	if nodeValue > 0 {
		self.quadEachCollision(box, eachFunc, nodeValue, xMin, xMax, yMin, yMax)
		return
	}

	// leaf/cell case
	var boxIndex int
	iterIndex := ^nodeValue
	for iterIndex != -1 {
		boxIndex, iterIndex = self.cellIters.Next(iterIndex)
		cbox := self.markedBoxes.CollisionNoDupAt(box, boxIndex)
		if cbox != nil && eachFunc(cbox) == SearchStop {
			self.helperIndex = 1
			return // search aborted
		}
	}
}

func (self *Quadtree) Remove(box Box) bool {
	_ = self.removeInternal(box)
	if self.helperIndex == -1 { return false }
	self.markedBoxes.RemoveBoxAt(self.helperIndex)
	return true
}

func (self *Quadtree) removeInternal(box Box) Box {
	self.helperIndex = -1
	xMin, xMax := self.workingArea.XMin(), self.workingArea.XMax()
	yMin, yMax := self.workingArea.YMin(), self.workingArea.YMax()
	_ = self.quadRemove(box, 0, xMin, xMax, yMin, yMax)
	if self.helperIndex == -1 { return nil }
	return self.markedBoxes.GetBoxAt(self.helperIndex)
}

func (self *Quadtree) quadRemove(box Box, quadIdx, xMin, xMax, yMin, yMax int) bool {
	centerInnerX := xMin + ((xMax - xMin) >> 1)
	centerInnerY := yMin + ((yMax - yMin) >> 1)
	xMinChild, xMaxChild, yMinChild, yMaxChild := xMin, centerInnerX, yMin, centerInnerY
	emptySub1 := self.nodeRemove(box, quadIdx + 0, xMinChild, xMaxChild, yMinChild, yMaxChild)
	xMinChild, xMaxChild = centerInnerX + 1, xMax
	emptySub2 := self.nodeRemove(box, quadIdx + 1, xMinChild, xMaxChild, yMinChild, yMaxChild)
	yMinChild, yMaxChild = centerInnerY + 1, yMax
	emptySub3 := self.nodeRemove(box, quadIdx + 2, xMinChild, xMaxChild, yMinChild, yMaxChild)
	xMinChild, xMaxChild = xMin, centerInnerX
	emptySub4 := self.nodeRemove(box, quadIdx + 3, xMinChild, xMaxChild, yMinChild, yMaxChild)
	return emptySub1 && emptySub2 && emptySub3 && emptySub4
}

// Returns true if the node is or becomes an empty/unitialized leaf
func (self *Quadtree) nodeRemove(box Box, nodeIdx, xMin, xMax, yMin, yMax int) bool {
	nodeValue := self.nodes[nodeIdx]
	if nodeValue == 0 { return true }
	if !BoxCollidesWith(box, xMin, xMax, yMin, yMax) { return false }
	if nodeValue > 0 {
		// interpret nodeValue as index to next quad of nodes
		emptySubtree := self.quadRemove(box, nodeValue, xMin, xMax, yMin, yMax)
		if !emptySubtree { return false }

		// mark as empty leaf
		self.nodes[nodeIdx] = 0
		self.nodes[nodeValue] = self.nodesFreeIndex
		self.nodesFreeIndex = nodeValue
		return true
	} else {
		// negative value, leaf/cell reached
		iterIndex := ^nodeValue
		prevIndex := -1
		for iterIndex != -1 {
			boxIndex, nextIter := self.cellIters.Next(iterIndex)
			if boxIndex == self.helperIndex {
				return self.removeCellIter(nodeIdx, iterIndex, prevIndex, nextIter)
			} else if self.helperIndex == -1 {
				if self.markedBoxes.BoxAtIndexEquals(box, boxIndex) {
					self.helperIndex = boxIndex
					return self.removeCellIter(nodeIdx, iterIndex, prevIndex, nextIter)
				}
			}

			// prepare next index
			prevIndex = iterIndex
			iterIndex = nextIter
		}

		return false // this cell is not empty
	}
}

// ---- helper functions ----

func (self *Quadtree) registerQuad() int {
	if self.nodesFreeIndex != -1 {
		// reuse existing free index
		newIndex := self.nodesFreeIndex
		self.nodesFreeIndex = self.nodes[newIndex]
		for i := 0; i < 4; i++ {
			self.nodes[newIndex + i] = 0 // mark as uninitialized leaf/cell
		}
		return newIndex
	} else {
		// allocate unless we still have capacity left
		newIndex := len(self.nodes)
		if cap(self.nodes) > newIndex + 3 {
			self.nodes = self.nodes[0 : newIndex + 4]
			for i := 0; i < 4; i++ {
				self.nodes[newIndex + i] = 0 // mark as uninitialized leaf/cell
			}
		} else {
			self.nodes = append(self.nodes, []int{0, 0, 0, 0}...)
		}
		return newIndex
	}
}

// returns true if after the removal the cell is now completely empty
func (self *Quadtree) removeCellIter(leafNodeIndex, cellIterIndex, prevIterIndex, nextIterIndex int) bool {
	self.cellIters.CutIter(cellIterIndex, prevIterIndex)
	if prevIterIndex == -1 {
		if nextIterIndex == -1 {
			self.nodes[leafNodeIndex] = 0 // uninitialized leaf node index
			return true
		} else {
			self.nodes[leafNodeIndex] = ^nextIterIndex
		}
	}
	return false
}

func (self *Quadtree) cellHasRoom(cellIterIndex int) bool {
	boxRoom := self.splitCutoff
	for cellIterIndex != -1 {
		boxRoom -= 1
		if boxRoom <= 0 { return false }
		_, cellIterIndex = self.cellIters.Next(cellIterIndex)
	}
	return true
}

// --- DEBUG ---

func (self *Quadtree) DebugString() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("root " + BoxString(self.workingArea) + "\n")
	currentNodes := make([]*nodePrintHelper, 0, 64)
	currentNodes = self.appendPrintHelperNodes(currentNodes, self.workingArea, 0, -1)
	nextNodes := make([]*nodePrintHelper, 0, 64)
	depth := 1

	for len(currentNodes) > 0 {
		for _, iHelper := range currentNodes {
			nodeValue := self.nodes[iHelper.index]
			if nodeValue > 0 {
				strBuilder.WriteString(iHelper.parentStr() + " -> inode @" + strconv.Itoa(iHelper.index) + " / depth " + strconv.Itoa(depth) + " / area " + BoxString(iHelper.box) + "\n")
				nextNodes = self.appendPrintHelperNodes(nextNodes, iHelper.box, nodeValue, iHelper.index)
			} else { // leaf
				strBuilder.WriteString(iHelper.parentStr() + " -> leaf @" + strconv.Itoa(iHelper.index) + " / depth " + strconv.Itoa(depth) + " / area " + BoxString(iHelper.box) + " ||")
				if nodeValue == 0 {
					strBuilder.WriteString(" []")
				} else {
					iterIndex := ^nodeValue
					for iterIndex != -1 {
						nodeValue, iterIndex = self.cellIters.Next(iterIndex)
						strBuilder.WriteString(" (" + strconv.Itoa(nodeValue) + ")" + BoxString(self.markedBoxes.GetBoxAt(nodeValue)))
					}
				}
				strBuilder.WriteString("\n")
			}
		}
		currentNodes = nextNodes
		nextNodes = nextNodes[0 : 0]
		depth += 1
	}
	return strBuilder.String()
}
type nodePrintHelper struct { box Box ; index int ; parent int }
func (self *nodePrintHelper) parentStr() string {
	if self.parent == -1 { return "root" }
	return strconv.Itoa(self.parent)
}
func (self *Quadtree) appendPrintHelperNodes(slice []*nodePrintHelper, box Box, index int, parent int) []*nodePrintHelper {
	centerInnerX := box.XMin() + ((box.XMax() - box.XMin()) >> 1)
	centerInnerY := box.YMin() + ((box.YMax() - box.YMin()) >> 1)
	xMin, xMax, yMin, yMax := box.XMin(), centerInnerX, box.YMin(), centerInnerY
	slice = append(slice, &nodePrintHelper { NewBox(xMin, xMax, yMin, yMax), index + 0, parent })
	xMin, xMax = centerInnerX + 1, box.XMax()
	slice = append(slice, &nodePrintHelper { NewBox(xMin, xMax, yMin, yMax), index + 1, parent })
	yMin, yMax = centerInnerY + 1, box.YMax()
	slice = append(slice, &nodePrintHelper { NewBox(xMin, xMax, yMin, yMax), index + 2, parent })
	xMin, xMax = box.XMin(), centerInnerX
	slice = append(slice, &nodePrintHelper { NewBox(xMin, xMax, yMin, yMax), index + 3, parent })
	return slice
}
