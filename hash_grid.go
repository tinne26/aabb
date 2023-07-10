package aabb

import "strings"
import "strconv"

// List of possible optimizations:
// - Specialized HashGrid32 for 32 bit tile sizes. Aligning working
//   area origin to 0 doesn't yield any significant optimizations.
//   But I might be able to not use limited indices, by having everything
//   revolve around 0, 0 and have the 4 quadrants be interleaved in
//   the map. I think I can do that with the regular hashGrid too,
//   unlike with regular grids. Indexing might get tricky, but should
//   be doable. Maybe I have to partition operations on all 4 spaces
//   though. But for an infinite grid, it might not be bad.
// - Do not initialize cellIter on the first addition, and use ^index
//   to refer to cell iters vs direct box indices.
// - Make MutateBox operate only on the changing areas, instead of
//   fully remove and fully add again.
// - Handle each search and collision case in a finer way, separating
//   single cell accesses from multiple horizontal cells, multiple
//   vertical cells, or finally all cells at once (tried it, but it
//   only seems to reduce execution time in 1-2%, too ugly)
// - Reduce memory usage by using int32 instead of int (tried briefly
//   but didn't observe any obvious performance difference, maybe you
//   need to add many many boxes in order for the index size to matter)

// A [Grid] that uses a hash instead of a full fledged slice. More
// memory efficient if the grid has a large number of empty cells,
// but it's between 10-15% slower than a regular grid. Still more
// than twice as fast as quadtrees for most use-cases.
//
// See also [EndlessHashGrid].
type HashGrid struct {
	xMin, xMax, yMin, yMax int
	cellWidth  int
	cellHeight int
	horzCells  int

	cellIters cellIterList
	markedBoxes markedBoxList
	cells map[int]int // each cell points to a cellIter
}

func NewHashGrid(workingArea Box, horzCells int, vertCells int) *HashGrid {
	if horzCells < 1 { panic("horzCells < 1") }
	if vertCells < 1 { panic("vertCells < 1") }
	areaWidth  := BoxWidth( workingArea)
	areaHeight := BoxHeight(workingArea)
	if horzCells > areaWidth  { panic("horzCells > workingArea width" ) }
	if vertCells > areaHeight { panic("vertCells > workingArea height") }
	if areaWidth  % horzCells != 0 { panic("workingArea width not multiple of horzCells" ) }
   if areaHeight % vertCells != 0 { panic("workingArea height not multiple of vertCells") }

	cellWidth  := areaWidth/horzCells
   cellHeight := areaHeight/vertCells
	return &HashGrid {
		xMin: workingArea.XMin(), xMax: workingArea.XMax(),
		yMin: workingArea.YMin(), yMax: workingArea.YMax(),
		cellWidth: cellWidth, cellHeight: cellHeight,
		horzCells: horzCells,
		cellIters: newCellIterList(64),
		markedBoxes: newMarkedBoxList(32),
		cells: make(map[int]int),
	}
}

func (self *HashGrid) DebugString() string {
	index := 0
	var strBuilder strings.Builder
	vertCells := ((self.yMax - self.yMin) + 1)/self.cellHeight
	for y := 0; y < vertCells; y++ {
		for x := 0; x < self.horzCells; x++ {
			strBuilder.WriteString("cell " + strconv.Itoa(index) + " (" + strconv.Itoa(x) + "X, " + strconv.Itoa(y) + "Y):")
			iter, found := self.cells[index]
			if found {
				for iter != -1 {
					boxIndex, nextIter := self.cellIters.Next(iter)
					strBuilder.WriteRune(' ')
					strBuilder.WriteString(BoxString(self.markedBoxes.GetBoxAt(boxIndex)))
					iter = nextIter
				}
				strBuilder.WriteRune('\n')
			} else {
				strBuilder.WriteString(" {empty}\n")
			}
			index += 1
		}
	}
	return strBuilder.String()
}

func (self *HashGrid) Clear() {
	self.markedBoxes.Clear()
	self.cellIters.Clear()
	for key, _ := range self.cells {
		delete(self.cells, key)
	}
}

func (self *HashGrid) Add(box Box) {
	boxIndex := self.markedBoxes.AddBox(box)
	if boxIndex < 0 { panic("boxIndex < 0") }
	self.innerAdd(box, boxIndex)
}

// Assumes the box is already registered, and adds it to the proper
// cells with the given index.
func (self *HashGrid) innerAdd(box Box, boxIndex int) {
	spaceBoxChecksWith(box, self.xMin, self.xMax, self.yMin, self.xMax)
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex, found := self.cells[cellIndex]
			if !found { cellIterIndex = -1 }
			self.cells[cellIndex] = self.cellIters.AddIterTo(boxIndex, cellIterIndex)
			cellIndex += 1
		}
		cellIndex += rowStride
	}
}

