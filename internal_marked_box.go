package aabb

type markedBox struct {
	box Box
	mark int
}

type markedBoxList struct {
	list []markedBox
	freeIndex int
	noDupIndex int // index used on searches to avoid duplicate results
}

func newMarkedBoxList(initialCapacity int) markedBoxList {
	return markedBoxList {
		list: make([]markedBox, 0, initialCapacity),
		freeIndex: -1,
		noDupIndex: 1,
	}
}

func (self *markedBoxList) Clear() {
	self.list = self.list[ : 0]
	self.freeIndex = -1
}

func (self *markedBoxList) Pack() {
	if self.freeIndex == -1 { return }

	self.freeIndex = -1
	i := 0
	listLen := len(self.list)
	for i < listLen {
		if self.list[i].box != nil { i += 1 ; continue }
		listLen -= 1
		self.list[listLen], self.list[i] = self.list[i], self.list[listLen]
		self.list = self.list[0 : listLen]
	}
}

func (self *markedBoxList) MutateBoxAt(boxIndex, xMin, xMax, yMin, yMax int) {
	(self.list[boxIndex].box).(MutableBox).MutateBox(xMin, xMax, yMin, yMax)
}

func (self *markedBoxList) GetBoxAt(boxIndex int) Box {
	return self.list[boxIndex].box
}

func (self *markedBoxList) ExtractBoxAt(boxIndex int) Box {
	boxRemoved := self.list[boxIndex].box
	self.RemoveBoxAt(boxIndex)
	return boxRemoved
}

func (self *markedBoxList) AddBox(box Box) int {
	if self.freeIndex != -1 {
		// reuse existing free index
		newIndex := self.freeIndex
		self.freeIndex = self.list[newIndex].mark
		self.list[newIndex] = markedBox { box, 0 }
		return newIndex
	} else { // no free index
		// allocate unless we still have capacity to expand
		newIndex := len(self.list)
		if cap(self.list) > newIndex {
			self.list = self.list[0 : newIndex + 1]
			self.list[newIndex] = markedBox{ box, 0 }
		} else {
			self.list = append(self.list, markedBox{ box, 0 })
		}
		return newIndex
	}
}

func (self *markedBoxList) RemoveBoxAt(boxIndex int) {
	self.list[boxIndex] = markedBox{ nil, self.freeIndex }
	self.freeIndex = boxIndex
}

func (self *markedBoxList) CollisionAt(box Box, boxIndex int) Box {
	targetBox := self.list[boxIndex].box
	if BoxesCollide(box, targetBox) && box != targetBox {
		return targetBox
	}
	return nil
}

func (self *markedBoxList) UnrolledCollisionAt(box Box, xMin, xMax, yMin, yMax, boxIndex int) Box {
	targetBox := self.list[boxIndex].box
	if BoxCollidesWith(targetBox, xMin, xMax, yMin, yMax) && box != targetBox {
		return targetBox
	}
	return nil
}

func (self *markedBoxList) BoxAtIndexEquals(box Box, boxIndex int) bool {
	return self.list[boxIndex].box == box
}

// Same as `CollisionAt`, but it ignores the collision if the current
// `noDupIndex` is the same as `markedBox.mark`. Use `IncNoDupIndex()`
// before starting a search to refresh the no-duplication filter.
func (self *markedBoxList) CollisionNoDupAt(box Box, boxIndex int) Box {
	selfBox := self.list[boxIndex]
	if selfBox.mark == self.noDupIndex { return nil }
	if BoxesCollide(box, selfBox.box) && box != selfBox.box {
		self.list[boxIndex].mark = self.noDupIndex
		return selfBox.box
	}
	return nil
}

func (self *markedBoxList) UnrolledCollisionNoDupAt(box Box, xMin, xMax, yMin, yMax, boxIndex int) Box {
	selfBox := self.list[boxIndex]
	if selfBox.mark == self.noDupIndex { return nil }
	if BoxCollidesWith(selfBox.box, xMin, xMax, yMin, yMax) && box != selfBox.box {
		self.list[boxIndex].mark = self.noDupIndex
		return selfBox.box
	}
	return nil
}

// Call before starting a search that will use `CollisionNoDupAt`,
// so the internal control for no duplicates can be adjusted.
func (self *markedBoxList) IncNoDupIndex() {
	self.noDupIndex += 1

	// handle overflows by resetting all box marks
	if self.noDupIndex < 0 {
		for i := 0; i < len(self.list); i++ {
			if self.list[i].box != nil {
				self.list[i].mark = 0
			}
		}
		self.noDupIndex = 1
	}
}
