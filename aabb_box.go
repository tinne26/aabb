package aabb

import "strconv"

// The core two-dimensional Box interface that allows performing
// broad-phase collision in AABB spaces.
type Box interface {
   XMin() int // included
   XMax() int // included
   YMin() int // included
   YMax() int // included
}
func BoxesAreEqual(a, b Box) bool {
	return a.XMin() == b.XMin() && a.XMax() == b.XMax() &&
	       a.YMin() == b.YMin() && a.YMax() == b.YMax()
}
func BoxesCollide(a, b Box) bool {
	return SegmentCollides(a.XMin(), a.XMax(), b.XMin(), b.XMax()) &&
          SegmentCollides(a.YMin(), a.YMax(), b.YMin(), b.YMax())
}
func BoxCollidesWith(box Box, xMinArea, xMaxArea, yMinArea, yMaxArea int) bool {
	return SegmentCollides(box.XMin(), box.XMax(), xMinArea, xMaxArea) &&
          SegmentCollides(box.YMin(), box.YMax(), yMinArea, yMaxArea)
}
func BoxContains(container, contained Box) bool {
	return SegmentContains(container.XMin(), container.XMax(), contained.XMin(), contained.XMax()) &&
	       SegmentContains(container.YMin(), container.YMax(), contained.YMin(), contained.YMax())
}

func AreaContains(xMinArea, xMaxArea, yMinArea, yMaxArea int, contained Box) bool {
	return SegmentContains(xMinArea, xMaxArea, contained.XMin(), contained.XMax()) &&
	       SegmentContains(yMinArea, yMaxArea, contained.YMin(), contained.YMax())
}

func BoxString(box Box) string {
	return "[ X " + strconv.Itoa(box.XMin()) + " " + strconv.Itoa(box.XMax()) + ", Y " + strconv.Itoa(box.YMin()) + " " + strconv.Itoa(box.YMax()) + " ]"
}

func BoxWidth( a Box) int { return a.XMax() - a.XMin() + 1 }
func BoxHeight(a Box) int { return a.YMax() - a.YMin() + 1 }
func BoxPanicIfInvalid(box Box) {
	if box.XMin() > box.XMax() { panic("box.XMin() > box.XMax() [" + strconv.Itoa(box.XMin()) + " > " + strconv.Itoa(box.XMax()) + "]") }
   if box.YMin() > box.YMax() { panic("box.YMin() > box.YMax() [" + strconv.Itoa(box.YMin()) + " > " + strconv.Itoa(box.YMax()) + "]") }
}

// Both segment min and max are included.
func SegmentCollides(aMin, aMax, bMin, bMax int) bool {
   if aMin <= bMin { return aMax >= bMin }
   return bMax >= aMin
}

func SegmentContains(containerMin, containerMax, containedMin, containedMax int) bool {
	return containerMin <= containerMin && containerMax >= containedMax
}

type MutableBox interface {
   Box
   MutateBox(xMin, xMax, yMin, yMax int)
}

// A default implementation for Box and MutableBox interfaces.
type RawBox struct { XMinField int; XMaxField int; YMinField int; YMaxField int }
func (self *RawBox) XMin() int { return self.XMinField }
func (self *RawBox) XMax() int { return self.XMaxField }
func (self *RawBox) YMin() int { return self.YMinField }
func (self *RawBox) YMax() int { return self.YMaxField }
func (self *RawBox) MutateBox(xMin, xMax, yMin, yMax int) {
   self.XMinField, self.XMaxField, self.YMinField, self.YMaxField = xMin, xMax, yMin, yMax
}

func NewBox(xMin, xMax, yMin, yMax int) *RawBox {
	box := &RawBox { xMin, xMax, yMin, yMax }
	BoxPanicIfInvalid(box)
	return box
}

func NewBoxCopy(box Box) *RawBox {
	return &RawBox { box.XMin(), box.XMax(), box.YMin(), box.YMax() }
}