func (self *HashGrid) Collision(box Box) Box {
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	xMin, xMax, yMin, yMax := box.XMin(), box.XMax(), box.YMin(), box.YMax()
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex, found := self.cells[cellIndex]
			if found {
				for cellIterIndex != -1 {
					var boxIndex int
					boxIndex, cellIterIndex = self.cellIters.Next(cellIterIndex)
					cbox := self.markedBoxes.UnrolledCollisionAt(box, xMin, xMax, yMin, yMax, boxIndex)
					if cbox != nil { return cbox }
				}
			}
			cellIndex += 1
		}
		cellIndex += rowStride
	}
	return nil
}

func (self *HashGrid) EachCollision(box Box, eachFunc func(Box) SearchControl) {
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	self.markedBoxes.IncNoDupIndex()
	xMin, xMax, yMin, yMax := box.XMin(), box.XMax(), box.YMin(), box.YMax()
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex, found := self.cells[cellIndex]
			if found {
				for cellIterIndex != -1 {
					var boxIndex int
					boxIndex, cellIterIndex = self.cellIters.Next(cellIterIndex)
					cbox := self.markedBoxes.UnrolledCollisionNoDupAt(box, xMin, xMax, yMin, yMax, boxIndex)
					if cbox == nil { continue }
					searchControl := eachFunc(cbox)
					if searchControl == SearchStop { return }
				}
			}
			cellIndex += 1
		}
		cellIndex += rowStride
	}
}

func (self *HashGrid) Remove(box Box) bool {
	boxIndex := self.innerRemove(box)
	if boxIndex == -1 { return false }
	self.markedBoxes.RemoveBoxAt(boxIndex)
	return true
}

// innerRemove removes first ocurrence of the given box within any cell,
// but doesn't remove the actual box inside self.boxes, it only returns
// the index of the box to be removed
func (self *HashGrid) innerRemove(box Box) int {
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	removedBoxIndex := -1
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex, found := self.cells[cellIndex]
			if found {
				if removedBoxIndex == -1 {
					removedBoxIndex = self.removeFirstEqualBoxInCell(box, cellIndex, cellIterIndex)
				} else {
					self.removeBoxInCellByIndex(removedBoxIndex, cellIndex, cellIterIndex)
				}
			}
			cellIndex += 1
		}
		cellIndex += rowStride
	}

	return removedBoxIndex
}

func (self *HashGrid) removeFirstEqualBoxInCell(box Box, cellIndex int, cellIterIndex int) int {
	prevCellIterIndex := -1
	for cellIterIndex != -1 {
		boxIndex, nextIterIndex := self.cellIters.Next(cellIterIndex)
		if self.markedBoxes.BoxAtIndexEquals(box, boxIndex) {
			self.removeCellIter(cellIndex, cellIterIndex, prevCellIterIndex, nextIterIndex)
			return boxIndex
		}
		prevCellIterIndex = cellIterIndex
		cellIterIndex = nextIterIndex
	}
	return -1
}

func (self *HashGrid) removeBoxInCellByIndex(knownBoxIndex int, cellIndex int, cellIterIndex int) {
	prevCellIterIndex := -1
	for cellIterIndex != -1 {
		boxIndex, nextIterIndex := self.cellIters.Next(cellIterIndex)
		if boxIndex == knownBoxIndex {
			self.removeCellIter(cellIndex, cellIterIndex, prevCellIterIndex, nextIterIndex)
			return
		}
		prevCellIterIndex = cellIterIndex
		cellIterIndex = nextIterIndex
	}
}

func (self *HashGrid) removeCellIter(cellIndex, cellIterIndex, prevCellIterIndex, nextIterIndex int) {
	self.cellIters.CutIter(cellIterIndex, prevCellIterIndex)
	if prevCellIterIndex == -1 {
		if nextIterIndex == -1 {
			delete(self.cells, cellIndex)
		} else {
			self.cells[cellIndex] = nextIterIndex
		}
	}
}

func (self *HashGrid) MutateBox(box MutableBox, xMin, xMax, yMin, yMax int) {
	boxIndex := self.innerRemove(box)
	if boxIndex == -1 {
		panic("can't update box, it can't be found in the HashGrid")
	}
	self.markedBoxes.MutateBoxAt(boxIndex, xMin, xMax, yMin, yMax)
	self.innerAdd(NewBox(xMin, xMax, yMin, yMax), boxIndex)
}

func (self *HashGrid) Stabilize() {
	// don't remove boxes, but clear everything else
	self.markedBoxes.Pack()
	self.cellIters.Clear()
	for key, _ := range self.cells {
		delete(self.cells, key)
	}

	// re-add each box
	for boxIndex, markedBoxObj := range self.markedBoxes.list {
		self.innerAdd(markedBoxObj.box, boxIndex)
	}
}

// --- helpers ---

func (self *HashGrid) getCellIndices(box Box) (int, int, int, int) {
	xMinCell := (box.XMin() - self.xMin)/self.cellWidth
   xMaxCell := (box.XMax() - self.xMin)/self.cellWidth
	yMinCell := (box.YMin() - self.yMin)/self.cellHeight
   yMaxCell := (box.YMax() - self.yMin)/self.cellHeight
	return xMinCell, xMaxCell, yMinCell, yMaxCell
}
