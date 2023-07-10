package aabb

import "strings"
import "strconv"

// Grids are one of the fastest approaches to collisions. They
// are not perfect, though:
//  - You have to know your working area in advance (and allocate
//    memory for it).
//  - Related to the previous, if your game levels have a lot of
//    empty space and empty cells, grids can be wasteful.
//  - Big or stretched objects will be referenced from multiple cells,
//    which is not ideal from a memory-usage perspective. In fact,
//    this is the main problem for grids if your game objects have a
//    lot of variance in size and/or aspect ratio.
//
// The first two problems can be alleviated with [HashGrid] and
// [EndlessHashGrid] in exchange for some performance.
type Grid struct {
	xMin, xMax, yMin, yMax int
	cellWidth  int
	cellHeight int
	horzCells  int

	cellIters cellIterList
	markedBoxes markedBoxList
	cells []int // each cell points to a cellIter
}

func NewGrid(workingArea Box, horzCells int, vertCells int) *Grid {
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
	cells := make([]int, horzCells*vertCells)
	for i, _ := range cells { cells[i] = -1 }
	return &Grid {
		xMin: workingArea.XMin(), xMax: workingArea.XMax(),
		yMin: workingArea.YMin(), yMax: workingArea.YMax(),
		cellWidth: cellWidth, cellHeight: cellHeight,
		horzCells: horzCells,
		cellIters: newCellIterList(64),
		markedBoxes: newMarkedBoxList(32),
		cells: cells,
	}
}

func (self *Grid) DebugString() string {
	index := 0
	var strBuilder strings.Builder
	vertCells := ((self.yMax - self.yMin) + 1)/self.cellHeight
	for y := 0; y < vertCells; y++ {
		for x := 0; x < self.horzCells; x++ {
			strBuilder.WriteString("cell " + strconv.Itoa(index) + " (" + strconv.Itoa(x) + "X, " + strconv.Itoa(y) + "Y):")
			iter := self.cells[index]
			if iter == -1 {
				strBuilder.WriteString(" {empty}\n")
			} else {
				for iter != -1 {
					boxIndex, nextIter := self.cellIters.Next(iter)
					strBuilder.WriteRune(' ')
					strBuilder.WriteString(BoxString(self.markedBoxes.GetBoxAt(boxIndex)))
					iter = nextIter
				}
				strBuilder.WriteRune('\n')
			}
			index += 1
		}
	}
	return strBuilder.String()
}

func (self *Grid) Clear() {
	self.markedBoxes.Clear()
	self.cellIters.Clear()
	for i, _ := range self.cells { self.cells[i] = i }
}

func (self *Grid) Add(box Box) {
	boxIndex := self.markedBoxes.AddBox(box)
	if boxIndex < 0 { panic("boxIndex < 0") }
	self.innerAdd(box, boxIndex)
}

// Assumes the box is already registered, and adds it to the proper
// cells with the given index.
func (self *Grid) innerAdd(box Box, boxIndex int) {
	spaceBoxChecksWith(box, self.xMin, self.xMax, self.yMin, self.yMax)
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex := self.cells[cellIndex]
			self.cells[cellIndex] = self.cellIters.AddIterTo(boxIndex, cellIterIndex)
			cellIndex += 1
		}
		cellIndex += rowStride
	}
}

func (self *Grid) Collision(box Box) Box {
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	xMin, xMax, yMin, yMax := box.XMin(), box.XMax(), box.YMin(), box.YMax()
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex := self.cells[cellIndex]
			for cellIterIndex != -1 {
				var boxIndex int
				boxIndex, cellIterIndex = self.cellIters.Next(cellIterIndex)
				cbox := self.markedBoxes.UnrolledCollisionAt(box, xMin, xMax, yMin, yMax, boxIndex)
				if cbox != nil { return cbox }
			}
			cellIndex += 1
		}
		cellIndex += rowStride
	}
	return nil
}

func (self *Grid) EachCollision(box Box, eachFunc func(Box) SearchControl) {
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	self.markedBoxes.IncNoDupIndex()
	xMin, xMax, yMin, yMax := box.XMin(), box.XMax(), box.YMin(), box.YMax()
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex := self.cells[cellIndex]
			for cellIterIndex != -1 {
				var boxIndex int
				boxIndex, cellIterIndex = self.cellIters.Next(cellIterIndex)
				cbox := self.markedBoxes.UnrolledCollisionNoDupAt(box, xMin, xMax, yMin, yMax, boxIndex)
				if cbox == nil { continue }
				searchControl := eachFunc(cbox)
				if searchControl == SearchStop { return }
			}
			cellIndex += 1
		}
		cellIndex += rowStride
	}
}

func (self *Grid) Remove(box Box) bool {
	boxIndex := self.innerRemove(box)
	if boxIndex == -1 { return false }
	self.markedBoxes.RemoveBoxAt(boxIndex)
	return true
}

// innerRemove removes first ocurrence of the given box within any cell,
// but doesn't remove the actual box inside self.boxes, it only returns
// the index of the box to be removed
func (self *Grid) innerRemove(box Box) int {
	xMinCell, xMaxCell, yMinCell, yMaxCell := self.getCellIndices(box)
	removedBoxIndex := -1
	cellIndex := yMinCell*self.horzCells + xMinCell
	rowStride := self.horzCells - xMaxCell + xMinCell - 1
	for y := yMinCell; y <= yMaxCell; y++ {
		for x := xMinCell; x <= xMaxCell; x++ {
			cellIterIndex := self.cells[cellIndex]
			if cellIterIndex != -1 {
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

func (self *Grid) removeFirstEqualBoxInCell(box Box, cellIndex int, cellIterIndex int) int {
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

func (self *Grid) removeBoxInCellByIndex(knownBoxIndex int, cellIndex int, cellIterIndex int) {
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

func (self *Grid) removeCellIter(cellIndex, cellIterIndex, prevCellIterIndex, nextIterIndex int) {
	self.cellIters.CutIter(cellIterIndex, prevCellIterIndex)
	if prevCellIterIndex == -1 {
		self.cells[cellIndex] = nextIterIndex
	}
}

func (self *Grid) MutateBox(box MutableBox, xMin, xMax, yMin, yMax int) {
	boxIndex := self.innerRemove(box)
	if boxIndex == -1 {
		panic("can't update box, it can't be found in the HashGrid")
	}
	self.markedBoxes.MutateBoxAt(boxIndex, xMin, xMax, yMin, yMax)
	self.innerAdd(self.markedBoxes.GetBoxAt(boxIndex), boxIndex)
}

func (self *Grid) Stabilize() {
	// don't remove boxes, but clear everything else
	self.markedBoxes.Pack()
	self.cellIters.Clear()
	for i, _ := range self.cells { self.cells[i] = -1 }

	// re-add each box
	for boxIndex, markedBoxObj := range self.markedBoxes.list {
		self.innerAdd(markedBoxObj.box, boxIndex)
	}
}

// --- helpers ---
func (self *Grid) getCellIndices(box Box) (int, int, int, int) {
	xMinCell := (box.XMin() - self.xMin)/self.cellWidth
   xMaxCell := (box.XMax() - self.xMin)/self.cellWidth
	yMinCell := (box.YMin() - self.yMin)/self.cellHeight
   yMaxCell := (box.YMax() - self.yMin)/self.cellHeight
	return xMinCell, xMaxCell, yMinCell, yMaxCell
}
