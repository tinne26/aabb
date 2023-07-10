package aabb

// Almost 20% slower than [HashGrid], but I hope future compiler optimizations
// will reduce the gap, as it's really hard to explain why this is so much
// slower than [HashGrid]. Still in line with or slightly faster than quadtrees,
// and fairly competitive approach all around. Not the best at anything, but
// not terrible at anything (though general [Grid]-like caveats apply).
type EndlessHashGrid struct {
	cellWidth  int
	cellHeight int
	cellIters cellIterList
	markedBoxes markedBoxList
	cells map[uint64]int // each cell points to a cellIter
}

func NewEndlessHashGrid(cellWidth int, cellHeight int) *EndlessHashGrid {
	if cellWidth  < 1 { panic("cellWidth < 1" ) }
	if cellHeight < 1 { panic("cellHeight < 1") }
	return &EndlessHashGrid {
		cellWidth  : cellWidth,
		cellHeight : cellHeight,
		cellIters  : newCellIterList(64),
		markedBoxes: newMarkedBoxList(32),
		cells      : make(map[uint64]int),
	}
}

func (self *EndlessHashGrid) Clear() {
	self.markedBoxes.Clear()
	self.cellIters.Clear()
	for key, _ := range self.cells {
		delete(self.cells, key)
	}
}

func (self *EndlessHashGrid) Add(box Box) {
	self.pointBoxToItsCells(box, self.markedBoxes.AddBox(box))
}

func (self *EndlessHashGrid) pointBoxToItsCells(box Box, boxIndex int) {
	BoxPanicIfInvalid(box)
	xMinCell, xMaxCell := getValuesCells(box.XMin(), box.XMax(), self.cellWidth)
	yMinCell, yMaxCell := getValuesCells(box.YMin(), box.YMax(), self.cellHeight)
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIndex := (uint64(uint32(x)) << 32) | uint64(uint32(y))
			cellIterIndex, found := self.cells[cellIndex]
			if !found { cellIterIndex = -1 } // -1 indicates no previous cellIter
			self.cells[cellIndex] = self.cellIters.AddIterTo(boxIndex, cellIterIndex)
		}
	}
}

func (self *EndlessHashGrid) Collision(box Box) Box {
	xMin, xMax, yMin, yMax := box.XMin(), box.XMax(), box.YMin(), box.YMax()
	xMinCell, xMaxCell := getValuesCells(xMin, xMax, self.cellWidth)
	yMinCell, yMaxCell := getValuesCells(yMin, yMax, self.cellHeight)
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIndex := (uint64(uint32(x)) << 32) | uint64(uint32(y))
			cellIterIndex, found := self.cells[cellIndex]
			if !found { continue }
			for cellIterIndex != -1 {
				var boxIndex int
				boxIndex, cellIterIndex = self.cellIters.Next(cellIterIndex)
				xx := self.markedBoxes.UnrolledCollisionAt(box, xMin, xMax, yMin, yMax, boxIndex)
				if xx != nil { return xx }
			}
		}
	}
	return nil
}

func (self *EndlessHashGrid) EachCollision(box Box, eachFunc func(Box) SearchControl) {
	self.markedBoxes.IncNoDupIndex()
	xMin, xMax, yMin, yMax := box.XMin(), box.XMax(), box.YMin(), box.YMax()
	xMinCell, xMaxCell := getValuesCells(xMin, xMax, self.cellWidth)
	yMinCell, yMaxCell := getValuesCells(yMin, yMax, self.cellHeight)
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIndex := (uint64(uint32(x)) << 32) | uint64(uint32(y))
			cellIterIndex, found := self.cells[cellIndex]
			if !found { continue }
			for cellIterIndex != -1 {
				var boxIndex int
				boxIndex, cellIterIndex = self.cellIters.Next(cellIterIndex)
				xx := self.markedBoxes.UnrolledCollisionNoDupAt(box, xMin, xMax, yMin, yMax, boxIndex)
				if xx == nil { continue }
				if eachFunc(xx) == SearchStop { return }
			}
		}
	}
}

func (self *EndlessHashGrid) Remove(box Box) bool {
	boxIndex := self.unregisterEqualBoxFromItsCells(box)
	if boxIndex == -1 { return false }
	self.markedBoxes.RemoveBoxAt(boxIndex)
	return true
}

func (self *EndlessHashGrid) unregisterEqualBoxFromItsCells(box Box) int {
	xMinCell, xMaxCell := getValuesCells(box.XMin(), box.XMax(), self.cellWidth)
	yMinCell, yMaxCell := getValuesCells(box.YMin(), box.YMax(), self.cellHeight)
	equalBoxIndex := -1
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIndex := (uint64(uint32(x)) << 32) | uint64(uint32(y))
			cellIterIndex, found := self.cells[cellIndex]
			if !found { continue }

			prevCellIterIndex := -1
			for cellIterIndex != -1 {
				boxIndex, nextIterIndex := self.cellIters.Next(cellIterIndex)
				if self.hasToApplyRemove(box, boxIndex, equalBoxIndex) {
					if equalBoxIndex == -1 { equalBoxIndex = boxIndex }
					self.cellIters.CutIter(cellIterIndex, prevCellIterIndex)
					if prevCellIterIndex == -1 {
						if nextIterIndex == -1 {
							delete(self.cells, cellIndex)
						} else {
							self.cells[cellIndex] = nextIterIndex
						}
					}
					break
				}
				prevCellIterIndex = cellIterIndex
				cellIterIndex = nextIterIndex
			}
		}
	}
	return equalBoxIndex
}

func (self *EndlessHashGrid) MutateBox(box MutableBox, xMin, xMax, yMin, yMax int) {
	boxIndex := self.unregisterEqualBoxFromItsCells(box)
	if boxIndex == -1 {
		panic("can't update box, it can't be found in the EndlessHashGrid")
	}
	self.markedBoxes.MutateBoxAt(boxIndex, xMin, xMax, yMin, yMax)
	self.pointBoxToItsCells(box, boxIndex)
}

func (self *EndlessHashGrid) Stabilize() {
	// don't remove boxes, but clear everything else
	self.markedBoxes.Pack()
	self.cellIters.Clear()
	for key, _ := range self.cells {
		delete(self.cells, key)
	}

	// re-add each box
	for boxIndex, markedBoxObj := range self.markedBoxes.list {
		self.pointBoxToItsCells(markedBoxObj.box, boxIndex)
	}
}

// ---- helper functions ----

func getValuesCells(vMin, vMax, cellLength int) (int32, int32) {
	return getValueCell(vMin, cellLength), getValueCell(vMax, cellLength)
}

func getValueCell(v int, cellLength int) int32 {
	if v >= 0 { return int32(v/cellLength) }
	return int32(-((-v + 1)/cellLength))
}

func (self *EndlessHashGrid) hasToApplyRemove(box Box, boxIndex int, equalBoxIndex int) bool {
	if equalBoxIndex != -1 { return boxIndex == equalBoxIndex }
	return self.markedBoxes.BoxAtIndexEquals(box, boxIndex)
}
