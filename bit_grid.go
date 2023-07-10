package aabb

import "strconv"

// Possible optimizations:
// - With go 1.18, maybe use generic types for uint8, uint16, uint32, uint64
// - Optimize for 32px grids

type bitGridBox struct {
   xBits uint32
   yBits uint32
   box Box
}

func (self *bitGridBox) CollidesWith(other *bitGridBox) bool {
	return (self.xBits & other.xBits) != 0 &&
	       (self.yBits & other.yBits) != 0 &&
			 BoxesCollide(self.box, other.box)
}

func (self *bitGridBox) RefreshGridBits(workingArea Box, horzCellSize int, vertCellSize int) {
   // safety checks
	spaceBoxChecks(self.box, workingArea)

	// compute x bits
   xMinCell := uint32((self.box.XMin() - workingArea.XMin())/horzCellSize)
   xMaxCell := uint32((self.box.XMax() - workingArea.XMin())/horzCellSize)
	if xMaxCell >= MaxBitGridCells { panic("xMaxCell >= MaxBitGridCells") }

   if xMaxCell == xMinCell {
		self.xBits = 1 << xMaxCell
   } else {
      self.xBits = ((1 << (xMaxCell + 1)) - 1) ^ ((1 << xMinCell) - 1)
   }

   // compute y bits
   yMinCell := uint32((self.box.YMin() - workingArea.YMin())/vertCellSize)
   yMaxCell := uint32((self.box.YMax() - workingArea.YMin())/vertCellSize)
	if yMaxCell >= MaxBitGridCells { panic("yMaxCell >= MaxBitGridCells") }

   if yMaxCell == yMinCell {
		self.yBits = 1 << yMaxCell
   } else {
      self.yBits = ((1 << (yMaxCell + 1)) - 1) ^ ((1 << yMinCell) - 1)
   }
}

// BitGrid implements the naive, O(n^2) algorithm to check collisions
// between Boxes... but with an optimization: instead of comparing Boxes
// one by one, it defines two flagsets for each [Box] so faster comparisons
// can be done through bitwise operations. It's more than twice as fast
// as the brute force approach.
type BitGrid struct {
   boxes []*bitGridBox
   workingArea Box
   horzCellSize int
   vertCellSize int
   tmpBitGridBox *bitGridBox
}

const MaxBitGridCells = 32
func NewBitGrid(workingArea Box, horzCells int, vertCells int) *BitGrid {
   if horzCells < 1 || horzCells > MaxBitGridCells {
      panicMsg := "horzCells out of valid range 1 <= "
      panicMsg += strconv.Itoa(horzCells) + " <= "
      panicMsg += strconv.Itoa(MaxBitGridCells)
      panic(panicMsg)
   }
   if vertCells < 1 || vertCells > MaxBitGridCells {
      panicMsg := "vertCells out of valid range 1 <= "
      panicMsg += strconv.Itoa(vertCells) + " <= "
      panicMsg += strconv.Itoa(MaxBitGridCells)
      panic(panicMsg)
   }

   workAreaWidth  := BoxWidth(workingArea)
   workAreaHeight := BoxHeight(workingArea)
   if workAreaWidth % horzCells != 0 {
      panic("workingArea width not multiple of horzCells")
   }
   if workAreaHeight % vertCells != 0 {
      panic("workingArea height not multiple of vertCells")
   }

   horzCellSize := workAreaWidth/horzCells
   vertCellSize := workAreaHeight/vertCells
   return &BitGrid {
      boxes: make([]*bitGridBox, 0, 64),
      workingArea: workingArea,
      horzCellSize: horzCellSize,
      vertCellSize: vertCellSize,
      tmpBitGridBox: &bitGridBox{},
   }
}

func (self *BitGrid) Clear() { self.boxes = self.boxes[0 : 0] }

func (self *BitGrid) Add(box Box) {
	spaceBoxChecks(box, self.workingArea)

   size := len(self.boxes)
   if cap(self.boxes) > size {
      self.boxes = self.boxes[0 : size + 1]
      if self.boxes[size] != nil {
         self.boxes[size].box = box
         self.boxes[size].RefreshGridBits(self.workingArea, self.horzCellSize, self.vertCellSize)
      } else {
         bitGridBox := &bitGridBox { 0, 0, box }
         bitGridBox.RefreshGridBits(self.workingArea, self.horzCellSize, self.vertCellSize)
         self.boxes[size] = bitGridBox
      }
   } else {
      bitGridBox := &bitGridBox { 0, 0, box }
      bitGridBox.RefreshGridBits(self.workingArea, self.horzCellSize, self.vertCellSize)
      self.boxes = append(self.boxes, bitGridBox)
   }
}

func (self *BitGrid) MutateBox(box MutableBox, xMin, xMax, yMin, yMax int) {
	// NOTE: no need for initial spaceBoxChecks, they happen in RefreshGridBits
	for _, otherBox := range self.boxes {
		if box == otherBox.box {
			(otherBox.box).(MutableBox).MutateBox(xMin, xMax, yMin, yMax)
			otherBox.RefreshGridBits(self.workingArea, self.horzCellSize, self.vertCellSize)
			return
		}
	}
	panic("box to be updated not found")
}

func (self *BitGrid) Stabilize() {
	for _, box := range self.boxes {
		box.RefreshGridBits(self.workingArea, self.horzCellSize, self.vertCellSize)
	}
}

func (self *BitGrid) Collision(box Box) Box {
   self.setupTmpBitGridBox(box)
   for _, otherBox := range self.boxes {
      if self.tmpBitGridBox.CollidesWith(otherBox) && box != otherBox.box {
         return otherBox.box
      }
   }
   return nil
}
func (self *BitGrid) EachCollision(box Box, eachFunc func(Box) SearchControl) {
	self.setupTmpBitGridBox(box)
   for _, otherBox := range self.boxes {
      if self.tmpBitGridBox.CollidesWith(otherBox) && box != otherBox.box {
			searchControl := eachFunc(otherBox.box)
         if searchControl == SearchStop { return }
      }
   }
}

func (self *BitGrid) Remove(box Box) bool {
   for i, otherBox := range self.boxes {
      if otherBox.box == box {
         lastIdx := len(self.boxes) - 1 // can't be negative, we have a box
         self.boxes[i], self.boxes[lastIdx] = self.boxes[lastIdx], self.boxes[i]
         self.boxes = self.boxes[0 : lastIdx]
         return true
      }
   }
   return false
}

func (self *BitGrid) setupTmpBitGridBox(box Box) {
	self.tmpBitGridBox.box = box
	self.tmpBitGridBox.RefreshGridBits(self.workingArea, self.horzCellSize, self.vertCellSize)
}
